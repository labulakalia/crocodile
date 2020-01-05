package schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	// Cron schedule loop
	Cron *cacheSchedule
)

const (
	defaultRPCTimeout       = time.Second * 3
	defaultHearbeatInterval = time.Second * 50
)

type task struct {
	name      string
	cronexpr  string
	close     chan struct{} // stop schedule
	running   bool          // task is running
	ctxcancel context.CancelFunc
	starttime int64 // run task time(ms)
}

type cacheSchedule struct {
	sync.RWMutex
	sch  map[string]*task
	sch1 sync.Map
}

// Init start run already exists task from db
func Init() error {
	Cron = &cacheSchedule{
		sch: make(map[string]*task),
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	eps, err := model.GetTasks(ctx)
	if err != nil {
		log.Error("GetTasks failed", zap.Error(err))
		return err
	}

	for _, t := range eps {
		Cron.Add(t.ID, t.Name, t.Cronexpr)
	}
	log.Info("init task success", zap.Int("Total", len(eps)))
	return nil
}

// Add task to schedule
func (s *cacheSchedule) Add(taskID, taskName string, cronExpr string) {
	s.Del(taskID)
	t := task{
		name:     taskName,
		cronexpr: cronExpr,
		close:    make(chan struct{}),
	}
	s.Lock()
	s.sch[taskID] = &t
	s.Unlock()
	log.Info("Add Task success", zap.String("taskid", taskID), zap.String("name", taskName))
	go s.runSchedule(taskID)
}

// del schedule
func (s *cacheSchedule) Del(id string) {
	t, ok := s.gettask(id)
	if ok {
		s.Lock()
		close(t.close)
		delete(s.sch, id)
		s.Unlock()
		log.Info("Del task success", zap.String("taskid", id))
		return
	}
}

// kill task
func (s *cacheSchedule) KillTask(taskid string) {
	t, exist := s.gettask(taskid)
	if !exist {
		log.Error("stoptask failed,task is not exist", zap.String("taskid", taskid))
		return
	}
	if t.ctxcancel != nil {
		t.ctxcancel()
	}
	return
}

func (s *cacheSchedule) gettask(id string) (*task, bool) {
	s.RLock()
	defer s.RUnlock()
	t, ok := s.sch[id]
	return t, ok
}

// add context cancel func to task
func (s *cacheSchedule) addctxcancel(taskid string, cancel context.CancelFunc) {
	t, exist := s.gettask(taskid)
	if !exist {
		log.Error("addctxcancel failed,task is not exist", zap.String("taskid", taskid))
		return
	}
	t.ctxcancel = cancel
}

// start run cronexpr schedule
func (s *cacheSchedule) runSchedule(taskid string) {
	t, exist := s.gettask(taskid)
	if !exist {
		return
	}
	log.Info("start run cronexpr", zap.Any("task", t), zap.String("id", taskid))

	sch, err := cron.ParseStandard(t.cronexpr)
	if err != nil {
		log.Error("ParseStandard", zap.Error(err))
		return
	}
	for {
		now := time.Now()
		next := sch.Next(now)
		select {
		case <-t.close:
			log.Info("Close Schedule", zap.String("taskID", taskid), zap.Any("task", t))
			return
		case <-time.After(next.Sub(now)):
			go s.RunTask(taskid)
		}
	}
}

func (s *cacheSchedule) RunTask(taskid string) {
	log.Info("start run task", zap.String("taskid", taskid))
	t, exist := s.gettask(taskid)
	if !exist {
		log.Error("taskid not exist", zap.String("taskid", taskid))
		return
	}

	if t.running {
		log.Info("task is running,so not run now", zap.String("taskid", taskid))
		return
	}
	t.running = true
	t.starttime = time.Now().Unix()
	defer func() {
		t.running = false
	}()
	tasklog, err := s.runTaskReal(taskid)
	if err != nil {
		log.Error("runtask failed", zap.String("taskid", taskid), zap.String("error", err.Error()))
		return
	}
	err = model.SaveLog(context.Background(), tasklog)
	if err != nil {
		log.Error("model.SaveLog failed", zap.String("error", err.Error()))
		return
	}
}

// start run task by execplanid
func (s *cacheSchedule) runTaskReal(id string) (*define.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	startTime := time.Now().UnixNano() / 1e6
	task, err := model.GetTaskByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "model.GetTaskByID")
	}

	resplog := define.Log{
		RunByTaskID: id,
	}

	if task.Run == 0 {
		return nil, errors.New(fmt.Sprintf("task %s forbid run", id))
	}
	// start run parent task
	if len(task.ParentTaskIds) != 0 {
		log.Info("start Run parenttasks", zap.Strings("taskids", task.ParentTaskIds))
		parentresplogs := s.runMultiTasks(task.ParentRunParallel, define.ParentTask, task.ParentTaskIds...)
		resplog.TaskResps = append(resplog.TaskResps, parentresplogs...)
	}
	// start run task
	log.Info("start run main task", zap.String("taskid", task.ID))
	runresp := s.runTask(task.ID, define.MasterTask)
	resplog.TaskResps = append(resplog.TaskResps, runresp)
	// start run childtasks
	if len(task.ChildTaskIds) != 0 {
		log.Info("start sun childtasks", zap.Strings("taskids", task.ChildTaskIds))
		childresplogs := s.runMultiTasks(task.ChildRunParallel, define.ChildTask, task.ChildTaskIds...)
		resplog.TaskResps = append(resplog.TaskResps, childresplogs...)
	}

	endTime := time.Now().UnixNano() / 1e6
	resplog.StartTime = startTime
	resplog.EndTime = endTime
	return &resplog, nil
}

