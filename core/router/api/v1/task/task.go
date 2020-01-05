package task

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
)

// CreateTask create a new worker
// POST /api/v1/task
// @params
// name
// taskType
// taskData
// parentTaskIds
// parentRunParallel
// childTaskIds
// childRunParallel
// remark
func CreateTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	//config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task := define.Task{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	_, err = json.Marshal(task.TaskData)
	if err != nil {
		log.Error("Marshal failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	if task.Name == "" {
		log.Error("task.Name is empty")
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
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
	task.ID = utils.GetID()
	task.CreateByUID = c.GetString("uid")
	task.Run = 1

	err = model.CreateTask(ctx, &task)
	if err != nil {
		log.Error("CreateTask failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron.Add(task.ID, task.Name, task.Cronexpr)
	resp.JSON(c, resp.Success, nil)
}

// ChangeTask change exist task
// Put /api/v1/task 
func ChangeTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task := define.Task{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
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

	err = model.ChangeTask(ctx, &task)
	if err != nil {
		log.Error("ChangeTask failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if task.Run == 0 {
		schedule.Cron.Del(task.ID)
	} else {
		schedule.Cron.Add(task.ID, task.Name, task.Cronexpr)
	}

	resp.JSON(c, resp.Success, nil)
}

// DeleteTask delete task
// DELETE /api/v1/task
func DeleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	taskid := c.Param("id")
	if utils.CheckID(taskid) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	exist, err := model.Check(ctx, model.TBTask, model.ID, taskid)
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
		exist, err = model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, taskid, uid)
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
	err = model.DeleteTask(ctx, taskid)
	if err != nil {
		log.Error("DeleteTask failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron.Del(taskid)
	resp.JSON(c, resp.Success, nil)

}

// GetTasks get all tasks
// GET /api/v1/tasks
func GetTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hgs, err := model.GetTasks(ctx)

	if err != nil {
		log.Error("GetTasks failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, hgs)
}

// GetTask get a task
// GET /api/v1/task/:id
func GetTask(c *gin.Context) {
	taskid := c.Param("id")
	if utils.CheckID(taskid) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	t, err := model.GetTaskByID(ctx, taskid)
	if err != nil {
		log.Error("GetTasks failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, t)
}

// RunTask start run task now
// PUT /api/v1/task/run/:id
func RunTask(c *gin.Context) {
	taskid := c.Param("id")
	if utils.CheckID(taskid) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	schedule.Cron.RunTask(taskid)
	resp.JSON(c, resp.Success, nil)
}

// KillTask kill running task
// PUT /api/v1/task/kill/:id
func KillTask(c *gin.Context) {
	log.Info("")
	taskid := c.Param("id")
	if utils.CheckID(taskid) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	schedule.Cron.KillTask(taskid)
	resp.JSON(c, resp.Success, nil)
}

// RunningTask get running task
// GET /api/v1/task/running
func RunningTask(c *gin.Context) {
	runningtasks := schedule.Cron.GetRunningtask()
	resp.JSON(c, resp.Success, runningtasks)
}

// LogTask get task log
// GET /api/v1/task/log/:id
func LogTask(c *gin.Context) {
	taskid := c.Param("id")
	if utils.CheckID(taskid) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	logs, err := model.GetLog(ctx, taskid)
	if err != nil {
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	resp.JSON(c, resp.Success, logs)
}
