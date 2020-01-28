package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	// "github.com/labulaka521/crocodile/common/errgroup"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/middleware"
	"github.com/labulaka521/crocodile/core/router/api/v1/host"
	"github.com/labulaka521/crocodile/core/router/api/v1/hostgroup"
	"github.com/labulaka521/crocodile/core/router/api/v1/task"
	"github.com/labulaka521/crocodile/core/router/api/v1/user"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// NewHTTPRouter create http.Server
func NewHTTPRouter() *http.Server {
	//gin.SetMode("release")
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
		rt.GET("/info", task.GetTask) // 获取指定task信息
		rt.POST("", task.CreateTask)
		rt.PUT("", task.ChangeTask)
		rt.DELETE("", task.DeleteTask)
		rt.PUT("/run", task.RunTask)
		rt.PUT("/kill", task.KillTask)
		rt.GET("/running", task.RunningTask)
		rt.GET("/log", task.LogTask)
		rt.GET("/log/ws", task.RealTimeLogTask)
	}
	rh := v1.Group("/host")
	{
		rh.GET("", host.GetHost)
		rh.PUT("", host.ChangeHostState)
		rh.DELETE("", host.DeleteHost)
	}

	httpSrv := &http.Server{
		Handler:      router,
		ReadTimeout:  config.CoreConf.Server.MaxHTTPTime.Duration,
		WriteTimeout: config.CoreConf.Server.MaxHTTPTime.Duration,
	}
	return httpSrv

}

// GetListen get listen addr by server or client
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

// Run start run http or grpc Server
func Run(mode define.RunMode, lis net.Listener) error {
	var (
		gRPCServer *grpc.Server
		httpServer *http.Server
		err        error
		m          cmux.CMux
	)

	gRPCServer, err = schedule.NewgRPCServer(mode)
	if err != nil {
		return err
	}

	m = cmux.New(lis)
	if mode == define.Server {
		httpServer = NewHTTPRouter()
		httpL := m.Match(cmux.HTTP1Fast())
		go httpServer.Serve(httpL)
		log.Info("start run http server", zap.String("addr", lis.Addr().String()))
	}

	grpcL := m.Match(cmux.Any())
	go gRPCServer.Serve(grpcL)
	log.Info("start run grpc server", zap.String("addr", lis.Addr().String()))

	go tryDisConn(gRPCServer, httpServer, mode)
	return m.Serve()
}

// tryDisConn will close grpc and http conn
// if time rather than 10s, will immediately close
func tryDisConn(gRPCServer *grpc.Server, httpServer *http.Server, mode define.RunMode) {

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGKILL, syscall.SIGSTOP,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT, syscall.SIGSYS,
	)

	select {
	case sig := <-signals:
		go func() {
			select {
			case <-time.After(time.Second * 10):
				log.Warn("Shutdown gracefully timeout, application will shutdown immediately.")
				os.Exit(0)
			}
		}()
		log.Info(fmt.Sprintf("get signal %s, application will shutdown.", sig))
		schedule.DoStopConn(mode)
		time.Sleep(time.Second)
		// g := errgroup.Group{}
		gRPCServer.Stop()
		if mode == define.Server {
			httpServer.Shutdown(context.Background())
		}
		// g.Wait()
		time.Sleep(time.Second * 11)
		fmt.Println("sleep")
		os.Exit(0)
	}

}
