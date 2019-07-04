package main

import (
	"crocodile/common/cfg"
	"crocodile/common/db/mysql"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/service/tasklog/handler"
	"crocodile/service/tasklog/model/tasklog"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro"
	"time"

	pbtasklog "crocodile/service/tasklog/proto/tasklog"
)

func main() {
	var (
		err error
	)
	cfg.Init()
	log.Init()

	// New Service
	service := micro.NewService(
		micro.Name("crocodile.srv.tasklog"),
		micro.Version("latest"),
		micro.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		micro.RegisterInterval(15*time.Second),
		micro.RegisterTTL(30*time.Second),
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
