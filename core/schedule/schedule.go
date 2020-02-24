package schedule

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/labulaka521/crocodile/common/errgroup"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/alarm"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/tasktype"
	"github.com/labulaka521/crocodile/core/utils/define"
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
	name        string
	cronexpr    string
	close       chan struct{}        // stop schedule
	running     bool                 // task is running
	stop        bool                 // stop run task
	ctxcancel   context.CancelFunc   // store cancelfunc could cancel all task by this cancel
	starttime   int64                // run task time(ms)
	endtime     int64                // end run task time(ms)
	logcaches   map[string]LogCacher // task runing logcaches
	taskids     []string             // save tasks。parent taskids、mainid、childids
	status      int                  // task run res fail: -1 success:1
	next        Next                 // it save a func Next by route policy
	Trigger     define.Trigger       // how to trigger task
	errTaskID   string               // run fail task's id
	errTask     string               // run fail task's id
	errCode     int                  // failed task return code
	errMsg      string               // task run failed errmsg
	errTasktype define.TaskRespType  // failed task type
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
	eps, _, err := model.GetTasks(ctx, 0, 0, "", "", "")
	if err != nil {
		log.Error("GetTasks failed", zap.Error(err))
		return err
	}

	for _, t := range eps {
		Cron.Add(t.ID, t.Name, t.Cronexpr, GetRoutePolicy(t.HostGroupID, t.RoutePolicy))
	}
	log.Info("init task success", zap.Int("Total", len(eps)))
	return nil
}

// Add task to schedule
func (s *cacheSchedule) Add(taskID, taskName string, cronExpr string, next Next) {
	log.Debug("Start Add task ", zap.String("name", taskName))

	t := task{
		name:     taskName,
		cronexpr: cronExpr,
		close:    make(chan struct{}),
		next:     next,
	}
	s.Lock()
	// 如果多个用户同时修改 确保不会冲突
	oldtask, exist := s.sch[taskID]
	if exist {
		close(oldtask.close)
		if oldtask.ctxcancel != nil {
			oldtask.ctxcancel()
		}
	}
	s.sch[taskID] = &t
	s.Unlock()
	log.Info("Add Task success", zap.String("taskid", taskID), zap.String("name", taskName))
	go s.runSchedule(taskID)
}

// Del schedule task
// if delete taskid,this taskid must be remove from other task's child or parent
func (s *cacheSchedule) Del(id string) {
	log.Info("start delete task", zap.String("taskid", id))
	task, exist := s.gettask(id)
	if exist {
		log.Debug("start clean ", zap.String("id", id))
		task.Lock()
		delete(s.sch, id)
		task.Unlock()
		go s.Clean(task)

	}

}

// Clean task
func (s *cacheSchedule) Clean(t *task) {
	log.Info("start clean task", zap.String("task", t.name))
	// if t.ctxcancel != nil {
	// 	t.ctxcancel()
	// }
	close(t.close)
	log.Info("clean task success", zap.String("task", t.name))
	return
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
	t, ok := s.sch[id]
	s.RUnlock()
	if ok && t.logcaches == nil {
		t.logcaches = make(map[string]LogCacher)
	}
	return t, ok
}

