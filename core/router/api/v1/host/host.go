package host

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

// GetHost return all registry gost
// @Summary get all hosts
// @Tags Host
// @Description get all registry host
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/host [get]
// @Security ApiKeyAuth
func GetHost(c *gin.Context) {
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

	hosts, count, err := model.GetHosts(ctx, q.Offset, q.Limit)

	if err != nil {
		log.Error("GetHost failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	resp.JSON(c, resp.Success, hosts, count)
}

// ChangeHostState stop host worker
// @Summary stop host worker
// @Tags Host
// @Description stop host worker
// @Param StopHost body define.GetID true "ID"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/host/stop [put]
// @Security ApiKeyAuth
func ChangeHostState(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	gethost := define.GetID{}
	err := c.ShouldBindJSON(&gethost)
	if err != nil {
		log.Error("c.ShouldBindJSON", zap.Error(err))
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(gethost.ID) != nil {
		log.Error("CheckID failed")
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	host, err := model.GetHostByID(ctx, gethost.ID)
	switch err.(type) {
	case nil:
		goto Next
	case define.ErrNotExist:
		resp.JSON(c, resp.ErrHostNotExist, nil)
		return
	default:

		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
Next:
	err = model.StopHost(ctx, gethost.ID, !host.Stop)
	if err != nil {
		log.Error("model.StopHost", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// DeleteHost delete host
// @Summary delete host
// @Tags Host
// @Description delete host
// @Param StopHost body define.GetID true "ID"
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/host [delete]
// @Security ApiKeyAuth
func DeleteHost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	gethost := define.GetID{}
	err := c.ShouldBindJSON(&gethost)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(gethost.ID) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	hostgroups, _, err := model.GetHostGroups(ctx, 0, 0)
	if err != nil {
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	for _, hostgroup := range hostgroups {
		for _, hid := range hostgroup.HostsID {
			if gethost.ID == hid {
				resp.JSON(c, resp.ErrDelHostUseByOtherHG, nil)
				return
			}
		}
	}

	err = model.DeleteHost(ctx, gethost.ID)
	if err != nil {
		log.Error("model.DeleteHost failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// GetSelect name,id
// @Summary Get Task Select
// @Tags Host
// @Produce json
// @Success 200 {object} resp.Response
// @Router /api/v1/host/select [get]
// @Security ApiKeyAuth
func GetSelect(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	data, err := model.GetNameID(ctx, model.TBHost)
	if err != nil {
		log.Error("model.GetNameID failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, data)
}
