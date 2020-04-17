package install

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/resp"
	"github.com/labulaka521/crocodile/core/version"
	"go.uber.org/zap"
)

// QueryIsInstall query system is installed
func QueryIsInstall(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()

	isinstall, err := model.QueryIsInstall(ctx)
	if err != nil {
		log.Error("model.QueryIsInstall failed", zap.Error(err))
		resp.JSON(c, resp.ErrInternalServer, nil)
	}
	if !isinstall {
		log.Debug("first running, need install...")
		resp.JSON(c, resp.NeedInstall, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// StartInstall install system
func StartInstall(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(),
		config.CoreConf.Server.DB.MaxQueryTime.Duration)
	defer cancel()
	isinstall, err := model.QueryIsInstall(ctx)
	if err != nil {
		resp.JSON(c, resp.ErrInternalServer, nil)
		return
	}
	if isinstall {
		resp.JSON(c, resp.IsInstall, nil)
		return
	}

	// get new create user
	adminuser := define.CreateAdminUser{}

	err = c.ShouldBindJSON(&adminuser)
	if err != nil {
		resp.JSON(c, resp.ErrBadRequest, nil)
		return
	}

	err = model.StartInstall(ctx, adminuser.Name, adminuser.Password)
	if err != nil {
		log.Error("model.StartInstall", zap.Error(err))
		resp.JSON(c, resp.ErrInstall, nil)
		return
	}
	resp.JSON(c, resp.Success, nil)
}

// QueryVersion query current version
func QueryVersion(c *gin.Context) {
	resp.JSON(c, resp.Success, version.Version)
}
