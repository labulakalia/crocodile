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

// CreateHostGroup create hostgroup
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	if hostgroup.Name == "" {
		log.Error("Hostgroup.Name is empty")
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	exist, err := model.Check(ctx, model.TBHostgroup, model.Name, hostgroup.Name)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrHostgroupExist, nil)
		return
	}

	hostgroup.ID = utils.GetID()
	hostgroup.CreateByUID = c.GetString("uid")

	err = model.CreateHostgroup(ctx, &hostgroup)
	if err != nil {
		log.Error("CreateHostgroup failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// ChangeHostGroup change hostgroup
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	// 判断ID是否存在
	exist, err := model.Check(ctx, model.TBHostgroup, model.ID, hostgroup.ID)
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
		exist, err = model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, hostgroup.ID, uid)
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

	err = model.ChangeHostGroup(ctx, &hostgroup)
	if err != nil {
		log.Error("ChangeHostGroup failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)

}

// DeleteHostGroup deletehostgroup
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	// 判断ID是否存在
	exist, err := model.Check(ctx, model.TBHostgroup, model.ID, hostgroup.ID)
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
		exist, err = model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, hostgroup.ID, uid)
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

	err = model.DeleteHostGroup(ctx, hostgroup.ID)
	if err != nil {
		log.Error("DeleteHostGroup failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// GetHostGroups get host groups
// GET /api/v1/hostgroup
//
func GetHostGroups(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hgs, err := model.GetHostGroups(ctx)
	if err != nil {
		log.Error("GetHostGroup failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, hgs)

}
