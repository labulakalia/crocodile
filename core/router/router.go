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

	"github.com/labulaka521/crocodile/core/config"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-contrib/pprof"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/middleware"
	"github.com/labulaka521/crocodile/core/router/api/v1/host"
	"github.com/labulaka521/crocodile/core/router/api/v1/hostgroup"
	"github.com/labulaka521/crocodile/core/router/api/v1/install"
	"github.com/labulaka521/crocodile/core/router/api/v1/notify"
	"github.com/labulaka521/crocodile/core/router/api/v1/task"
	"github.com/labulaka521/crocodile/core/router/api/v1/user"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/utils/asset"
	"github.com/labulaka521/crocodile/core/utils/define"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	_ "github.com/labulaka521/crocodile/core/docs" // init swagger docs
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// NewHTTPRouter create http.Server
func NewHTTPRouter() *http.Server {
	//gin.SetMode("release")
	router := gin.New()

	fs := &assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
		Prefix:    "web/crocodile",
	}

	router.StaticFS("/crocodile", fs)
	router.GET("/static/*url", func(c *gin.Context) {
		pre, exist := c.Params.Get("url")
		if !exist {
			return
		}
		c.Redirect(http.StatusMovedPermanently, "/crocodile/static"+pre)
	})
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/crocodile/favicon.ico")
	})
	router.GET("/index.html", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/crocodile")
	})

	pprof.Register(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//gin.SetMode(gin.ReleaseMode)
	//,
	router.Use(gin.Recovery(), middleware.ZapLogger(), middleware.PermissionControl(), middleware.Oprtation())

	v1 := router.Group("/api/v1")
	ru := v1.Group("/user")
	{
		ru.POST("/registry", user.RegistryUser) // only admin // 管理员创建了新的用户。。。
		ru.GET("/info", user.GetUser)
		ru.GET("/all", user.GetUsers)          // only admin
		ru.PUT("/admin", user.AdminChangeUser) // only admin  // 管理员修改了某某用户
		ru.PUT("/info", user.ChangeUserInfo)   // 某某修改了个人信息
		ru.POST("/login", user.LoginUser)
		ru.POST("/logout", user.LogoutUser) // 某某注销登陆
		ru.GET("/select", user.GetSelect)
		ru.GET("/alarmstatus", user.GetAlarmStatus)
		ru.GET("/operate", user.GetOperateLog)
	}
	rhg := v1.Group("/hostgroup")
	{
		rhg.GET("", hostgroup.GetHostGroups)
		rhg.POST("", hostgroup.CreateHostGroup)
		rhg.PUT("", hostgroup.ChangeHostGroup)
		rhg.DELETE("", hostgroup.DeleteHostGroup)
		rhg.GET("/select", hostgroup.GetSelect)
		rhg.GET("/hosts", hostgroup.GetHostsByIHGID)
	}
	rt := v1.Group("/task")
	{
		rt.GET("", task.GetTasks)
		rt.GET("/info", task.GetTask)
		rt.POST("", task.CreateTask)
		rt.POST("/clone", task.CloneTask)
		rt.PUT("", task.ChangeTask)
		rt.DELETE("", task.DeleteTask)
		rt.PUT("/run", task.RunTask)
		rt.PUT("/kill", task.KillTask)
		rt.GET("/running", task.GetRunningTask)
		rt.DELETE("/log", task.CleanTaskLog)
		rt.GET("/log", task.LogTask)
		rt.GET("/log/tree", task.LogTreeData)
		rt.GET("/log/websocket", task.RealRunTaskLog)
		rt.GET("/status/websocket", task.RealRunTaskStatus)

		rt.GET("/cron", task.ParseCron)
		rt.GET("/select", task.GetSelect)
	}
	rh := v1.Group("/host")
	{
		rh.GET("", host.GetHost)
		rh.PUT("/stop", host.ChangeHostState)
		rh.DELETE("", host.DeleteHost)
		rh.GET("/select", host.GetSelect)
	}

	rn := v1.Group("/notify")
	{
		rn.GET("", notify.GetNotify)
		rn.PUT("", notify.ReadNotify)
	}
	ri := v1.Group("/install")
	{
		ri.GET("/status", install.QueryIsInstall)
		ri.POST("", install.StartInstall)
	}
	// if nor find router, will rediret to /crocodile/
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/crocodile/")
	})

	httpSrv := &http.Server{
		Handler: router,
		// ReadTimeout:  config.CoreConf.Server.MaxHTTPTime.Duration,
		// WriteTimeout: config.CoreConf.Server.MaxHTTPTime.Duration,
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
		if os.Getenv("PORT") != "" {
			addr = ":" + os.Getenv("PORT")
		} else {
			addr = fmt.Sprintf(":%d", config.CoreConf.Server.Port)
		}

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
	////
	grpcL := m.Match(cmux.Any())
	go gRPCServer.Serve(grpcL)
	log.Info("start run grpc server", zap.String("addr", lis.Addr().String()))

	//deploy heroku need custom port
	//if mode == define.Server {
	//	grpclis, err := net.Listen("tcp", ":8080")
	//	if err != nil {
	//		log.Error("net.Listen failed", zap.Error(err))
	//		return err
	//	}
	//	go gRPCServer.Serve(grpclis)
	//	log.Info("start run grpc server", zap.String("addr", grpclis.Addr().String()))
	//} else {
	//	grpcL := m.Match(cmux.Any())
	//	go gRPCServer.Serve(grpcL)
	//	log.Info("start run grpc server", zap.String("addr", lis.Addr().String()))
	//}

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

		// g := errgroup.Group{}
		log.Debug("Start Stop GrpcServer")
		gRPCServer.Stop()
		if mode == define.Server {
			log.Debug("Start Stop HttpServer")
			httpServer.Shutdown(context.Background())
		}
		// g.Wait()
		//time.Sleep(time.Second * 11)
		os.Exit(0)
	}

}
