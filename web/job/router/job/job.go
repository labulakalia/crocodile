package job

import (
	"context"
	"crocodile/common/bind"
	"crocodile/common/cfg"
	"crocodile/common/e"
	"crocodile/common/registry"
	"crocodile/common/response"

	"crocodile/common/wrapper"

	pbjob "crocodile/service/job/proto/job"
	pbtasklog "crocodile/service/tasklog/proto/tasklog"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/client"
	"time"
)

var (
	JobClient pbjob.JobService
	Logclient pbtasklog.TaskLogService
)

func Init() {
	c := client.NewClient(
		client.Retries(3),
		client.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
	)
	JobClient = pbjob.NewJobService("crocodile.srv.job", client.DefaultClient)
	Logclient = pbtasklog.NewTaskLogService("crocodile.srv.tasklog", c)
}

type QueryTask struct {
	TaskName string `form:"taskname" json:"taskname" validate:"required"`
}

func CreateJob(c *gin.Context) {
	var (
		app       response.Gin
		ctx       context.Context
		task      pbjob.Task
		resp      *pbjob.Response
		err       error
		code      int32
		loginuser string
		exists    bool
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	task = pbjob.Task{}

	if err = bind.BindJson(c, &task); err != nil {
		logging.Errorf("BindJson Err:%v", err)
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	if loginuser, exists = c.Keys["user"].(string); !exists {
		code = e.ERR_TOKEN_INVALID
		app.Response(code, nil)
		return
	}
	task.Createdby = loginuser

	resp, err = JobClient.CreateJob(ctx, &task)
	if err != nil {
		code = e.ERR_CREATE_JOB_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)

}
func DeleteJob(c *gin.Context) {
	var (
		app response.Gin
		ctx context.Context

		task   pbjob.Task
		resp   *pbjob.Response
		err    error
		code   int32
		exists bool
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}

	task, exists = c.Keys["data"].(pbjob.Task)
	if !exists {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	resp, err = JobClient.DeleteJob(ctx, &task)
	if err != nil {
		code = e.ERR_DELETE_JOB_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)

}
func ChangeJob(c *gin.Context) {
	var (
		app    response.Gin
		ctx    context.Context
		task   pbjob.Task
		resp   *pbjob.Response
		err    error
		code   int32
		exists bool
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	logging.Infof("%+v", c.Keys["data"])
	task, exists = c.Keys["data"].(pbjob.Task)
	if !exists {
		logging.Error("Not Exits data")
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	resp, err = JobClient.ChangeJob(ctx, &task)
	if err != nil {
		code = e.ERR_CHANGE_ACTUAT_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)

}

func GetJob(c *gin.Context) {
	var (
		app       response.Gin
		ctx       context.Context
		querytask QueryTask
		reqtask   *pbjob.Task
		resp      *pbjob.Response
		err       error
		code      int32
		res       []*pbjob.Task
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}

	querytask = QueryTask{}

	_ = bind.BindQuery(c, &querytask)

	reqtask = &pbjob.Task{
		Taskname: querytask.TaskName,
	}
	resp, err = JobClient.GetJob(ctx, reqtask)
	if err != nil {
		logging.Errorf("Get Job Err: %v", err)
		code = e.ERR_GET_JOB_FAIL
		app.Response(code, nil)
		return
	}
	res = []*pbjob.Task{}
	for _, task := range resp.Tasks {
		res = append(res, task)
	}
	app.Response(resp.Code, res)
}

func RunJob(c *gin.Context) {
	var (
		app  response.Gin
		ctx  context.Context
		task pbjob.Task

		resp   *pbjob.Response
		err    error
		code   int32
		exists bool
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}

	task, exists = c.Keys["data"].(pbjob.Task)
	if !exists {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	resp, err = JobClient.RunJob(ctx, &task)
	if err != nil {
		code = e.ERR_RUN_JOB_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)
}

func KillJob(c *gin.Context) {
	var (
		app  response.Gin
		ctx  context.Context
		task pbjob.Task

		resp   *pbjob.Response
		err    error
		code   int32
		exists bool
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}

	task, exists = c.Keys["data"].(pbjob.Task)
	if !exists {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	resp, err = JobClient.KillJob(ctx, &task)
	if err != nil {
		code = e.ERR_KILL_JOB_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)
}

type QueryLog struct {
	Taskname string    `form:"taskname" validate:"required"`
	Fromtime time.Time `form:"fromtime" validate:"required"`
	Totime   time.Time `form:"totime" validate:"required"`
	Limit    int32     `form:"limit" validate:"required"`
	Offset   int32     `form:"offset"`
}

type GetLog struct {
	Id        uint64 `json:"id"`
	Taskname  string `json:"taskname"`
	Command   string `json:"command"`
	Cronexpr  string `json:"cronexpr"`
	Createdby string `json:"createdby"`
	Timeout   int32  `json:"timeout"`
	Actuator  string `json:"actuator"`
	Runhost   string `json:"runhost"`
	Starttime string `json:"starttime"`
	Endtime   string `json:"endtime"`
	Output    string `json:"output"`
	Err       string `json:"err"`
}

type RespLog struct {
	Logs  []*GetLog `json:"logs"`
	Count int32     `json:"count"`
}

func GetJobLog(c *gin.Context) {
	var (
		app        response.Gin
		ctx        context.Context
		querylog   QueryLog
		reqtasklog pbtasklog.QueryLog
		respLog    *pbtasklog.RespLog
		fromtime   *timestamp.Timestamp
		totime     *timestamp.Timestamp
		err        error
		code       int32
		getlogs    []*GetLog
	)
	ctx, ok := wrapper.ContextWithSpan(c)
	if ok == false {
		logging.Error("get context err")
		ctx = context.Background()
	}
	ctx, _ = context.WithTimeout(ctx, time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	querylog = QueryLog{}
	if err = bind.BindQuery(c, &querylog); err != nil {
		logging.Errorf("Bind Query Err:%v", err)
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	fromtime, err = ptypes.TimestampProto(querylog.Fromtime)
	if err != nil {
		logging.Errorf("Bind Query Err:%v", err)
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}
	totime, err = ptypes.TimestampProto(querylog.Totime)
	if err != nil {
		logging.Errorf("Bind Query Err:%v", err)
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	reqtasklog = pbtasklog.QueryLog{
		Taskname: querylog.Taskname,
		Fromtime: fromtime,
		Totime:   totime,
		Offset:   querylog.Offset,
		Limit:    querylog.Limit,
	}
	if respLog, err = Logclient.GetLog(ctx, &reqtasklog); err != nil {
		code = e.ERR_GET_JOB_LOG_FAIL
		app.Response(code, nil)
		return
	}
	for _, log := range respLog.Logs {
		getlog := GetLog{}
		getlog.Id = log.Id
		getlog.Taskname = log.Taskname
		getlog.Command = log.Command
		getlog.Cronexpr = log.Cronexpr
		getlog.Createdby = log.Createdby
		getlog.Timeout = log.Timeout
		getlog.Actuator = log.Actuator
		getlog.Runhost = log.Runhost
		getlog.Starttime = ptypes.TimestampString(log.Starttime)
		getlog.Endtime = ptypes.TimestampString(log.Endtime)
		getlog.Output = log.Output
		getlog.Err = log.Err
		getlogs = append(getlogs, &getlog)
	}

	code = e.SUCCESS
	app.Response(code, RespLog{Count: respLog.Count, Logs: getlogs})
}
