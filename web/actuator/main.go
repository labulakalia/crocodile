package main

import (
	"crocodile/common/cfg"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/web/actuator/router"
	"crocodile/web/actuator/router/actuator"
	"github.com/labulaka521/logging"
	"github.com/micro/cli"
	"github.com/micro/go-micro/web"
	"time"
)

func main() {
	cfg.Init()
	log.Init()
	// create new web service
	service := web.NewService(
		web.Name("crocodile.web.actuator"),
		web.Version("latest"),
		web.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		web.RegisterInterval(15*time.Second),
		web.RegisterTTL(30*time.Second),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(c *cli.Context) {
			actuator.Init()
		}),
	); err != nil {
		logging.Fatal(err)
	}

	// register html handler
	service.Handle("/", router.NewRouter())

	// run service
	if err := service.Run(); err != nil {
		logging.Fatal(err)
	}
}
