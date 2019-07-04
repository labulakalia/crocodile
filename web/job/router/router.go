package router

import (
	"crocodile/common/middle"
	"crocodile/web/job/router/job"
	"github.com/gin-gonic/gin"
)

func NewRouter() (r *gin.Engine) {
	var (
		apiv1job *gin.RouterGroup
	)
	gin.DisableConsoleColor()
	r = gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middle.MiddleJwt())

	apiv1job = r.Group("/job")
	apiv1job.Use(JobControl())
	{
		apiv1job.POST("/create", job.CreateJob)
		apiv1job.DELETE("/delete", job.DeleteJob)
		apiv1job.PUT("/change", job.ChangeJob)
		apiv1job.GET("/list", job.GetJob)
		apiv1job.POST("/kill", job.KillJob)
		apiv1job.POST("/run", job.RunJob)
		apiv1job.GET("/log", job.GetJobLog)
	}
	return r
}
