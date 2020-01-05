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

// GetHost get all hosts, online and offline host
// GET /api/v1/host
func GetHost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hosts, err := model.GetHost(ctx)

	if err != nil {
		log.Error("GetHost failed", zap.String("error", err.Error()))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, hosts)
}

// ChangeHostState stop run worker
// PUT /api/v1/host/stop
func ChangeHostState(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hosttask := define.GetTaskid{}
	err := c.ShouldBindJSON(&hosttask)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(hosttask.ID) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	host, err := model.GetHostByID(ctx,hosttask.ID)
	if err != nil {
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}

	err = model.StopHost(ctx , hosttask.ID, host.Stop ^ 1)
	if err != nil {
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// DeleteHost delete host from
// DELETE /api/v1/host
func DeleteHost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
	config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hosttask := define.GetTaskid{}
	err := c.ShouldBindJSON(&hosttask)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	if utils.CheckID(hosttask.ID) != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}
	err = model.DeleteHost(ctx, hosttask.ID)
	if err != nil {
		resp.JSON(c,resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}
