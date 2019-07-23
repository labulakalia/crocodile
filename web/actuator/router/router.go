package router

import (
	"crocodile/common/middle"
	"crocodile/common/wrapper"
	"crocodile/web/actuator/router/actuator"
	"github.com/gin-gonic/gin"
)

func NewRouter() (r *gin.Engine) {
	var (
		apiv1actuator *gin.RouterGroup
	)
	gin.DisableConsoleColor()
	r = gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middle.MiddleJwt())
	r.Use(wrapper.TracerWrapper)
	apiv1actuator = r.Group("/actuator")
	apiv1actuator.Use(ActuatorControl())
	//apiv1actuator.Group()
	{
		apiv1actuator.POST("/create", actuator.CreateActuator)
		apiv1actuator.DELETE("/delete", actuator.DeleteActuator)
		apiv1actuator.PUT("/change", actuator.ChangeActuator)
		apiv1actuator.GET("/list", actuator.GetActuator)
		apiv1actuator.GET("/executeip", actuator.GetALLExecuteIP)
	}
	apiv1actuator1 := r.Group("/debug")
	{
		apiv1actuator1.POST("/create", actuator.CreateActuator)
	}
	return r
}