// saveLog save task resp log
func (s *cacheSchedule) storelog(runbyid string, t *task) error {
	// read all log
	// put logcache to locachepool
	tasklog := &define.Log{
		Name:         t.name,
		RunByTaskID:  runbyid,
		StartTime:    t.starttime,
		EndTime:      t.endtime,
		TotalRunTime: int(t.endtime - t.starttime),
		Status:       t.status,
		Trigger:      t.Trigger,
		ErrCode:      t.errCode,
		ErrMsg:       t.errMsg,
		ErrTasktype:  t.errTasktype,
		ErrTaskID:    t.errTaskID,
		ErrTask:      t.errTask,
		TaskResps:    make([]*define.TaskResp, 0, len(t.logcaches)),
	}

	for name, logcache := range t.logcaches {
		var (
			taskresp define.TaskResp
			ok       bool
		)
		// if ok,code runnhost tasktype valid
		taskresp, ok = logcache.Get().(define.TaskResp)

		id, tasktype, _ := splitname(name)
		// taskresp.Task = t.name
		if ok {
			taskresp.TaskID = id
			taskresp.TaskTypeStr = taskresp.TaskType.String()
			taskresp.LogData = logcache.ReadAll()

		} else {
			taskresp = define.TaskResp{
				TaskType:    tasktype,
				TaskTypeStr: tasktype.String(),
				TaskID:      id,
				LogData:     logcache.ReadAll(),
			}
		}
		if logcache.GetTaskStatus() == define.TsWait {
			taskresp.Status = define.TsCancel.String()
		} else {
			taskresp.Status = logcache.GetTaskStatus().String()
		}

		// taskresp
		logcache.Close()
		// Put coolpool after logcache close

		cachepool.Put(logcache)
		tasklog.TaskResps = append(tasklog.TaskResps, &taskresp)
	}

	go alarm.JudgeNotify(tasklog)
	// save log
	err := model.SaveLog(context.Background(), tasklog)

	return err
}

// runSchedule start run cronexpr schedule
func (s *cacheSchedule) runSchedule(taskid string) {
	t, exist := s.gettask(taskid)
	if !exist {
		return
	}
	log.Info("start run cronexpr", zap.Any("task", t), zap.String("id", taskid))

	expr, err := cronexpr.Parse(t.cronexpr)
	if err != nil {
		log.Error("cronexpr parse failed", zap.Error(err))
		return
	}
	var (
		last time.Time
		next time.Time
	)
	last = time.Now()
	for {
		log.Debug("all task", zap.Any("tasks", s.sch))
		next = expr.Next(last)
		select {
		case <-t.close:
			log.Info("close Schedule", zap.String("taskID", taskid), zap.Any("task", t))
			return
		case <-time.After(next.Sub(last)):
			last = next
			if t.stop {
				log.Error("task is stop run", zap.String("task", t.name))
			} else {
				go s.RunTask(taskid, define.Auto)
			}
		}
	}
}

func generatename(id string, tasktype define.TaskRespType) string {
	return id + "_" + strconv.Itoa(int(tasktype))
}

func splitname(taskname string) (string, define.TaskRespType, error) {
	res := strings.Split(taskname, "_")
	if len(res) != 2 {
		return "", 0, fmt.Errorf("split %s failed", taskname)
	}
	id := res[0]
	tasktype, err := strconv.Atoi(res[1])
	if err != nil {
		return "", 0, nil
	}
	return id, define.TaskRespType(tasktype), nil
}

