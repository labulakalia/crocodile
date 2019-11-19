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

// POST /api/v1/task
// @params
// name
// taskType
// taskData
// parentTaskIds
// parentRunParallel
// childTaskIds
// childRunParallel
// execPlanID
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
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}

	_, err = json.Marshal(task.TaskData)
	if err != nil {
		log.Error("Marshal failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}

	if task.Name == "" {
		log.Error("task.Name is empty")
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}
	exist, err := model.Check(ctx, model.TB_task, model.Name, task.Name)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.Json(c, resp.ErrTaskExist, nil)
		return
	}
	task.Id = utils.GetId()
	task.CreateByUId = c.GetString("uid")
	task.Run = 1

	err = model.CreateTask(ctx, &task)
	if err != nil {
		log.Error("CreateTask failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron.Add(task.Id, task.CronExpr)
	resp.Json(c, resp.Success, nil)
}

// PATCH /api/v1/task
func ChangeTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task := define.Task{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}
	exist, err := model.Check(ctx, model.TB_task, model.ID, task.Id)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}

	if !exist {
		resp.Json(c, resp.ErrHostgroupNotExist, nil)
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
		exist, err = model.Check(ctx, model.TB_hostgroup, model.IDCreateByUID, task.Id, uid)
		if err != nil {
			log.Error("IsExist failed", zap.String("error", err.Error()))
			resp.Json(c, resp.ErrInternalServer, nil)
			return
		}

		if !exist {
			resp.Json(c, resp.ErrUnauthorized, nil)
			return
		}
	}

	err = model.ChangeTask(ctx, &task)
	if err != nil {
		log.Error("ChangeTask failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron.Add(task.Id, task.CronExpr)
	resp.Json(c, resp.Success, nil)
}

// DELETE /api/v1/task
func DeleteTask(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	task := define.Task{}
	err := c.ShouldBindJSON(&task)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}
	exist, err := model.Check(ctx, model.TB_task, model.ID, task.Id)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}

	if !exist {
		resp.Json(c, resp.ErrHostgroupNotExist, nil)
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
		exist, err = model.Check(ctx, model.TB_hostgroup, model.IDCreateByUID, task.Id, uid)
		if err != nil {
			log.Error("IsExist failed", zap.String("error", err.Error()))
			resp.Json(c, resp.ErrInternalServer, nil)
			return
		}

		if !exist {
			resp.Json(c, resp.ErrUnauthorized, nil)
			return
		}
	}
	err = model.DeleteTask(ctx, task.Id)
	if err != nil {
		log.Error("DeleteTask failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	schedule.Cron.Del(task.Id)
	resp.Json(c, resp.Success, nil)

}

// GET /api/v1/tasks
func GetTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hgs, err := model.GetTasks(ctx)
	if err != nil {
		log.Error("GetTasks failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, hgs)

}

// 立即运行
// GET /api/v1/task/run
func RunTask(c *gin.Context) {

}

// 查看任务日志
// GET /api/v1/task/logs
func LogsTask(c *gin.Context) {

}