// run multi tasks
func (s *cacheSchedule) runMultiTasks(RunParallel int, tasktype define.TaskRespType, taskids ...string) []*define.TaskResp {
	taskresp := make([]*define.TaskResp, 0, len(taskids))

	if RunParallel == 1 {
		var wg sync.WaitGroup
		wg.Add(len(taskids))
		for _, id := range taskids {
			go func(id string) {
				runresp := s.runTask(id, tasktype)
				taskresp = append(taskresp, runresp)
				wg.Done()
			}(id)
		}
		wg.Wait()
	} else {
		for _, id := range taskids {
			runresp := s.runTask(id, tasktype)
			taskresp = append(taskresp, runresp)
		}
	}

	return taskresp

}

// realy run task
func (s *cacheSchedule) runTask(id string, tasktype define.TaskRespType) *define.TaskResp {
	taskresp := &define.TaskResp{
		TaskID:   id,
		Code:     resp.ErrInternalServer,
		ErrMsg:   resp.GetMsg(resp.ErrInternalServer),
		TaskType: tasktype,
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task, err := model.GetTaskByID(ctx, id)
	if err != nil {
		log.Error("model.GetTaskByID failed", zap.String("taskid", id), zap.Error(err))
		return taskresp
	}

	hg, err := model.GetHostGroupID(ctx, task.HostGroupID)
	if err != nil {
		log.Error("model.GetTaskByID failed", zap.String("taskid", id), zap.Error(err))
		return taskresp
	}

	if len(hg.HostsID) == 0 {
		taskresp.Code = resp.ErrRPCNotValidHost
		taskresp.ErrMsg = resp.GetMsg(resp.ErrRPCNotValidHost)
		return taskresp
	}

	conn, err := tryGetRCCConn(ctx, hg)
	if err != nil {
		log.Error("tryGetRpcConn failed", zap.String("error", err.Error()))
		taskresp.Code = resp.ErrRPCNotConnHost
		taskresp.ErrMsg = resp.GetMsg(resp.ErrRPCNotConnHost)
		return taskresp
	}

	tdata, err := json.Marshal(task.TaskData)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		return taskresp
	}
	taskreq := &pb.TaskReq{
		TaskId:   task.ID,
		TaskType: int32(task.TaskType),
		TaskData: tdata,
		Timeout:  int32(task.Timeout),
	}
	var (
		taskcancel context.CancelFunc
		taskctx    context.Context
	)

	if task.Timeout > 0 {
		taskctx, taskcancel = context.WithTimeout(context.Background(), time.Second*time.Duration(task.Timeout))

	} else {
		taskctx, taskcancel = context.WithCancel(context.Background())
	}
	s.addctxcancel(task.ID, taskcancel)
	var errmsg []byte
	taskclient := pb.NewTaskClient(conn)
	rpcTaskResp, err := taskclient.RunTask(taskctx, taskreq)

	if err != nil {
		log.Error("RunTask failed", zap.Error(err))
		errcode := dealRPCErr(err)
		taskresp.Code = int32(errcode)
		taskresp.ErrMsg = resp.GetMsg(errcode)
	} else {
		var genresp bytes.Buffer
		errmsg = rpcTaskResp.ErrMsg
		genresp.Write(rpcTaskResp.RespData)
		taskresp.RespData = genresp.String()
		genresp.Reset()

		genresp.Write(rpcTaskResp.ErrMsg)
		taskresp.ErrMsg = genresp.String()
		taskresp.Code = rpcTaskResp.Code
	}
	if errmsg == nil {
		errmsg = []byte("")
	}

	taskresp.WorkerHost = conn.Target()
	return taskresp
}

// get all running task
func (s *cacheSchedule) GetRunningtask() []*define.RunTask {
	runtasks := []*define.RunTask{}
	s.RLock()
	for taskid, task := range s.sch {
		if !task.running {
			continue
		}
		runtask := define.RunTask{
			ID:            taskid,
			Name:          task.name,
			StartTime:     utils.UnixToStr(task.starttime),
			StartTimeUnix: task.starttime,
			RunTime:       int(time.Now().Unix() - task.starttime),
		}
		runtasks = append(runtasks, &runtask)
	}
	return runtasks
}
