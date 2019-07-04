package main

import (
	"crocodile/common/cfg"
	"crocodile/common/db/mysql"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/service/job/handler"
	"crocodile/service/job/model/task"
	"crocodile/service/job/scheduler"
	"database/sql"
	"github.com/labulaka521/logging"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"time"

	pbjob "crocodile/service/job/proto/job"
)

func main() {
	var (
		err  error
		db   *sql.DB
		h    *handler.Job
		exit chan int
	)

	// New Service
	cfg.Init()
	log.Init()

	service := micro.NewService(
		micro.Name("crocodile.srv.job"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		micro.AfterStart(func() error {
			exit = make(chan int, 1)
			// 启动调度中心
			go scheduler.Loop(exit, db)
			return nil
		}),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			task.Init(service.Client())
		}),
	)

	db = mysql.New(cfg.MysqlConfig.DSN, cfg.MysqlConfig.MaxIdleConnection, cfg.MysqlConfig.MaxIdleConnection)

	h = &handler.Job{
		&task.Service{
			DB: db,
		},
	}

	// Register Handler
	err = pbjob.RegisterJobHandler(service.Server(), h)
	if err != nil {
		logging.Fatalf("RegisterJobHandler Err: %v", err)
	}

	// Run service
	if err = service.Run(); err != nil {
		logging.Fatal(err)
	}
	exit <- 0

	logging.Info("Exiting Job Service...")
}
