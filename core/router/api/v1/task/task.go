package task

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorhill/cronexpr"
	"github.com/gorilla/websocket"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
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
	//config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task := define.CreateTask{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	_, err = cronexpr.Parse(task.Cronexpr)
	if err != nil {
		log.Error("cronexpr.Parse failed", zap.Error(err))
		resp.JSON(c, resp.ErrCronExpr, nil)
		return
	}

	// TODO 检查任务数据
	exist, err := model.Check(ctx, model.TBTask, model.Name, task.Name)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrTaskExist, nil)
		return
	}
	// task.CreateByUID = c.GetString("uid")
	task.Run = true
	id := utils.GetID()
	err = model.CreateTask(ctx, id, task.Name, task.TaskType, task.TaskData, true, task.ParentTaskIds, task.ParentRunParallel,
		task.ChildTaskIds, task.ChildRunParallel, task.Cronexpr, task.Timeout, task.AlarmUserIds, task.RoutePolicy,
		task.ExpectCode, task.ExpectContent, task.AlarmStatus, c.GetString("uid"), task.HostGroupID, task.Remark,
	)
	if err != nil {
		log.Error("CreateTask failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	event := schedule.EventData{
		TaskID: id,
		TE:     schedule.AddEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	//log.Debug("start Add Schedule Cron", zap.String("taskid", id))
	//schedule.Cron.Add(id, task.Name, task.Cronexpr,
	//	schedule.GetRoutePolicy(task.HostGroupID, task.RoutePolicy))
	resp.JSON(c, resp.Success, nil)
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

	task := define.ChangeTask{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	_, err = cronexpr.Parse(task.Cronexpr)
	if err != nil {
		log.Error("cronexpr.Parse failed", zap.Error(err))
		resp.JSON(c, resp.ErrCronExpr, nil)
		return
	}

	exist, err := model.Check(ctx, model.TBTask, model.ID, task.ID)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	if !exist {
		resp.JSON(c, resp.ErrTaskNotExist, nil)
		return
	}

	uid := c.GetString("uid")

	// 获取用户的类型
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}

	// 这里只需要确定如果rule的用户类型是否为Admin
	if role != define.AdminUser {
		// 判断ID的创建人是否为uid
		exist, err = model.Check(ctx, model.TBTask, model.IDCreateByUID, task.ID, uid)
		if err != nil {
			log.Error("IsExist failed", zap.Error(err))
			resp.JSON(c, resp.ErrInternalServer, nil)
			return
		}

		if !exist {
			resp.JSON(c, resp.ErrUnauthorized, nil)
			return
		}
	}

	err = model.ChangeTask(ctx, task.ID, task.Run, task.TaskType, task.TaskData, task.ParentTaskIds, task.ParentRunParallel,
		task.ChildTaskIds, task.ChildRunParallel, task.Cronexpr, task.Timeout, task.AlarmUserIds, task.RoutePolicy,
		task.ExpectCode, task.ExpectContent, task.AlarmStatus, task.HostGroupID, task.Remark,
	)
	if err != nil {
		log.Error("ChangeTask failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	event := schedule.EventData{
		TaskID: task.ID,
		TE:     schedule.AddEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	//schedule.Cron.Add(task.ID, task.Name, task.Cronexpr,
	//	schedule.GetRoutePolicy(task.HostGroupID, task.RoutePolicy))

	resp.JSON(c, resp.Success, nil)
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(deletetask.ID) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	exist, err := model.Check(ctx, model.TBTask, model.ID, deletetask.ID)
	if err != nil {
		log.Error("model.Check failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	if !exist {
		log.Warn("unauthorized ", zap.String("id", deletetask.ID))
		resp.JSON(c, resp.ErrHostgroupNotExist, nil)
		return
	}

	uid := c.GetString("uid")

	// 获取用户的类型
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}

	// 这里只需要确定如果rule的用户类型是否为Admin
	if role != define.AdminUser {
		// 判断ID的创建人是否为uid
		exist, err = model.Check(ctx, model.TBTask, model.IDCreateByUID, deletetask.ID, uid)
		if err != nil {
			log.Error("model.Check failed", zap.Error(err))
			resp.JSON(c, resp.ErrInternalServer, nil)
			return
		}

		if !exist {
			resp.JSON(c, resp.ErrUnauthorized, nil)
			return
		}
	}

	usecount, err := model.TaskIsUse(ctx, deletetask.ID)
	if err != nil {
		log.Error("model.TaskIsUse failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if usecount > 0 {
		log.Warn("task can delete,use by other task", zap.String("taskid", deletetask.ID), zap.Int("use count", usecount))
		resp.JSON(c, resp.ErrTaskUseByOtherTask, nil)
		return
	}

	err = model.DeleteTask(ctx, deletetask.ID)
	if err != nil {
		log.Error("model.DeleteTask failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	//schedule.Cron.Del(deletetask.ID)
	//_, err = model.	CleanTaskLog(ctx, "", deletetask.ID, time.Now().UnixNano()/1e6)
	//if err != nil {
	//	log.Error("model.CleanTaskLog failed", zap.Error(err))
	//	resp.JSON(c, resp.ErrInternalServer, nil)
	//	return
	//}
	event := schedule.EventData{
		TaskID: deletetask.ID,
		TE:     schedule.DeleteEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	resp.JSON(c, resp.Success, nil)

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
	hgs, count, err := model.GetTasks(ctx, q.Offset, q.Limit, "", q.PSName, createby)
	if err != nil {
		log.Error("GetTasks failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, hgs, count)
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(getid.ID) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	t, err := model.GetTaskByID(ctx, getid.ID)

	switch err.(type) {
	case nil:
		resp.JSON(c, resp.Success, t)
	case define.ErrNotExist:
		resp.JSON(c, resp.ErrTaskNotExist, nil)
	default:
		log.Error("GetTasks failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
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
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	runtask := define.GetID{}
	err := c.ShouldBindJSON(&runtask)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(runtask.ID) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	uid := c.GetString("uid")

	// 获取用户的类型
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}

	// 这里只需要确定如果rule的用户类型是否为Admin
	if role != define.AdminUser {
		// 判断ID的创建人是否为uid
		exist, err := model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, runtask.ID, uid)
		if err != nil {
			log.Error("model.Check failed", zap.Error(err))
			resp.JSON(c, resp.ErrInternalServer, nil)
			return
		}

		if !exist {
			resp.JSON(c, resp.ErrUnauthorized, nil)
			return
		}
	}
	//go schedule.Cron.RunTask(runtask.ID, define.Manual)

	event := schedule.EventData{
		TaskID: runtask.ID,
		TE:     schedule.RunEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	resp.JSON(c, resp.Success, nil)
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(runtask.ID) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	event := schedule.EventData{
		TaskID: runtask.ID,
		TE:     schedule.KillEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	//schedule.Cron.KillTask(runtask.ID)
	resp.JSON(c, resp.Success, nil)
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
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	allrunningtasks, err := schedule.Cron2.GetRunningTask()
	if err != nil {
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if len(runningtasks) < q.Offset {
		runningtasks = []*define.RunTask{}
	} else if len(allrunningtasks) >= q.Offset && len(allrunningtasks) < q.Offset+q.Limit {
		runningtasks = allrunningtasks[q.Offset:]
	} else {
		runningtasks = allrunningtasks[q.Offset : q.Offset+q.Limit]
	}

	resp.JSON(c, resp.Success, runningtasks, len(runningtasks))
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

	name := c.Query("name")
	statusstr := c.Query("status")

	status, err := strconv.Atoi(statusstr)
	if err != nil {
		log.Warn("get params status is not int", zap.Error(err))
	}
	if status < -1 || status > 1 {
		status = 0
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	var (
		q define.Query
	)

	err = c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	logs, count, err := model.GetLog(ctx, name, status, q.Offset, q.Limit)
	if err != nil {
		log.Error("GetLog failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, logs, count)
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
func LogTreeData(c *gin.Context) {
	getid := define.GetID{}
	err := c.BindQuery(&getid)
	if err != nil {
		log.Error("c.BindQuery", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	starttime := c.Query("start_time")
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	if starttime == "" {
		log.Error("can't get start_time")
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	starttimeint, err := strconv.ParseInt(starttime, 10, 64)
	if err != nil {
		log.Error("strconv.ParseInt", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	TaskTreeStatus, err := model.GetTreeLog(ctx, getid.ID, starttimeint)
	if err != nil {
		log.Error("model.GetTreeLog", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, TaskTreeStatus)
}

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
	getid := define.GetID{}
	err = c.BindQuery(&getid)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	realid := c.Query("realid")
	taskruntype, err := strconv.Atoi(c.Query("type"))
	if err != nil {
		log.Error("can get valid task type", zap.Error(err))
		conn.WriteMessage(websocket.TextMessage, []byte("can get task type"))
		return
	}

	task, ok := schedule.Cron2.GetTask(getid.ID)
	if !ok {
		log.Error("can get taskid", zap.String("taskid", getid.ID))
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("can get taskid %s", getid.ID)))
		return
	}
	var offset int64
	for {
		output, err := task.GetTaskRealLog(define.TaskRespType(taskruntype), realid, offset)
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
			ok, err := schedule.Cron2.IsRunning(getid.ID)
			if err != nil {
				log.Error("Cron2.IsRunning failed", zap.Error(err))
				return
			}
			if !ok {
				log.Warn("task is not running ", zap.String("taskid", getid.ID))
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
		resp.JSON(c, resp.ErrBadRequest, nil)
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	cronbyte, err := base64.StdEncoding.DecodeString(reqep.CronExpr)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	var respTimes []string
	nextN := cronexpr.MustParse(string(cronbyte)).NextN(time.Now(), 10)
	for _, nextTime := range nextN {
		respTimes = append(respTimes, nextTime.Format("2006-01-02 15:04:05"))
	}
	resp.JSON(c, resp.Success, respTimes)

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
	data, err := model.GetNameID(ctx, model.TBTask)
	if err != nil {
		log.Error("model.GetNameID failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, data)
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	exist, err := model.Check(ctx, model.TBTask, model.Name, clonetask.Name)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrTaskExist, nil)
		return
	}

	task, err := model.GetTaskByID(ctx, clonetask.ID)
	switch err.(type) {
	case nil:
		goto Next
	case define.ErrNotExist:
		log.Error("get task failed", zap.Error(err))
		resp.JSON(c, resp.ErrTaskNotExist, nil)
		return
	default:
		log.Error("model.GetTaskByID failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
Next:
	id := utils.GetID()
	if id == "" {
		log.Error("utils.GetID return empty", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	err = model.CreateTask(ctx,
		id,
		clonetask.Name,
		task.TaskType,
		task.TaskData,
		task.Run,
		task.ParentTaskIds,
		task.ParentRunParallel,
		task.ChildTaskIds,
		task.ChildRunParallel,
		task.Cronexpr,
		task.Timeout,
		task.AlarmUserIds,
		task.RoutePolicy,
		task.ExpectCode,
		task.ExpectContent,
		task.AlarmStatus,
		c.GetString("uid"),
		task.HostGroupID,
		fmt.Sprintf("从任务%s克隆", task.Name))
	if err != nil {
		log.Error(" model.CreateTask failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	event := schedule.EventData{
		TaskID: id,
		TE:     schedule.AddEvent,
	}
	res, err := json.Marshal(event)
	if err != nil {
		log.Error("json.Marshal failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron2.PubTaskEvent(res)
	resp.JSON(c, resp.Success, nil)
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	// TODO 检查任务数据
	exist, err := model.Check(ctx, model.TBTask, model.Name, cleanlog.Name)
	if err != nil {
		log.Error("IsExist failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if !exist {
		log.Error("task is not exist", zap.String("name", cleanlog.Name))
		resp.JSON(c, resp.ErrTaskNotExist, nil)
		return
	}

	// 获取用户的类型
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}

	// 这里只需要确定如果rule的用户类型是否为Admin
	if role != define.AdminUser {
		// 判断任务的创建人是否为当前用户
		exist, err = model.Check(ctx, model.TBTask, model.NameCreateByUID, cleanlog.Name, c.GetString("uid"))
		if err != nil {
			log.Error("IsExist failed", zap.Error(err))
			resp.JSON(c, resp.ErrInternalServer, nil)
			return
		}

		if !exist {
			resp.JSON(c, resp.ErrUnauthorized, nil)
			return
		}
	}

	deletetime := (time.Now().UnixNano() - int64(time.Hour)*24*cleanlog.PreDay) / 1e6
	delcount, err := model.CleanTaskLog(ctx, cleanlog.Name, "", deletetime)
	if err != nil {
		log.Error("model.CleanTaskLog failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	type del struct {
		DelCount int64 `json:"delcount"`
	}

	resp.JSON(c, resp.Success, del{DelCount: delcount})
}
