package task

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorhill/cronexpr"
	"github.com/gorilla/websocket"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/middleware"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
)

// CreateTask create new task
// @Summary create new task
// @Tags Task
// @Produce json
// @Param Task body define.CreateTask true "create task"
// @Success 200 {object} resp.Response
// @Router /api/v1/task [post]
// @Security ApiKeyAuth
func CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task := model.Task{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}
	_, err = cronexpr.Parse(task.Cronexpr)
	if err != nil {
		log.Error("cronexpr.Parse failed", zap.Error(err))
		resp.JSONv2(c, define.ErrCronExpr{
			Value: task.Cronexpr,
		})
		return
	}

	id, err := model.CreateTaskv2(ctx, &task)
	if err != nil {
		resp.JSONv2(c, err)
	}
	event := schedule.EventData{
		TaskID: id,
		TE:     schedule.AddEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	//log.Debug("start Add Schedule Cron", zap.String("taskid", id))
	//schedule.Cron.Add(id, task.Name, task.Cronexpr,
	//	schedule.GetRoutePolicy(task.HostGroupID, task.RoutePolicy))
	resp.JSONv2(c, nil)
}

// ChangeTask change task
// @Summary change task
// @Tags Task
// @Produce json
// @Param Task body define.ChangeTask true "change task"
// @Success 200 {object} resp.Response
// @Router /api/v1/task [put]
// @Security ApiKeyAuth
func ChangeTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task := model.Task{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSONv2(c, err, nil)
		return
	}

	_, err = cronexpr.Parse(task.Cronexpr)
	if err != nil {
		log.Error("cronexpr.Parse failed", zap.Error(err))
		resp.JSONv2(c, define.ErrCronExpr{}, nil)
		return
	}

	// 获取用户的类型
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}

	// 这里只需要确定如果rule的用户类型是否为Admin
	if role != define.AdminUser {
		ctx = context.WithValue(ctx, "uid", c.GetString("uid"))
		// 判断ID的创建人是否为uid
		// exist, err = model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, task.ID, uid)
		// if err != nil {
		// 	log.Error("IsExist failed", zap.Error(err))
		// 	resp.JSONv2(c, resp.ErrInternalServer, nil)
		// 	return
		// }

		// if !exist {
		// 	resp.JSONv2(c, resp.ErrUnauthorized, nil)
		// 	return
		// }
	}
	err = model.ChangeTaskv2(ctx, &task)
	if err != nil {
		log.Error("change task failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}

	// err = model.ChangeTask(ctx, task.ID, task.Run, task.TaskType, task.TaskData, task.ParentTaskIds, task.ParentRunParallel,
	// 	task.ChildTaskIds, task.ChildRunParallel, task.Cronexpr, task.Timeout, task.AlarmUserIds, task.RoutePolicy,
	// 	task.ExpectCode, task.ExpectContent, task.AlarmStatus, task.HostGroupID, task.Remark,
	// )
	// if err != nil {
	// 	log.Error("ChangeTask failed", zap.Error(err))
	// 	resp.JSONv2(c, resp.ErrInternalServer, nil)
	// 	return
	// }
	// model.
	event := schedule.EventData{
		TaskID: task.ID,
		TE:     schedule.AddEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	//schedule.Cron.Add(task.ID, task.Name, task.Cronexpr,
	//	schedule.GetRoutePolicy(task.HostGroupID, task.RoutePolicy))

	resp.JSONv2(c, nil)
}

// DeleteTask delete task
// @Summary delete task
// @Tags Task
// @Produce json
// @Param Task body define.GetID true "delete task"
// @Success 200 {object} resp.Response
// @Router /api/v1/task [delete]
// @Security ApiKeyAuth
func DeleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	deletetask := define.GetID{}
	err := c.ShouldBindJSON(&deletetask)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}

	// 获取用户的类型
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}

	// 这里只需要确定如果rule的用户类型是否为Admin
	if role != define.AdminUser {
		ctx = context.WithValue(ctx, "uid", c.GetString("uid"))
		// 判断ID的创建人是否为uid
		// exist, err = model.Check(ctx, model.TBTask, model.IDCreateByUID, deletetask.ID, uid)
		// if err != nil {
		// 	log.Error("model.Check failed", zap.Error(err))
		// 	resp.JSONv2(c, resp.ErrInternalServer, nil)
		// 	return
		// }

		// if !exist {
		// 	resp.JSONv2(c, resp.ErrUnauthorized, nil)
		// 	return
		// }
	}

	err = model.DeleteTaskv2(ctx, deletetask.ID)
	if err != nil {
		log.Error("model.DeleteTask failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}

	event := schedule.EventData{
		TaskID: deletetask.ID,
		TE:     schedule.DeleteEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSONv2(c, resp.ErrInternalServer)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	resp.JSONv2(c, nil)

}

// GetTasks get tasks
// @Summary get tasks
// @Tags Task
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param psname query string false "PreSearchName"
// @Param self query bool false "Self Create Task"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task [get]
// @Security ApiKeyAuth
func GetTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	type GetQuery struct {
		define.Query
		PSName string `form:"psname"`
		Self   bool   `form:"self"`
	}
	var (
		q   GetQuery
		err error
	)

	err = c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	var createby string
	if q.Self {
		createby = c.GetString("uid")
	}

	data, count, err := model.GetTasksv2(ctx, q.Offset, q.Limit, q.PSName, createby)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}
	var (
		hgids   []string
		userids []string
	)
	for _, task := range data {
		hgids = append(hgids, task.HostgroupID)
		userids = append(userids, task.CreateUID)
	}

	var idnamerelaton = make(map[string]string)
	err = model.GetIDNameDict(ctx, hgids, &model.HostGroup{}, &idnamerelaton)
	if err != nil {
		resp.JSONv2(c, fmt.Errorf("get hostgroup id name failed: %w", err))
		return
	}
	err = model.GetIDNameDict(ctx, userids, &model.User{}, &idnamerelaton)
	if err != nil {
		resp.JSONv2(c, fmt.Errorf("get user id name failed: %w", err))
		return
	}
	resp.JSONv2(c, nil, data, count)
}

