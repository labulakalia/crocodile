package main

import (
	"crocodile/common/cfg"
	"crocodile/common/db/mysql"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/common/wrapper"
	"crocodile/service/actuator/model/actuator"
	"github.com/labulaka521/logging"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"

	"crocodile/service/actuator/handler"
	"database/sql"
	"github.com/micro/go-micro"
	"time"

	pbactuator "crocodile/service/actuator/proto/actuator"
	goopentracing "github.com/opentracing/opentracing-go"
)

func main() {
	var (
		err error
		db  *sql.DB
		h   *handler.Actua
	)

	// New Service
	cfg.Init()
	log.Init()
	t, io, err := wrapper.NewTracer("crocodile.srv.actuator", "")
	if err != nil {
		logging.Fatal(err)
	}
	defer io.Close()
	goopentracing.SetGlobalTracer(t)
	// New Service
	service := micro.NewService(
		micro.Name("crocodile.srv.actuator"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
		micro.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		micro.WrapHandler(opentracing.NewHandlerWrapper()),
	)

	// Initialise service
	service.Init()
	db = mysql.New(cfg.MysqlConfig.DSN, cfg.MysqlConfig.MaxIdleConnection, cfg.MysqlConfig.MaxIdleConnection)

	h = &handler.Actua{
		&actuator.Service{
			DB: db,
		},
	}

	// Register Handler
	err = pbactuator.RegisterActuatorHandler(service.Server(), h)
	if err != nil {
		logging.Fatal(err)
	}

	// Run service
	if err = service.Run(); err != nil {
		logging.Fatal(err)
	}
}
