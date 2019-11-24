package router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/middleware"
	"github.com/labulaka521/crocodile/core/router/api/v1/hostgroup"
	"github.com/labulaka521/crocodile/core/router/api/v1/task"
	"github.com/labulaka521/crocodile/core/router/api/v1/user"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func NewHttpRouter() *http.Server {
	router := gin.New()
	//gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Recovery(), middleware.ZapLogger(), middleware.PermissionControl())

	v1 := router.Group("/api/v1")
	ru := v1.Group("/user")
	{
		ru.POST("", user.RegistryUser)
		ru.GET("", user.GetUser)
		ru.PUT("", user.ChangeUser)
		ru.GET("/infos", user.GetUsers)
		ru.POST("/login", user.LoginUser)
	}
	rhg := v1.Group("/hostgroup")
	{
		rhg.GET("", hostgroup.GetHostGroups)
		rhg.POST("", hostgroup.CreateHostGroup)
		rhg.PUT("", hostgroup.ChangeHostGroup)
		rhg.DELETE("", hostgroup.DeleteHostGroup)
	}
	rt := v1.Group("/task")
	{
		rt.GET("", task.GetTasks)
		rt.POST("", task.CreateTask)
		rt.PUT("", task.ChangeTask)
		rt.DELETE("", task.DeleteTask)
	}
	httpSrv := &http.Server{
		Handler:      router,
		ReadTimeout:  config.CoreConf.Server.MaxHttpTime.Duration,
		WriteTimeout: config.CoreConf.Server.MaxHttpTime.Duration,
	}
	return httpSrv

}

func GetListen(mode define.RunMode) (net.Listener, error) {
	var (
		addr string
	)
	switch mode {
	case define.Server:
		addr = fmt.Sprintf(":%d", config.CoreConf.Server.Port)
	case define.Client:
		addr = fmt.Sprintf(":%d", config.CoreConf.Client.Port)
	default:
		return nil, errors.New("Unsupport mode")
	}
	lis, err := net.Listen("tcp", addr)

	return lis, err
}

func Run(mode define.RunMode, lis net.Listener) error {
	var (
		gRPCServer *grpc.Server
		err        error
		m          cmux.CMux
	)

	gRPCServer, err = schedule.NewgRPCServer(mode)
	if err != nil {
		return err
	}

	m = cmux.New(lis)
	if mode == define.Server {
		httpServer := NewHttpRouter()
		httpL := m.Match(cmux.HTTP1Fast())
		go httpServer.Serve(httpL)
		log.Info("Start Run Http Server", zap.String("addr", lis.Addr().String()))
	}
	grpcL := m.Match(cmux.Any())
	go gRPCServer.Serve(grpcL)
	log.Info("Start Run gRPC Server", zap.String("addr", lis.Addr().String()))

	return m.Serve()
}