// GetTask get task info
// @Summary get tasks
// @Tags Task
// @Param ID query string true "id"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/info [get]
// @Security ApiKeyAuth
func GetTask(c *gin.Context) {
	getid := define.GetID{}
	err := c.ShouldBindQuery(&getid)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	t, err := model.GetTaskByIDv2(ctx, getid.ID)
	resp.JSONv2(c, err, t)

}

// RunTask start run task now
// @Summary get tasks
// @Tags Task
// @Param Task query define.GetID true "id"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/run [put]
// @Security ApiKeyAuth
func RunTask(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(),
	// 	config.CoreConf.Server.DB.MaxQueryTime.Duration)
	// defer cancel()

	runtask := define.GetID{}
	err := c.ShouldBindJSON(&runtask)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}
	// uid := c.GetString("uid")

	// // 获取用户的类型
	// var role define.Role
	// if v, ok := c.Get("role"); ok {
	// 	role = v.(define.Role)
	// }

	// TODO 操作日志

	event := schedule.EventData{
		TaskID: runtask.ID,
		TE:     schedule.RunEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	resp.JSONv2(c, nil)
}

// KillTask kill running task
// @Summary kill running task
// @Tags Task
// @Param Task query define.GetID true "id"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/kill [put]
// @Security ApiKeyAuth
func KillTask(c *gin.Context) {
	runtask := define.GetID{}
	err := c.ShouldBindJSON(&runtask)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}

	event := schedule.EventData{
		TaskID: runtask.ID,
		TE:     schedule.KillEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	//schedule.Cron.KillTask(runtask.ID)
	resp.JSONv2(c, nil)
}

// GetRunningTask return running task
// @Summary get tasks
// @Tags Task
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/running [get]
// @Security ApiKeyAuth
func GetRunningTask(c *gin.Context) {
	var (
		q            define.Query
		err          error
		runningtasks []*define.RunTask
	)

	err = c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	allrunningtasks, err := schedule.Cron2.GetRunningTask()
	if err != nil {
		resp.JSONv2(c, err)
	}
	if len(runningtasks) < q.Offset {
		runningtasks = []*define.RunTask{}
	} else if len(allrunningtasks) >= q.Offset && len(allrunningtasks) < q.Offset+q.Limit {
		runningtasks = allrunningtasks[q.Offset:]
	} else {
		runningtasks = allrunningtasks[q.Offset : q.Offset+q.Limit]
	}

	resp.JSONv2(c, nil, runningtasks, len(runningtasks))
}

