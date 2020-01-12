package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/labulaka521/crocodile/common/errgroup"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/alarm"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/tasktype"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// Cron schedule loop
	Cron *cacheSchedule
)

// task running status
type task struct {
	sync.RWMutex
	once sync.Once
	name      string
	cronexpr  string
	close     chan struct{}        // stop schedule
	running   bool                 // task is running
	ctxcancel context.CancelFunc   // store cancelfunc could cancel all task by this cancel
	starttime int64                // run task time(ms)
	endtime   int64                // end run task time(ms)
	logcaches map[string]LogCacher // task runing logcaches

	// // unexpectlogcache LogCacher            // unexecpt log,if has log, will exist bug
	// runninghost     string              // current run task on host
	// runningtask     string              // running task id parent,master,child task
	// runningtasktype define.TaskRespType // running task id parent,master,child task type
	errTaskID   string              // run fail task's id
	errCode     int                 // failed task return code
	errMsg      string              // task run failed errmsg
	errTasktype define.TaskRespType // failed task type
}

type cacheSchedule struct {
	sync.RWMutex
	sch map[string]*task
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

// Del schedule task
// if delete taskid,this taskid must be remove from other task's child or parent
func (s *cacheSchedule) Del(id string) {
	t, ok := s.gettask(id)
	if ok {
		if t.running && t.ctxcancel != nil {
			t.ctxcancel()
		}
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
	if ok && t.logcaches == nil {
		t.logcaches = make(map[string]LogCacher)
	}
	return t, ok
}

// saveLog save task resp log
func (s *cacheSchedule) saveLog(t *task) error {
	// read all log
	// put logcache to locachepool
	return nil
}

// runSchedule start run cronexpr schedule
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

// RunTask start run a task
func (s *cacheSchedule) RunTask(taskid string) {
	var masterlogcache LogCacher
	log.Info("start run task", zap.String("taskid", taskid))
	masterlogcache = cachepool.Get().(LogCacher) // this log cache is main
	t, exist := s.gettask(taskid)
	if !exist {
		log.Error("this is bug, taskid not exist", zap.String("taskid", taskid), zap.Any("sch", s.sch))
		// logcache.WriteStringf("taskid %s not exist", taskid)
		cachepool.Put(masterlogcache)
		return
	}
	t.Lock()
	t.logcaches[taskid] = masterlogcache
	t.Unlock()

	// if master task is running,will do not run this time
	if t.running {
		log.Info("task is running,so not run now", zap.String("taskid", taskid))
		masterlogcache.WriteStringf("taskid %s is running, so not run now", taskid)
		return
	}
	t.running = true
	t.starttime = time.Now().Unix()
	defer func() {
		t.running = false
		t.endtime = time.Now().Unix()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	// save control ctx
	t.ctxcancel = cancel
	task, err := model.GetTaskByID(context.Background(), taskid)
	if err != nil {
		log.Error("model.GettaskById failed", zap.String("id", taskid), zap.Error(err))
		masterlogcache.WriteStringf("can not get task %s from db", taskid)
		return
	}

	// TODO delete judge run, onlu use it in cron
	if task.Run == 0 {
		log.Error("model.GettaskById failed", zap.Error(err))
		masterlogcache.WriteStringf("task %s[%s] is forbid run", task.Name, taskid)
		return
	}

	// if exist a err task,will stop all task
	g := errgroup.WithCancel(ctx)
	g.GOMAXPROCS(1)
	// parent tasks
	g.Go(func(ctx context.Context) error {
		return s.runMultiTasks(ctx, task.ParentRunParallel, define.ParentTask, task.ID, task.ParentTaskIds...)
	})
	// master task
	g.Go(func(ctx context.Context) error {
		return s.runTask(ctx, task.ID, task.ID, define.MasterTask)
	})
	// childs task
	g.Go(func(ctx context.Context) error {
		return s.runMultiTasks(ctx, task.ChildRunParallel, define.ChildTask, task.ID, task.ChildTaskIds...)
	})
	err = g.Wait()
	if err != nil {
		log.Error("run failed", zap.Error(err))
	}
	// TODO save log
	s.saveLog(t)
}

// run multi tasks
// if hash one task err, will exit all task
// TODO: task run err whether influence  other task
// TODO: multi task set RunParallel total
func (s *cacheSchedule) runMultiTasks(ctx context.Context, RunParallel int,
	tasktype define.TaskRespType, runbyid string, taskids ...string) error {
	if len(taskids) == 0 {
		return nil
	}
	log.Info("Start Run Task", zap.Strings("taskids", taskids), zap.String("tasktype", tasktype.String()))
	var maxproc int
	if RunParallel == 1 {
		maxproc = len(taskids)
	} else {
		maxproc = 1
	}
	g := errgroup.WithCancel(ctx)
	g.GOMAXPROCS(maxproc)
	for _, id := range taskids {
		g.Go(func(ctx context.Context) error {
			return s.runTask(ctx, id, runbyid, tasktype)
		})
	}
	err := g.Wait()
	return err
}

// runTask start run task,log will store
func (s *cacheSchedule) runTask(ctx context.Context, id, /*real run task id*/
	runbyid /*run by id*/ string, taskresptype define.TaskRespType) error {
	var (
		// real need task status
		// realtask *task
		logcache LogCacher

		// recverr      error
		taskrespcode = tasktype.DefaultExitCode
		// recv grpc stream
		taskresp *pb.TaskResp
		// error
		err error
		// task data
		taskdata *define.Task
		// task's hostgroup
		hg *define.HostGroup
		// worker conn
		conn *grpc.ClientConn
		// task run data
		tdata []byte
		// recv grpc stream
		taskrespstream pb.Task_RunTaskClient
		// grpc client
		taskclient pb.TaskClient
		taskreq    *pb.TaskReq
		cancel     context.CancelFunc
		taskctx    context.Context
	)
	var output []byte
	// whenn func exit,check res and judge whatever alarm

	runbytask, exist := s.gettask(runbyid)
	if !exist {
		// if happen,this is a bug,
		log.Error("this is a bug,task should exist,but can not get task,", zap.String("taskid", runbyid), zap.Any("allschedule", s.sch))
		err = fmt.Errorf("[bug] can not get taskid %s from schuedle: %v", id, s.sch)
		return err
	}
	if taskresptype == define.MasterTask {
		runbytask.RLock()
		logcache, exist = runbytask.logcaches[id]
		runbytask.RUnlock()
		if !exist {
			// it happen,it is a bug
			warnbug := fmt.Sprintf("[bug] can get master task's %s logcache from logcaches: %v", id, runbytask.logcaches)
			log.Error(warnbug)
			logcache = cachepool.Get().(LogCacher)
			logcache.WriteString(warnbug)
			runbytask.Lock()
			runbytask.logcaches[id] = logcache
			runbytask.Unlock()
		}
	} else {
		logcache = cachepool.Get().(LogCacher)
		runbytask.Lock()
		runbytask.logcaches[id] = logcache
		runbytask.Unlock()
	}


	queryctx, querycancel := context.WithTimeout(ctx,
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer querycancel()


	// TODO cache task run data and hostgroup
	taskdata, err = model.GetTaskByID(queryctx, id)
	if err != nil {
		log.Error("model.GetTaskByID failed", zap.String("taskid", id), zap.Error(err))
		logcache.WriteStringf("Get Task id %s failed: %v", id, err)
		// return err
		goto Check
	}
	logcache.WriteStringf("Start Prepare Task %s[%s]", taskdata.Name, id)
	logcache.WriteStringf("Start Conn Worker Host For Task %s[%s]", taskdata.Name, id)
	// Get hsotgroup
	hg, err = model.GetHostGroupID(queryctx, taskdata.HostGroupID)
	if err != nil {
		log.Error("model.GetTaskByID failed", zap.String("taskid", id), zap.Error(err))
		logcache.WriteStringf("Get Task %s[%s]'s HostGroup failed: %v", taskdata.Name, id, err)
		goto Check
	}
	// check hostgroup
	if len(hg.HostsID) == 0 {
		err = fmt.Errorf("HostGroup %s[%s] not valid worker", hg.Name, taskdata.HostGroupID)
		log.Error("Get failed", zap.Error(err))
		logcache.WriteString(err.Error())
		goto Check
	}
	// try conn hg's host
	conn, err = tryGetRCCConn(ctx, hg)
	if err != nil {
		log.Error("tryGetRpcConn failed", zap.String("error", err.Error()))
		logcache.WriteStringf("Get Conn from hostgroup %s[%s] failed: %v", taskdata.HostGroupID, hg.Name, err)
		goto Check
	}
	// runbytask.runninghost = conn.Target()

	logcache.WriteStringf("Success Conn Worker Host[%s]", conn.Target())

	logcache.WriteStringf("Start Get Task %s[%s] Run Data", taskdata.Name, id)
	// Marshal task data
	tdata, err = json.Marshal(taskdata.TaskData)
	if err != nil {
		log.Error("json.Marshal", zap.Error(err))
		logcache.WriteStringf("Marshal task %s[%s]'s RunData [%v] failed: %v", taskdata.Name, id, taskdata.TaskData, err)
		goto Check
	}

	// task run data
	taskreq = &pb.TaskReq{
		TaskId:   id,
		TaskType: int32(taskdata.TaskType),
		TaskData: tdata,
	}

	logcache.WriteStringf("Success Get Task %s[%s] Run Data ", taskdata.Name, id)

	logcache.WriteStringf("Start Run Task %s[%s] On Host[%s]", taskdata.Name, id, conn.Target())

	// taskctx only use RunTask
	if taskdata.Timeout > 0 {
		taskctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(taskdata.Timeout))
	} else {
		taskctx, cancel = context.WithCancel(ctx)
	}

	defer cancel()
	taskclient = pb.NewTaskClient(conn)

	taskrespstream, err = taskclient.RunTask(taskctx, taskreq)
	if err != nil {
		log.Error("Run task failed", zap.Error(err))
		logcache.WriteStringf("Run Task %s[%s] TaskData [%v] failed:%v", taskdata.Name, id, taskreq, err)
		goto Check
	}

	// RunTask default resp code

	logcache.WriteStringf("---------------- Task %s[%s] Start Output----------------", taskdata.Name, id)
	// defer logcache.WriteStringf("---------------- Task %s[%s] Start Output----------------", realtask.name, id)
	for {
		// Recv return err is nil or io.EOF
		// the last lastrecv must be return code 3 byte
		taskresp, err = taskrespstream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
				taskrespcode = logcache.GetCode()
				break
			}
			taskrespcode = dealRPCErr(err)

			logcache.WriteStringf("Task %s[%s] Run Fail: %v", taskdata.Name, id, resp.GetMsg(taskrespcode))
			// Alarm
			log.Error("Recv failed", zap.Error(err))
			err = resp.GetMsgErr(taskrespcode)
			break
		}
		logcache.Write(taskresp.GetResp())
		output = append(output, taskresp.GetResp()...)
	}
	logcache.WriteStringf("---------------- Task %s[%s] End Output-------------------", taskdata.Name, id)
	// return err
	goto Check
Check:
	var errmsg string
	if err != nil {
		errmsg = " Error" + err.Error()
	}
	logcache.WriteStringf("\nTask %s[%s] Run Over Return Code: %d"+errmsg, taskdata.Name, id, taskrespcode)
	
	var alarmerr error
	// if a task fail other task will return context.Canceled,but it can not alarm
	// because the first err task always alarm,so other task do not alarm
	// and the first err task's errmsg will save tasking
	runbytask.once.Do(func() {
		// check task resp
		alarmerr := alarm.CheckAlarm(id, runbyid, taskresptype, taskrespcode, output, err)
		if alarmerr != nil {
			//  task run err
			runbytask.errTaskID = id
			runbytask.errCode = taskrespcode
			runbytask.errMsg = alarmerr.Error()
			runbytask.errTasktype = taskresptype
		}
	})

	return alarmerr
}

// GetRunningtask return all running task
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
