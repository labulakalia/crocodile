package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/version"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"os"

	"sync"
	"time"
)

var (
	Cron *cacheSchedule
)

const (
	DefaultRpcTimeout = time.Second * 3
)

type task struct {
	cronexpr string
	close    chan struct{}
	running  bool // task is running
}

type cacheSchedule struct {
	sync.RWMutex
	sch map[string]*task
}

// start run already exists task from db
func InitServer() error {
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
		Cron.Add(t.Id, t.CronExpr)
	}
	log.Info("Init task success", zap.Int("Total", len(eps)))
	return nil
}

func InitClient(port int) error {

	conn, err := NewgRPCConn(config.CoreConf.Client.ServerAddr)
	if err != nil {
		return err
	}
	hbClient := pb.NewHeartbeatClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRpcTimeout)
	defer cancel()
	hostname, _ := os.Hostname()
	regHost := pb.RegistryReq{
		Port:      int32(port),
		Hostname:  hostname,
		Version:   version.Version,
		Hostgroup: config.CoreConf.Client.HostGroup,
	}
	_, err = hbClient.RegistryHost(ctx, &regHost)
	if err != nil {
		return err
	}
	log.Info("Host Registry Success")
	return nil
}

// add task to schedule
func (s *cacheSchedule) Add(taskId string, cronExpr string) {
	s.Del(taskId)
	t := task{
		cronexpr: cronExpr,
		close:    make(chan struct{}),
	}
	s.Lock()
	s.sch[taskId] = &t
	s.Unlock()
	log.Info("Add Task success", zap.String("taskid", taskId))
	go s.runSchedule(taskId)
}

// start run cronexpr schedule
func (s *cacheSchedule) runSchedule(taskid string) {
	t, exist := s.gettask(taskid)
	if !exist {
		return
	}
	log.Info("Start Run Cronexpr", zap.Any("task", t), zap.String("id", taskid))

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
			go func() {
				if t.running {
					return
				}
				t.running = true
				defer func() {
					t.running = false
				}()
				tasklog, err := s.RunTask(taskid)
				if err != nil {
					log.Error("RunTask failed", zap.String("taskid", taskid), zap.String("error", err.Error()))
					return
				}
				err = model.SaveLog(context.Background(), tasklog)
				if err != nil {
					log.Error("model.SaveLog failed", zap.String("error", err.Error()))
					return
				}

			}()
		}
	}
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

func (s *cacheSchedule) gettask(id string) (*task, bool) {
	s.Lock()
	defer s.Unlock()

	t, ok := s.sch[id]
	return t, ok
}

// start run task by execplanid
func (s *cacheSchedule) RunTask(id string) (*define.Log, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task, err := model.GetTaskByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "model.GetTaskByID")
	}
	startTime := time.Now().Unix()
	resplog := define.Log{
		RunByTaskId: id,
	}

	if task.Run == 0 {
		return nil, errors.New(fmt.Sprintf("task %s forbid run", id))
	}
	if len(task.ParentTaskIds) != 0 {
		log.Info("Start Run ParentTasks", zap.Strings("taskids", task.ParentTaskIds))
		parentresplogs := s.runMultiTasks(task.ParentRunParallel, define.ParentTask, task.ParentTaskIds...)
		resplog.TaskResps = append(resplog.TaskResps, parentresplogs...)
	}
	taskresp, err := s.runTask(task.Id)
	if err != nil {
		log.Error("runTask failed", zap.String("taskid", task.Id), zap.String("error", err.Error()))
		taskresp = &define.TaskResp{
			TaskId:   id,
			Code:     -1,
			ErrMsg:   []byte("runTask failed: " + err.Error()),
			Data:     nil,
			TaskType: define.MasterTask,
		}
	}
	resplog.TaskResps = append(resplog.TaskResps, taskresp)

	if len(task.ChildTaskIds) != 0 {
		log.Info("Start Run ChildTasks", zap.Strings("taskids", task.ChildTaskIds))
		childresplogs := s.runMultiTasks(task.ChildRunParallel, define.ChildTask, task.ChildTaskIds...)
		resplog.TaskResps = append(resplog.TaskResps, childresplogs...)
	}
	endTime := time.Now().Unix()
	resplog.TotalRunTime = int(endTime - startTime)
	resplog.StartTimne = utils.UnixToStr(startTime)
	resplog.EndTime = utils.UnixToStr(endTime)
	return &resplog, nil
}

