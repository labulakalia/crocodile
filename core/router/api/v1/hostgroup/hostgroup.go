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
// @Summary create hostgroup
// @Tags HostGroup
// @Description create new hostgroup
// @Produce json
// @Param User body define.CreateHostGroup true "HostGroup"
// @Success 200 {object} resp.Response
// @Router /api/v1/hostgroup [post]
// @Security ApiKeyAuth
func CreateHostGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	hg := define.CreateHostGroup{}

	err := c.ShouldBindJSON(&hg)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	exist, err := model.Check(ctx, model.TBHostgroup, model.Name, hg.Name)
	if err != nil {
		log.Error("IsExist failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if exist {
		resp.JSON(c, resp.ErrHostgroupExist, nil)
		return
	}

	err = model.CreateHostgroup(ctx, hg.Name, hg.Remark, c.GetString("uid"), hg.HostsID)
	if err != nil {
		log.Error("CreateHostgroup failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// ChangeHostGroup change hostgroup
// @Summary change hostgroup
// @Tags HostGroup
// @Description change hostgroup
// @Produce json
// @Param User body define.ChangeHostGroup true "HostGroup"
// @Success 200 {object} resp.Response
// @Router /api/v1/hostgroup [put]
// @Security ApiKeyAuth
func ChangeHostGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	hg := define.ChangeHostGroup{}

	err := c.ShouldBindJSON(&hg)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	// 判断ID是否存在
	exist, err := model.Check(ctx, model.TBHostgroup, model.ID, hg.ID)
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
		exist, err = model.Check(ctx, model.TBHostgroup, model.IDCreateByUID, hg.ID, uid)
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

	err = model.ChangeHostGroup(ctx, hg.HostsID, hg.ID, hg.Remark)
	if err != nil {
		log.Error("ChangeHostGroup failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)

}

// DeleteHostGroup delete hostgroup
// @Summary delete hostgroup
// @Tags HostGroup
// @Description delete hostgroup
// @Produce json
// @Param User body define.GetID true "HostGroup"
// @Success 200 {object} resp.Response
// @Router /api/v1/hostgroup [delete]
// @Security ApiKeyAuth
func DeleteHostGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	hostgroup := define.GetID{}

	err := c.ShouldBindJSON(&hostgroup)
	if err != nil {
		log.Error("ShouldBindJSON failed", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(hostgroup.ID) != nil {
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

// GetHostGroups get all hostgroup
// @Summary get all hostgroup
// @Tags HostGroup
// @Description get all hostgroup
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/hostgroup [get]
// @Security ApiKeyAuth
func GetHostGroups(c *gin.Context) {
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
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}

	hgs, count, err := model.GetHostGroups(ctx, q.Limit, q.Offset)
	if err != nil {
		log.Error("GetHostGroup failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, hgs, count)
}

// GetSelect return name,id
// @Summary get name,id
// @Tags HostGroup
// @Description get select option
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/hostgroup [get]
// @Security ApiKeyAuth
func GetSelect(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	data, err := model.GetNameID(ctx, model.TBHostgroup)
	if err != nil {
		log.Error("model.GetNameID", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, data)
}

// GetHostsByIHGID get host detail by hostgroup id
// @Summary get host detail by hostgroup id
// @Tags HostGroup
// @Description get all hostgroup
// @Param id query string false "ID"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/hostgroup/hosts [get]
// @Security ApiKeyAuth
func GetHostsByIHGID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	getid := define.GetID{}
	err := c.BindQuery(&getid)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	hosts, err := model.GetHostsByHGID(ctx, getid.ID)
	switch err.(type) {
	case nil:
		resp.JSON(c, resp.Success, hosts)
	case define.ErrNotExist:
		resp.JSON(c, resp.ErrHostgroupNotExist, nil)
	default:
		log.Error("model.GetHostsByHGID", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
}