// LogTask get task log
// @Summary get tasks
// @Tags Task
// @Param taskname query int false "taskName"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Param status query int false "Status"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/log [get]
// @Security ApiKeyAuth
func LogTask(c *gin.Context) {
	type Log struct {
		Name   string `form:"name"`
		Status int    `form:"status" binding:"gte=-1,lte=1"`
		define.Query
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	var (
		q Log
	)

	err := c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
		resp.JSONv2(c, err)
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	logs, count, err := model.GetLogv2(ctx, q.Name, q.Status, q.Offset, q.Limit)
	resp.JSONv2(c, err, logs, count)
}

// LogTreeData get log tree
// @Summary get tasks log tree data
// @Tags Task
// @Param id query int false "ID"
// @Param start_time query int false "StartTime"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/log/tree [get]
// @Security ApiKeyAuth
// func LogTreeData(c *gin.Context) {
// 	// TODO
// 	getid := define.GetID{}
// 	err := c.BindQuery(&getid)
// 	if err != nil {
// 		log.Error("c.BindQuery", zap.Error(err))
// 		resp.JSONv2(c, err)
// 		return
// 	}

// 	starttime := c.Query("start_time")
// 	ctx, cancel := context.WithTimeout(context.Background(),
// 		config.CoreConf.Server.DB.MaxQueryTime.Duration)
// 	defer cancel()
// 	if starttime == "" {
// 		log.Error("can't get start_time")
// 		resp.JSONv2(c, err)
// 		return
// 	}
// 	starttimeint, err := strconv.ParseInt(starttime, 10, 64)
// 	if err != nil {
// 		log.Error("strconv.ParseInt", zap.Error(err))
// 		resp.JSONv2(c, resp.ErrBadRequest, nil)
// 		return
// 	}
// 	TaskTreeStatus, err := model.GetTreeLog(ctx, getid.ID, starttimeint)
// 	if err != nil {
// 		log.Error("model.GetTreeLog", zap.Error(err))
// 		resp.JSONv2(c, resp.ErrInternalServer, nil)
// 		return
// 	}
// 	resp.JSONv2(c, resp.Success, TaskTreeStatus)
// }

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	defaultSendTTL = 2 * time.Second
)

// RealRunTaskLog return real time log
// GET /api/v1/task/log/websocket?id=manid&realid=ididididid&type=
func RealRunTaskLog(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("Upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()
	type tasklog struct {
		define.GetID
		RealID string              `form:"realid" binding:"required"`
		Type   define.TaskRespType `form:"type" binding:"required"`
	}
	query := tasklog{}
	err = c.BindQuery(&query)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}
	// realid := c.Query("realid")
	// taskruntype, err := strconv.Atoi(c.Query("type"))
	// if err != nil {
	// 	log.Error("can get valid task type", zap.Error(err))
	// 	conn.WriteMessage(websocket.TextMessage, []byte("can get task type"))
	// 	return
	// }

	task, ok := schedule.Cron2.GetTask(query.ID)
	if !ok {
		log.Error("can get taskid", zap.String("taskid", query.ID))
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("can get taskid %s", query.ID)))
		return
	}
	var offset int64
	for {
		output, err := task.GetTaskRealLog(define.TaskRespType(query.Type), query.RealID, offset)
		if err == nil {
			offset++
			err = conn.WriteMessage(websocket.TextMessage, output)
			if err != nil {
				log.Error("WriteMessage failed", zap.Error(err))
				return
			}
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Error("ReadMessage failed", zap.Error(err))
				return
			}
			time.Sleep(time.Millisecond * 10)
			continue
		}
		if errors.Is(err, io.EOF) {
			log.Debug("read task log over")
			// conn.WriteMessage(websocket.TextMessage, []byte("task run finished"))
			return
		} else if errors.Is(err, schedule.ErrNoGetLog) {
			log.Debug("can not get new data, please wait some time")
			// if can get data,check task is running ,is task is stop then close websocket
			ok, err := schedule.Cron2.IsRunning(query.ID)
			if err != nil {
				log.Error("Cron2.IsRunning failed", zap.Error(err))
				return
			}
			if !ok {
				log.Warn("task is not running ", zap.String("taskid", query.ID))
				return
			}
			time.Sleep(time.Second)
		} else {
			var erroutput []byte
			if errors.Is(err, redis.Nil) {
				erroutput = []byte("task is run finished")
			} else {
				log.Error("read task log failed", zap.Error(err))
				erroutput = []byte(err.Error())
			}
			err = conn.WriteMessage(websocket.TextMessage, erroutput)
			if err != nil {
				log.Error("WriteMessage failed", zap.Error(err))

			}
			return
		}
	}
}

