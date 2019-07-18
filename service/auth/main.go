package main

import (
	"crocodile/common/cfg"
	"crocodile/common/db/mysql"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/common/wrapper"
	"crocodile/service/auth/handler"
	"crocodile/service/auth/model/user"
	pbauth "crocodile/service/auth/proto/auth"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/trace/opentracing"
	goopentracing "github.com/opentracing/opentracing-go"
	"time"
)

func main() {
	var (
		h   *handler.Auth
		err error
	)

	// New Service
	cfg.Init()
	log.Init()
	t, io, err := wrapper.NewTracer("crocodile.srv.auth", "")
	if err != nil {
		logging.Fatal(err)
	}
	defer io.Close()
	goopentracing.SetGlobalTracer(t)

	service := micro.NewService(
		micro.Name("crocodile.srv.auth"),
		micro.Version("latest"),
		// 注册的有效时长
		micro.RegisterTTL(time.Second*30),
		// 每隔15秒注册一次
		micro.RegisterInterval(time.Second*15),
		micro.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		micro.WrapHandler(opentracing.NewHandlerWrapper()),
	)
	// Initialise service
	service.Init()

	s := &user.Service{
		DB: mysql.New(cfg.MysqlConfig.DSN, cfg.MysqlConfig.MaxIdleConnection, cfg.MysqlConfig.MaxIdleConnection),
	}

	h = &handler.Auth{s}

	// Register Handler
	if err = pbauth.RegisterAuthHandler(service.Server(), h); err != nil {
		logging.Fatalf("RegisterAuthHandler Err: %v", err)
	}

	// Run service
	if err := service.Run(); err != nil {
		logging.Fatal(err)
	}
}
