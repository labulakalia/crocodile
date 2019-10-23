package router

import (
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/middleware"
	"github.com/labulaka521/crocodile/core/router/api/v1/user"
	"go.uber.org/zap"
)

func InitRouter() {
	router := gin.New()
	//gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Recovery(), middleware.ZapLogger(), middleware.PermissionControl())

	ru := router.Group("/api/v1/user")
	{
		ru.POST("/login", user.LoginUser)
		ru.POST("/registry", user.RegistryUser)
		ru.GET("/info", user.GetUser)
		ru.GET("/infos", user.GetUsers)
		ru.PUT("/change", user.ChangeUser)
	}
	err := router.Run("127.0.0.1:8080")
	log.Info("Run Server failed", zap.Error(err))
}