// RunTask start run a task
func (s *cacheSchedule) RunTask(taskid string, trigger define.Trigger) {
	var (
		masterlogcache LogCacher
		ctx            context.Context
		cancel         context.CancelFunc
		err            error
		task           *define.GetTask
		g              *errgroup.Group
	)
	log.Info("start run task", zap.String("taskid", taskid))
	masterlogcache = cachepool.Get().(LogCacher) // this log cache is main
	masterlogcache.SetTaskStatus(define.TsWait)
	t, exist := s.gettask(taskid)
	if !exist {
		log.Error("this is bug, taskid not exist", zap.String("taskid", taskid), zap.Any("sch", s.sch))
		// logcache.WriteStringf("taskid %s not exist", taskid)
		cachepool.Put(masterlogcache)
		return
	}

	t.errTaskID = ""
	t.errTask = ""
	t.errCode = 0
	t.errMsg = ""
	t.errTasktype = 0

	t.Trigger = trigger
	masterlogcache.Clean()
	mastername := generatename(taskid, define.MasterTask)
	t.Lock()
	t.logcaches[mastername] = masterlogcache
	t.Unlock()

	// if master task is running,will do not run this time
	if t.running {
		log.Warn("task is running,so not run now", zap.String("task", t.name))
		cachepool.Put(masterlogcache)
		return
	}
	t.running = true
	t.starttime = time.Now().UnixNano() / 1e6

	ctx, cancel = context.WithCancel(context.Background())
	// save control ctx
	t.ctxcancel = cancel
	defer cancel()
	task, err = model.GetTaskByID(context.Background(), taskid)
	if err != nil {
		log.Error("can't get main taskID from dataBase", zap.String("task", task.Name))
		cachepool.Put(masterlogcache)
		return
	}

	if !task.Run {
		log.Error("main task is forbid run", zap.String("task", task.Name))
		cachepool.Put(masterlogcache)
		return
	}

	t.taskids = make([]string, 0, 1+len(task.ParentTaskIds)+len(task.ChildTaskIds))

	for _, id := range task.ParentTaskIds {
		name := generatename(id, define.ParentTask)
		logcache := cachepool.Get().(LogCacher)
		logcache.SetTaskStatus(define.TsWait)
		logcache.Clean()
		t.Lock()
		t.logcaches[name] = logcache
		t.Unlock()
		t.taskids = append(t.taskids, name)
	}

	t.taskids = append(t.taskids, mastername)

	for _, id := range task.ChildTaskIds {
		name := generatename(id, define.ChildTask)
		logcache := cachepool.Get().(LogCacher)
		logcache.SetTaskStatus(define.TsWait)
		logcache.Clean()
		t.Lock()
		t.logcaches[name] = logcache
		t.Unlock()
		t.taskids = append(t.taskids, name)
	}

	// if exist a err task,will stop all task
	g = errgroup.WithCancel(ctx)
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
		log.Error("task run failed", zap.String("taskid", taskid), zap.Error(err))
	}
	goto Over
Over:
	t.running = false
	t.endtime = time.Now().UnixNano() / 1e6

	if t.errTaskID == "" {
		t.status = 1
	}
	err = s.storelog(taskid, t)
	if err != nil {
		log.Error("save task log failed", zap.Error(err))
	}
}