func (s *cacheSchedule) runMultiTasks(RunParallel int, tasktype define.TaskRespType, taskids ...string) []*define.TaskResp {
	taskresp := make([]*define.TaskResp, 0, len(taskids))

	if RunParallel == 1 {
		var wg sync.WaitGroup
		wg.Add(len(taskids))
		for _, id := range taskids {
			go func(id string) {
				resp, err := s.runTask(id)
				if err != nil {
					log.Error("runTask failed", zap.String("taskid", id), zap.String("error", err.Error()))
					resp = &define.TaskResp{
						TaskId: id,
						Code:   -1,
						ErrMsg: []byte("runTask failed: " + err.Error()),
						Data:   nil,
					}
				}
				resp.TaskType = tasktype
				taskresp = append(taskresp, resp)

				wg.Done()
			}(id)
		}
		wg.Wait()
	} else {
		for _, id := range taskids {
			resp, err := s.runTask(id)
			if err != nil {
				log.Error("runTask failed", zap.String("taskid", id), zap.String("error", err.Error()))
				resp = &define.TaskResp{
					TaskId: id,
					Code:   -1,
					ErrMsg: []byte("runTask failed: " + err.Error()),
					Data:   nil,
				}
			}
			resp.TaskType = tasktype
			taskresp = append(taskresp, resp)
		}
	}

	return taskresp

}

// chan log
func (s *cacheSchedule) runTask(id string) (*define.TaskResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task, err := model.GetTaskByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "model.GetTaskByID")
	}
	log.Info("Start runTask", zap.String("taskid", task.Id))
	hg, err := model.GetHostGroupID(ctx, task.HostGroupId)
	if err != nil {
		return nil, errors.Wrap(err, "model.GetHostGroupID")
	}

	if len(hg.Addrs) == 0 {
		return nil, errors.New("hostgroup not exist run host")
	}

	conn, err := tryGetRpcConn(ctx, hg)
	if err != nil {
		return nil, errors.Wrap(err, "s.TryGetConn")
	}
	taskclient := pb.NewTaskClient(conn)
	tdata, err := json.Marshal(task.TaskData)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}
	taskreq := &pb.TaskReq{
		TaskId:   task.Id,
		TaskType: int32(task.TaskType),
		TaskData: tdata,
		Timeout:  int32(task.Timeout),
	}
	taskctx := context.Background()
	if task.Timeout != 0 {
		taskctx, cancel = context.WithTimeout(context.Background(), time.Second*time.Duration(task.Timeout))
		defer cancel()
	}
	rpcTaskResp, err := taskclient.RunTask(taskctx, taskreq)
	var errmsg []byte
	if err != nil {
		statusErr, ok := status.FromError(err)
		if !ok {
			return nil, err
		}
		if statusErr.Code() != codes.DeadlineExceeded {
			fmt.Printf("%+v\n", err)
			return nil, err
		}
		errmsg = []byte("task run timeout")
	} else {
		errmsg = rpcTaskResp.ErrMsg
	}
	if errmsg == nil {
		errmsg = []byte("")
	}
	if rpcTaskResp.RespData == nil {
		rpcTaskResp.RespData = []byte("")
	}
	log.Info("runTask Resp", zap.Any("resp", rpcTaskResp))

	logresp := define.TaskResp{
		TaskId: id,
		Code:   rpcTaskResp.Code,
		ErrMsg: errmsg,
		Data:   rpcTaskResp.RespData,
	}
	return &logresp, nil
}

// get rpc conn
func tryGetRpcConn(ctx context.Context, hg *define.HostGroup) (*grpc.ClientConn, error) {
	i := 0
	for i < len(hg.Addrs) {
		i++
		addr, err := model.RandAddrByGostGroup(ctx, hg)
		if err != nil {
			log.Error("model.RandHostByGostGroup failed", zap.String("error", err.Error()))
			continue
		}
		conn, err := NewgRPCConn(addr)
		if err != nil {
			log.Error("GetRpcConn failed", zap.String("error", err.Error()))
			continue
		}
		// idle
		if conn.GetState() <= connectivity.Ready {
			return conn, nil
		}
		conn.Close()
	}
	return nil, errors.New("can not get valid grpc conn")
}
