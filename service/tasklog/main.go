package main

import (
	"crocodile/common/cfg"
	"crocodile/common/db/mysql"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/common/wrapper"
	"crocodile/service/tasklog/handler"
	"crocodile/service/tasklog/model/tasklog"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	"time"

	pbtasklog "crocodile/service/tasklog/proto/tasklog"
	goopentracing "github.com/opentracing/opentracing-go"
)

func main() {
	var (
		err error
	)
	cfg.Init()
	log.Init()
	t, io, err := wrapper.NewTracer("crocodile.srv.tasklog", "")
	if err != nil {
		logging.Fatal(err)
	}
	defer io.Close()
	goopentracing.SetGlobalTracer(t)

	// New Service
	service := micro.NewService(
		micro.Name("crocodile.srv.tasklog"),
		micro.Version("latest"),
		micro.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		micro.RegisterInterval(15*time.Second),
		micro.RegisterTTL(30*time.Second),
		micro.WrapHandler(opentracing.NewHandlerWrapper()),
	)

	// Initialise service
	service.Init()

	db := mysql.New(cfg.MysqlConfig.DSN, cfg.MysqlConfig.MaxIdleConnection, cfg.MysqlConfig.MaxIdleConnection)

	h := &handler.TaskLog{
		Service: &tasklog.Service{
			DB: db,
		},
	}

	// Register Handler
	if err = pbtasklog.RegisterTaskLogHandler(service.Server(), h); err != nil {
		logging.Fatal("RegistryTaskLogHandler Err:%V", err)
	}

	// Run service
	if err := service.Run(); err != nil {
		logging.Fatal(err)
	}
}
