package main

import (
	"crocodile/common/cfg"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/common/wrapper"
	"crocodile/web/auth/router/user"
	"time"

	"crocodile/web/auth/router"

	"github.com/labulaka521/logging"
	"github.com/micro/cli"

	"github.com/micro/go-micro/web"
	goopentracing "github.com/opentracing/opentracing-go"
)

func main() {
	cfg.Init()
	log.Init()

	t, io, err := wrapper.NewTracer("crocodile.web.auth", "")
	if err != nil {
		logging.Fatal(err)
	}
	defer io.Close()
	goopentracing.SetGlobalTracer(t)
	// create new web service

	service := web.NewService(
		web.Name("crocodile.web.auth"),
		web.Version("latest"),
		web.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		web.RegisterInterval(15*time.Second),
		web.RegisterTTL(30*time.Second),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(c *cli.Context) {
			user.Init()
		}),
	); err != nil {
		logging.Fatal(err)
	}

	// register html handler
	// 路由的开头也必须是auth开头才可以
	service.Handle("/", router.NewRouter())

	// run service
	if err := service.Run(); err != nil {
		logging.Fatal(err)
	}
}
