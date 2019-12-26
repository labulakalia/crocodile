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

// get hosts
func GetHost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	hosts, err := model.GetHost(ctx)

	if err != nil {
		log.Error("GetHost failed", zap.String("error", err.Error()))
		resp.Json(c, resp.ErrInternalServer, nil)
		return
	}
	resp.Json(c, resp.Success, hosts)
}

// 暂停主机分配任务
// PUT /api/v1/host
func StopHost(c *gin.Context) {
	panic("implentment me")
}

// 删除主机
// Delete /api/v1/host
// 需要从所有的主机组中找出主机id并删除
func DeleteHost(c *gin.Context) {
	panic("implentment me")
}