// RealRunTaskStatus  Get Task Status
// GET /api/v1/task/status/ws?id=manid
func RealRunTaskStatus(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("Upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	getid := define.GetID{}
	err = c.BindQuery(&getid)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(define.ErrBadRequest{}.Error()))
		return
	}

	log.Debug("start get real task status", zap.String("taskid", getid.ID))

	_, token, err := conn.ReadMessage()
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("get token fail"))
		return
	}
	_, _, pass := middleware.CheckToken(string(token))
	if !pass {
		conn.WriteMessage(websocket.TextMessage, []byte("check token auth fail"))
		return
	}
	task, ok := schedule.Cron2.GetTask(getid.ID)
	if !ok {
		log.Error("can not get task", zap.String("taskid", getid.ID))
		return
	}
	timer := time.NewTimer(time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			taskrunstatus, finish, err := task.GetTaskTreeStatatus()
			if err != nil {
				log.Error("task.GetTaskTreeStatatus failed", zap.Error(err))
				return
			}

			err = conn.WriteJSON(taskrunstatus)
			if err != nil {
				log.Error("WriteJSON failed", zap.Error(err))
				return
			}
			// if task status has one of  running,wait,so return status
			// otherwise close websocket
			if finish {
				return
			}
			_, _, err = conn.ReadMessage()
			if err != nil {
				log.Error("ReadMessage failed", zap.Error(err))
				return
			}
			timer.Reset(defaultSendTTL)
		}
	}
}

// ParseCron parse cronexpr
// @Summary parse cronexpr
// @Tags Task
// @Param expr query string true "Expr"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/cron [get]
// @Security ApiKeyAuth
func ParseCron(c *gin.Context) {
	type reqexpr struct {
		CronExpr string `form:"expr" binding:"required"`
	}
	reqep := reqexpr{}
	err := c.ShouldBindQuery(&reqep)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}
	cronbyte, err := base64.StdEncoding.DecodeString(reqep.CronExpr)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}

	expr, err := cronexpr.Parse(string(cronbyte))
	if err != nil {
		resp.JSONv2(c, define.ErrCronExpr{
			Value: string(cronbyte),
		})
		return
	}
	var (
		last time.Time
		next time.Time
	)
	last = time.Now()
	resptimes := []string{}
	for {
		next = expr.Next(last)
		last = next
		resptimes = append(resptimes, next.Format("2006-01-02 15:04:05"))
		if len(resptimes) == 10 {
			break
		}
	}
	resp.JSONv2(c, nil, resptimes)
}

// GetSelect name,id
// @Summary Get Task Select
// @Tags Task
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/select [get]
// @Security ApiKeyAuth
func GetSelect(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	data, err := model.GetIDNameOption(ctx, nil, &model.Task{})
	resp.JSONv2(c, err, data)
}

// CloneTask clone task
// @Summary create a task by copy old task
// @Tags Task
// @Param Task body define.IDName true "clone task"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/clone [post]
// @Security ApiKeyAuth
func CloneTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	clonetask := define.IDName{}

	err := c.ShouldBindJSON(&clonetask)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}

	task, err := model.GetTaskByIDv2(ctx, clonetask.ID)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}
	id, err := model.CreateTaskv2(ctx, task)
	if err != nil {
		resp.JSONv2(c, err)
		return
	}
	event := schedule.EventData{
		TaskID: id,
		TE:     schedule.AddEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	resp.JSONv2(c, nil)
}

// CleanTaskLog clean old task log
// @Summary create a task by copy old task
// @Tags Task
// @Param Log body define.Cleanlog true "clean task log"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/clone [delete]
// @Security ApiKeyAuth
func CleanTaskLog(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	cleanlog := define.Cleanlog{}

	err := c.ShouldBindJSON(&cleanlog)
	if err != nil {
		log.Error("c.ShouldBindJson failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}

	deletetime := time.Now().Add(time.Hour * time.Duration(-24*cleanlog.PreDay))
	delcount, err := model.CleanLogv2(ctx, cleanlog.ID, deletetime)
	resp.JSONv2(c,err,nil,delcount)
}