// run multi tasks
// if hash one task err, will exit all task
// TODO: task run err whether influence  other task
// TODO: multi task set RunParallel total
func (s *cacheSchedule) runMultiTasks(ctx context.Context, RunParallel bool,
	tasktype define.TaskRespType, runbyid string, taskids ...string) error {
	if len(taskids) == 0 {
		return nil
	}
	log.Info("Start Run Task", zap.Strings("taskids", taskids), zap.String("tasktype", tasktype.String()))
	var maxproc int
	if RunParallel {
		maxproc = len(taskids)
	} else {
		maxproc = 1
	}
	g := errgroup.WithCancel(ctx)
	g.GOMAXPROCS(maxproc)
	for _, id := range taskids {
		taskid := id
		g.Go(func(ctx context.Context) error {
			return s.runTask(ctx, taskid, runbyid, tasktype)
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
		taskdata *define.GetTask
		// worker conn
		conn *grpc.ClientConn
		// task run data
		tdata []byte
		// recv grpc stream
		taskrespstream pb.Task_RunTaskClient
		// grpc client
		taskclient pb.TaskClient
		taskreq    *pb.TaskReq
		ctxcancel  context.CancelFunc
		taskctx    context.Context
		realtask   *task
		output     []byte
	)

	// when func exit,check res and judge whatever alarm

	runbytask, exist := s.gettask(runbyid)
	if !exist {
		// if happen,this is a bug,
		log.Error("this is a bug,task should exist,but can not get task,", zap.String("taskid", runbyid), zap.Any("allschedule", s.sch))
		err = fmt.Errorf("[bug] can not get taskid %s from schuedle: %v", id, s.sch)
		return err
	}

	if taskresptype == define.MasterTask {
		realtask = runbytask
	} else {
		realtask, exist = s.gettask(id)
		if !exist {
			log.Error("this is a bug,task should exist,but can not get task,", zap.String("taskid", runbyid), zap.Any("allschedule", s.sch))
			err = fmt.Errorf("[bug] can not get taskid %s from schuedle: %v", id, s.sch)
			return err
		}
	}
	name := generatename(id, taskresptype)
	runbytask.RLock()
	logcache, exist = runbytask.logcaches[name]
	runbytask.RUnlock()
	if !exist {
		// it happen,it is a bug
		warnbug := fmt.Sprintf("[bug] can get master task's %s logcache from logcaches: %v", id, runbytask.logcaches)
		log.Error(warnbug)
		logcache = cachepool.Get().(LogCacher)
		logcache.WriteString(warnbug)
		runbytask.Lock()
		runbytask.logcaches[name] = logcache
		runbytask.Unlock()
	}

	logcache.SetTaskStatus(define.TsRun)

	queryctx, querycancel := context.WithTimeout(ctx,
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer querycancel()

	// TODO cache task run data and hostgroup
	taskdata, err = model.GetTaskByID(queryctx, id)
	if err != nil {
		log.Error("model.GetTaskByID failed", zap.String("taskid", id), zap.Error(err))
		logcache.WriteStringf("Get Task id %s failed: %v", id, err)
		goto Check
	}
	logcache.WriteStringf("Start Prepare Task %s[%s]", taskdata.Name, id)
	logcache.WriteStringf("Start Conn Worker Host For Task %s[%s]", taskdata.Name, id)

	conn, err = tryGetRCCConn(ctx, realtask.next)
	if err != nil {
		log.Error("tryGetRpcConn failed", zap.String("error", err.Error()))
		err = fmt.Errorf("Get Rpc Conn Failed From Hostgroup %s[%s] Err: %v",
			taskdata.HostGroup, taskdata.HostGroupID, err)
		logcache.WriteStringf(err.Error())
		goto Check
	}

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

	logcache.WriteStringf("Success Get Task %s[%s] Run Data", taskdata.Name, id)

	logcache.WriteStringf("Start Run Task %s[%s] On Host[%s]", taskdata.Name, id, conn.Target())

	// taskctx only use RunTask
	if taskdata.Timeout > 0 {
		taskctx, ctxcancel = context.WithTimeout(ctx, time.Second*time.Duration(taskdata.Timeout))
	} else {
		taskctx, ctxcancel = context.WithCancel(ctx)
	}

	defer ctxcancel()
	taskclient = pb.NewTaskClient(conn)

	taskrespstream, err = taskclient.RunTask(taskctx, taskreq)
	if err != nil {
		log.Error("Run task failed", zap.Error(err))
		logcache.WriteStringf("Run Task %s[%s] TaskData [%v] failed:%v", taskdata.Name, id, taskreq, err)
		goto Check
	}

	// RunTask default resp code

	logcache.WriteStringf("Task %s[%s]  Output----------------", taskdata.Name, id)
	for {
		// Recv return err is nil or io.EOF
		// the last lastrecv must be return code 3 byte
		taskresp, err = taskrespstream.Recv()
		if err != nil {
			if err == io.EOF {
				err = nil
				taskrespcode = logcache.GetCode()
				goto Check
			}
			err = DealRPCErr(err)
			// taskrespcode =
			logcache.WriteStringf("Task %s[%s] Run Fail: %v", taskdata.Name, id, err.Error())
			// Alarm
			log.Error("Recv failed", zap.Error(err))
			// err = resp.GetMsgErr(taskrespcode)
			goto Check
		}
		logcache.Write(taskresp.GetResp())
		output = append(output, taskresp.GetResp()...)
	}
	// logcache.WriteStringf("Task %s[%s] End Output-------------------\n", taskdata.Name, id)

Check:
	// logcache.WriteStringf("Task %s[%s] Run Over\n", taskdata.Name, id)
	// 记录任务的状态

	// save task returncode,tasktype,if task could find run host,run host will be save hear
	tmptaskresp := define.TaskResp{
		Task:     realtask.name,
		Code:     taskrespcode,
		TaskType: taskresptype,
	}
	if conn != nil {
		// if conn worker failed,can not get worker host
		tmptaskresp.RunHost = conn.Target()
	}
	logcache.Save(tmptaskresp)
	// 当终止任务时，第一个任务取消的任务不经过这里处理，后续的任务才会经过这里处理
	// 所以需要判断t.errTaskId 为空时才经过这里处理
	if err != nil && runbytask.errTaskID != "" {
		select {
		case <-ctx.Done():
			log.Error("task is cancel", zap.String("task", realtask.name))
			logcache.WriteStringf("task %s[%s] is canceled", realtask.name, id)
			logcache.SetTaskStatus(define.TsCancel)
			return nil
		default:
		}
	}
	var alarmerr error
	// if a task fail other task will return context.Canceled,but it can not alarm
	// because the first err task always alarm,so other task do not alarm
	// and the first err task's errmsg will save tasking

	// check task resp code and resp content
	judgeres := func() error {
		if err != nil {
			return err
		}
		if taskdata.ExpectCode != taskrespcode {
			return fmt.Errorf("%s task %s[%s] resp code is %d,want resp code %d", taskresptype.String(), id, taskdata.Name, taskrespcode, taskdata.ExpectCode)
		}
		if taskdata.ExpectContent != "" {
			if !strings.Contains(string(output), taskdata.ExpectContent) {
				return fmt.Errorf("%s task %s[%s] resp context not contains expect content: %s", taskresptype.String(), id, taskdata.Name, taskdata.ExpectContent)
			}
		}
		return nil
	}
	alarmerr = judgeres()
	if alarmerr != nil {
		//  task run err
		// 只运行到这里一次
		// runbytask.status = -1 // task fail
		log.Error("task run fail", zap.String("task", realtask.name), zap.Error(err))
		if runbytask.errTaskID == "" {
			runbytask.status = -1
			runbytask.errTaskID = id
			runbytask.errTask = realtask.name
			runbytask.errCode = taskrespcode
			runbytask.errMsg = alarmerr.Error()
			runbytask.errTasktype = taskresptype
			// log.Error("-----------------task run fail", zap.String("task", realtask.name), zap.Error(err))
			logcache.SetTaskStatus(define.TsFail)
		}

	} else {
		log.Error("task run success", zap.String("task", realtask.name))
		logcache.SetTaskStatus(define.TsFinish)
		// 如有任务失败，那么还未运行的任务可以标记为取消
		// 如果失败的任务还存在并行任务，那么
	}
	return alarmerr
}

//  sort running task
type runningTask []*define.RunTask

func (rt runningTask) Len() int { return len(rt) }
func (rt runningTask) Less(i, j int) bool {
	ii, err := strconv.Atoi(rt[i].ID)
	if err != nil {
		log.Error("change ID to int failed", zap.String("id", rt[i].ID))
	}
	jj, err := strconv.Atoi(rt[j].ID)
	if err != nil {
		log.Error("change ID to int failed", zap.String("id", rt[j].ID))
	}
	return ii < jj
}
func (rt runningTask) Swap(i, j int) { rt[i], rt[j] = rt[j], rt[i] }

// GetRunningtask return all running task
func (s *cacheSchedule) GetRunningtask() []*define.RunTask {
	runtasks := runningTask{}
	s.RLock()
	defer s.RUnlock()
	for taskid, task := range s.sch {
		if !task.running {
			continue
		}
		// task.running
		runtask := define.RunTask{
			ID:           taskid,
			Name:         task.name,
			StartTimeStr: utils.UnixToStr(task.starttime / 1e3),
			StartTime:    task.starttime,
			RunTime:      int(time.Now().Unix() - task.starttime/1e3),
			Trigger:      task.Trigger.String(),
			Cronexpr:     task.cronexpr,
		}
		runtasks = append(runtasks, &runtask)
	}
	sort.Sort(runtasks)
	return runtasks
}

// GetRunTaskStaus return
func (s *cacheSchedule) GetRunTaskStaus(taskid string) []*define.TaskStatusTree {
	retTasksStatus := define.GetTasksTreeStatus()

	parentTasksStatus := retTasksStatus[0]

	taskinfo, exist := s.gettask(taskid)
	if !exist {
		return nil
	}
	mainTaskStatus := retTasksStatus[1]
	mainTaskStatus.Name = taskinfo.name
	mainTaskStatus.ID = taskid

	childTasksStatus := retTasksStatus[2]

	task, exist := s.gettask(taskid)
	if !exist {
		return nil
	}
	var status = define.TsNoData
	var isSet = false
	var change = false
	log.Debug("start get task run status", zap.Strings("ids", task.taskids))
	for _, name := range task.taskids {
		// if !task.running {
		// 	log.Error("task is not run", zap.String("name", name))
		// 	return nil
		// }
		task.RLock()
		logcache := task.logcaches[name]
		task.RUnlock()
		id, _, _ := splitname(name)

		taskinfo, _ := s.gettask(id)

		if name == generatename(taskid, define.MasterTask) {
			parentTasksStatus.Status = status.String()
			isSet = false
			status = define.TsNoData
			change = true
			// main task
			// log.Debug("start get main status" + id + ":" + logcache.GetTaskStatus().String())
			mainTaskStatus.Status = logcache.GetTaskStatus().String()
			mainTaskStatus.TaskType = define.MasterTask
		} else if change == false {
			// parent taskids
			// 如果全部为wait就设置为wait
			// 如果全部成功那么就设置为finish

			// 如果任务存在run那么就设置为run
			// 如果任务有失败那么就设置为fail
			if !isSet {
				if logcache.GetTaskStatus() == define.TsRun || logcache.GetTaskStatus() == define.TsFail {
					isSet = true
					status = logcache.GetTaskStatus()
				} else {
					status = logcache.GetTaskStatus()
				}

			}
			// log.Debug("start get parent status" + id + ":" + status.String())
			parentaskStatus := &define.TaskStatusTree{
				Name:     taskinfo.name,
				ID:       id,
				TaskType: define.ParentTask,
				Status:   logcache.GetTaskStatus().String(),
			}
			parentTasksStatus.TaskType = define.ParentTask
			parentTasksStatus.Children = append(parentTasksStatus.Children, parentaskStatus)

		} else {
			if !isSet {
				if logcache.GetTaskStatus() == define.TsRun || logcache.GetTaskStatus() == define.TsFail {
					isSet = true
					status = logcache.GetTaskStatus()
				} else {
					status = logcache.GetTaskStatus()
				}
			}
			// log.Debug("start get parent status" + id + ":" + status.String())

			// child taskids
			childaskStatus := &define.TaskStatusTree{
				Name:     taskinfo.name,
				ID:       id,
				TaskType: define.ChildTask,
				Status:   logcache.GetTaskStatus().String(),
			}
			childTasksStatus.TaskType = define.ChildTask
			childTasksStatus.Children = append(childTasksStatus.Children, childaskStatus)
		}
	}
	childTasksStatus.Status = status.String()

	log.Debug("TaskRunStatus", zap.Any("status", retTasksStatus))
	return retTasksStatus
}

// GetLogCache get log cache
func (s *cacheSchedule) GetRunTaskLogCache(taskid, realtaskid string, tasktype define.TaskRespType) (LogCacher, error) {
	t, exist := s.gettask(taskid)
	if !exist {
		return nil, fmt.Errorf("can get task id %s", taskid)
	}
	if !t.running {
		return nil, fmt.Errorf("main task %s[%s] is not running,can't get running log", t.name, taskid)
	}
	name := generatename(realtaskid, tasktype)
	t.RLock()
	logcache, exist := t.logcaches[name]
	t.RUnlock()
	if !exist {
		return nil, fmt.Errorf("can get task's logcache id %s", realtaskid)
	}
	return logcache, nil
}
