package execute

import (
	"context"
	pbexecutor "crocodile/service/executor/proto/executor"
	pbjob "crocodile/service/job/proto/job"
	pbtasklog "crocodile/service/tasklog/proto/tasklog"
	"crocodile/third_party/github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/client"
	"os/exec"
	"sync"
	"time"
)

var (
	// 执行过程中的信息
	ExecutingTable map[string]*TaskExecuteInfo
	lock           *sync.RWMutex
	tasklogclient  pbtasklog.TaskLogService
)

func Init(client client.Client) {
	ExecutingTable = make(map[string]*TaskExecuteInfo)
	lock = &sync.RWMutex{}
	tasklogclient = pbtasklog.NewTaskLogService("crocodile.srv.tasklog", client)
}

// 执行中的信息
type TaskExecuteInfo struct {
	Task       *pbjob.Task
	CancelFunc context.CancelFunc
}

func NewTaskExecuteInfo(cancelFunc context.CancelFunc, task *pbjob.Task) (taskExecuteInfo *TaskExecuteInfo) {
	taskExecuteInfo = &TaskExecuteInfo{
		Task:       task,
		CancelFunc: cancelFunc,
	}
	return
}

// 运行任务
func RunTask(ctx context.Context, executeMsg *pbexecutor.ExecuteMsg) (err error) {
	var (
		cancelFunc context.CancelFunc
		cancelCtx  context.Context
		exits      bool

		startTime *timestamp.Timestamp
		endTime   *timestamp.Timestamp
		taskLog   *pbtasklog.TaskResLog
		cmd       *exec.Cmd
		output    []byte
	)

	logging.Infof("Start Run Task %s", executeMsg.Task.Taskname)
	// 在运行中的任务不允许再次调度
	// 加读锁 因为map不是并发安全的
	lock.RLock()
	_, exits = ExecutingTable[executeMsg.Task.Taskname]
	lock.RUnlock()
	if exits {
		logging.Errorf("Task %s:%s Is Running", executeMsg.Task.Taskname, executeMsg.Task.Id)
		return
	}

	// 未设置任务的超时时间
	if executeMsg.Task.Timeout == 0 {
		cancelCtx, cancelFunc = context.WithCancel(ctx)
	} else {
		cancelCtx, cancelFunc = context.WithTimeout(ctx, time.Duration(executeMsg.Task.Timeout)*time.Second)
	}

	// 更新执行表
	// 需要加写锁
	lock.Lock()
	ExecutingTable[executeMsg.Task.Taskname] = NewTaskExecuteInfo(cancelFunc, executeMsg.Task)
	lock.Unlock()

	// 执行的日志
	startTime, _ = ptypes.TimestampProto(time.Now())
	taskLog = &pbtasklog.TaskResLog{
		Executemsg: executeMsg,
		StartTime:  startTime,
	}

	// 运行任务的命令
	cmd = exec.CommandContext(cancelCtx, "/bin/bash", "-c", executeMsg.Task.Command)
	output, err = cmd.CombinedOutput()

	taskLog.Output = string(output)
	if err != nil {
		taskLog.Err = err.Error()
	}
	endTime, _ = ptypes.TimestampProto(time.Now())
	taskLog.EndTime = endTime
	// 发送日志
	go sendLog(ctx, taskLog)
	return
}

// 强杀任务
func KillTask(executeMsg *pbexecutor.ExecuteMsg) (err error) {
	var (
		taskExecuteInfo *TaskExecuteInfo
		exits           bool
	)
	lock.RLock()
	if taskExecuteInfo, exits = ExecutingTable[executeMsg.Task.Taskname]; !exits {
		return
	}
	lock.RUnlock()
	taskExecuteInfo.CancelFunc()
	return
}

func sendLog(ctx context.Context, taskLog *pbtasklog.TaskResLog) {
	var (
		exits     bool
		simplelog pbtasklog.SimpleLog
		err       error
	)
	// 删除执行表中的任务
	// 加锁
	lock.Lock()
	if _, exits = ExecutingTable[taskLog.Executemsg.Task.Taskname]; exits {
		delete(ExecutingTable, taskLog.Executemsg.Task.Taskname)
	}
	lock.Unlock()

	logging.Debugf("Send Task %s Log", taskLog.Executemsg.Task.Taskname)

	simplelog = pbtasklog.SimpleLog{
		Taskname:  taskLog.Executemsg.Task.Taskname,
		Command:   taskLog.Executemsg.Task.Command,
		Cronexpr:  taskLog.Executemsg.Task.Cronexpr,
		Createdby: taskLog.Executemsg.Task.Createdby,
		Timeout:   taskLog.Executemsg.Task.Timeout,
		Actuator:  taskLog.Executemsg.Task.Actuator,
		Runhost:   taskLog.Executemsg.Runhost,
		Starttime: taskLog.StartTime,
		Endtime:   taskLog.EndTime,
		Output:    taskLog.Output,
		Err:       taskLog.Err,
	}
	if _, err = tasklogclient.CreateLog(ctx, &simplelog); err != nil {
		logging.Errorf("Create Log Err:%v", err)
	}
}
