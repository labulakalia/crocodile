package hostgroup

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
)

// HostGroup

// POST /api/v1/hostgroup
// @params
// name
// hosts [] option
// remark option
func CreateHostGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	hostgroup := define.HostGroup{}

	err := c.ShouldBindJSON(&hostgroup)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}

	if hostgroup.Name == "" {
		log.Error("Hostgroup.Name is empty")
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}

	exist, err := model.Check(ctx, model.TB_hostgroup, model.Name, hostgroup.Name)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.Json(c, resp.ErrHostgroupExist, nil)
		return
	}

	hostgroup.Id = utils.GetId()
	hostgroup.CreateByUId = c.GetString("uid")

	err = model.CreateHostgroup(ctx, &hostgroup)
	if err != nil {
		log.Error("CreateHostgroup failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, nil)
}

// PUT /api/v1/hostgroup
// @params
// id
// name
// hosts option
// remark option
func ChangeHostGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	hostgroup := define.HostGroup{}

	err := c.ShouldBindJSON(&hostgroup)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}
	// 判断ID是否存在
	exist, err := model.Check(ctx, model.TB_hostgroup, model.ID, hostgroup.Id)
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
		exist, err = model.Check(ctx, model.TB_hostgroup, model.IDCreateByUID, hostgroup.Id, uid)
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

	err = model.ChangeHostGroup(ctx, &hostgroup)
	if err != nil {
		log.Error("ChangeHostGroup failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, nil)

}

// DELETE /api/v1/hostgroup
// id
func DeleteHostGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	hostgroup := define.HostGroup{}

	err := c.ShouldBindJSON(&hostgroup)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.Json(c, resp.ErrBadRequest, nil)
		return
	}
	// 判断ID是否存在
	exist, err := model.Check(ctx, model.TB_hostgroup, model.ID, hostgroup.Id)
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
		exist, err = model.Check(ctx, model.TB_hostgroup, model.IDCreateByUID, hostgroup.Id, uid)
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

	err = model.DeleteHostGroup(ctx, hostgroup.Id)
	if err != nil {
		log.Error("DeleteHostGroup failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, nil)
}

// GET /api/v1/hostgroup
//
func GetHostGroups(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hgs, err := model.GetHostGroups(ctx)
	if err != nil {
		log.Error("GetHostGroup failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, hgs)

}
