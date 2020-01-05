package host

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
)

// GetHost get all hosts, online and offline host
// GET 
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

// StopHost stop run worker
// PUT /api/v1/host/stop
func StopHost(c *gin.Context) {
	panic("implentment me")

}

// DeleteHost delete host from 
// DELETE /api/v1/host
func DeleteHost(c *gin.Context) {
	panic("implentment me")
}
