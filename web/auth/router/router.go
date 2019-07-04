package router

import (
	"crocodile/common/middle"
	"crocodile/web/auth/router/user"
	"github.com/gin-gonic/gin"
)

func NewRouter() (r *gin.Engine) {
	var (
		apiv1user *gin.RouterGroup
	)
	r = gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middle.MiddleJwt())

	apiv1user = r.Group("/auth")
	apiv1user.Use(UserControl())
	{
		apiv1user.GET("/info", user.GetUser)
		apiv1user.GET("/infos", user.GetUsers)
		apiv1user.PUT("/info", user.ChangeUser)
		apiv1user.POST("/create", user.UserCreate)
		apiv1user.POST("/login", user.UserLogin)
		apiv1user.POST("/logout", user.Logout)
		apiv1user.DELETE("/delete", user.DeleteUser)
	}
	return r
}
