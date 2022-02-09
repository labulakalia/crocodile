package notify

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"go.uber.org/zap"
)

// GetNotify get self notify
func GetNotify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	uid := c.GetString("uid")
	notifys, err := model.GetNotify(ctx, uid)
	resp.JSONv2(c, err, notifys)
	return
}

// ReadNotify make notify status is read
func ReadNotify(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	type notifyid struct {
		ID int `json:"id"`
	}
	nuid := notifyid{}
	err := c.ShouldBindJSON(&nuid)
	if err != nil {
		log.Error("c.ShouldBindJSON failed", zap.Error(err))
		resp.JSONv2(c, err)
		return
	}
	err = model.NotifyRead(ctx, nuid.ID, c.GetString("uid"))
	if err != nil {
		log.Error("model.NotifyRead failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}
