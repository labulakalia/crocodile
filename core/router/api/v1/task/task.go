package task

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
		resp.JSON(c, resp.ErrBadRequest, nil)
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
	err = model.CreateTask(ctx, id, task.Name, task.TaskType, task.TaskData, task.ParentTaskIds, task.ParentRunParallel,
		task.ChildTaskIds, task.ChildRunParallel, task.Cronexpr, task.Timeout, task.AlarmUserIds, task.RoutePolicy,
		task.ExpectCode, task.ExpectContent, task.AlarmStatus, c.GetString("uid"), task.HostGroupID, task.Remark,
	)
	if err != nil {
		log.Error("CreateTask failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	log.Debug("start Add Schedule Cron", zap.String("taskid", id))
	schedule.Cron.Add(id, task.Name, task.Cronexpr,
		schedule.GetRoutePolicy(task.HostGroupID, task.RoutePolicy))
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
		resp.JSON(c, resp.ErrBadRequest, nil)
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
		exist, err = model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, task.ID, uid)
		if err != nil {
			log.Error("IsExist failed", zap.String("error", err.Error()))
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
		log.Error("ChangeTask failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron.Add(task.ID, task.Name, task.Cronexpr,
		schedule.GetRoutePolicy(task.HostGroupID, task.RoutePolicy))

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
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	if !exist {
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
		exist, err = model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, deletetask.ID, uid)
		if err != nil {
			log.Error("IsExist failed", zap.String("error", err.Error()))
			resp.JSON(c, resp.ErrInternalServer, nil)
			return
		}

		if !exist {
			resp.JSON(c, resp.ErrUnauthorized, nil)
			return
		}
	}
	err = model.DeleteTask(ctx, deletetask.ID)
	if err != nil {
		log.Error("DeleteTask failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron.Del(deletetask.ID)
	resp.JSON(c, resp.Success, nil)

}

// GetTasks get tasks
// @Summary get tasks
// @Tags Task
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task [get]
// @Security ApiKeyAuth
func GetTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	var (
		q   define.Query
		err error
	)

	err = c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	hgs, err := model.GetTasks(ctx, q.Offset, q.Limit)

	if err != nil {
		log.Error("GetTasks failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, hgs)
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
	if err != nil {
		log.Error("GetTasks failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, t)
}

// RunTask start run task now
// GetTask get task info
// @Summary get tasks
// @Tags Task
// @Param Task query define.GetID true "id"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/run [put]
// @Security ApiKeyAuth
func RunTask(c *gin.Context) {
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
	go schedule.Cron.RunTask(runtask.ID, define.Manual)
	resp.JSON(c, resp.Success, nil)
}

// KillTask kill running task
// GetTask kill running task
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
	schedule.Cron.KillTask(runtask.ID)
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
		q   define.Query
		err error
	)

	err = c.BindQuery(&q)
	if err != nil {
		log.Error("BindQuery offset failed", zap.Error(err))
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}
	runningtasks := schedule.Cron.GetRunningtask()
	if len(runningtasks) < q.Offset {
		runningtasks = []*define.RunTask{}
	} else if len(runningtasks) >= q.Offset && len(runningtasks) < q.Offset+q.Limit {
		runningtasks = runningtasks[q.Offset:len(runningtasks)]
	} else {
		runningtasks = runningtasks[q.Offset : q.Offset+q.Limit]
	}

	resp.JSON(c, resp.Success, runningtasks)
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
	getname := define.GetName{}
	err := c.BindQuery(&getname)
	if err != nil {
		log.Error("BindQuery", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

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
	logs, err := model.GetLog(ctx, getname.Name, status, q.Offset, q.Limit)
	if err != nil {
		log.Error("GetLog failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, logs)
}

// LogTreeData get log tree
// @Summary get tasks log tree data
// @Tags Task
// @Param id query int false "ID"
// @Param start_time query int false "StartTime"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/task/log/tree[get]
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
	timeout        = 5 * time.Second
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
	tasktype, err := strconv.Atoi(c.Query("type"))
	if err != nil {
		log.Error("can get valid task type", zap.Error(err))
		conn.WriteMessage(websocket.TextMessage, []byte("can get task type"))
		return
	}

	logcache, err := schedule.Cron.GetRunTaskLogCache(getid.ID, realid, define.TaskRespType(tasktype))
	if err != nil {
		log.Error("GetRunTaskLogCache failed", zap.Error(err))
		err = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	offset := 0

	var out = make([]byte, 1024)
	for {
		n, err := logcache.ReadOnly(out, offset)
		if err == nil {
			if n > 0 {
				offset += n
				err = conn.WriteMessage(websocket.TextMessage, out[:n])
				if err != nil {
					log.Error("WriteMessage failed", zap.Error(err))
					return
				}
				_, _, err := conn.ReadMessage()
				if err != nil {
					log.Error("ReadMessage failed", zap.Error(err))
					return
				}
			}
			time.Sleep(time.Second)
			continue
		}
		if err == io.EOF {
			log.Debug("read task log over")
			// conn.WriteMessage(websocket.TextMessage, []byte("task run finished"))
			return
		} else if err == schedule.ErrNoReadData {
			log.Debug("can not get new data, please wait some time")
			time.Sleep(time.Second)
		} else {
			log.Error("read task log failed", zap.Error(err))
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

	timer := time.NewTimer(time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			taskrunstatus := schedule.Cron.GetRunTaskStaus(getid.ID)
			if taskrunstatus == nil {
				// conn.WriteMessage(websocket.TextMessage, []byte("task run finish"))
				log.Error("GetRunTaskStaus failed")
				return
			}

			err := conn.WriteJSON(taskrunstatus)
			if err != nil {
				log.Error("WriteJSON failed", zap.Error(err))
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
	cronbase64 := c.Query("expr")
	cronbyte, err := base64.StdEncoding.DecodeString(cronbase64)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	expr, err := cronexpr.Parse(string(cronbyte))
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
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
	resp.JSON(c, resp.Success, resptimes)
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
		log.Error("model.GetNameID", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, data)
}
