package hostgroup

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
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
		resp.JSONv2(c, err)
		return
	}
	err = model.CreateHostgroupv2(ctx, hg.Name, hg.Remark, c.GetString("uid"), hg.HostsID)
	resp.JSONv2(c, err)
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
		resp.JSONv2(c, err)
		return
	}
	// 获取用户的类型
	var role define.Role
	if v, ok := c.Get("role"); ok {
		role = v.(define.Role)
	}
	// var currentUID string
	// 这里只需要确定如果rule的用户类型是否为Admin
	if role != define.AdminUser {
		ctx = context.WithValue(ctx, "uid", c.GetString("uid"))
	}

	err = model.ChangeHostGroupv2(ctx, hg.HostsID, hg.ID, hg.Remark)
	resp.JSONv2(c, err, nil)

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
	}

	err = model.DeleteHostGroupv2(ctx, hostgroup.ID)
	resp.JSONv2(c, err)
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
		resp.JSONv2(c, err, nil)
		return
	}

	if q.Limit == 0 {
		q.Limit = define.DefaultLimit
	}

	hgs, count, err := model.GetHostGroups(ctx, q.Limit, q.Offset)
	resp.JSONv2(c, err, hgs, count)
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
	data, err := model.GetIDNameOption(ctx, nil, &model.HostGroup{})
	resp.JSONv2(c, err, data)
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
		resp.JSONv2(c, err, nil)
		return
	}
	hosts, err := model.GetHostsByHGIDv2(ctx, getid.ID)
	resp.JSONv2(c, err, hosts)
}
