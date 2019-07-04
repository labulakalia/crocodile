package main

import (
	"crocodile/common/cfg"
	"crocodile/common/log"
	"crocodile/common/registry"
	"crocodile/web/job/router"
	"crocodile/web/job/router/job"
	"github.com/labulaka521/logging"
	"github.com/micro/cli"
	"time"

	"github.com/micro/go-micro/web"
)

func main() {
	cfg.Init()
	log.Init()
	// create new web service
	service := web.NewService(
		web.Name("crocodile.web.job"),
		web.Version("latest"),
		web.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
		web.RegisterInterval(15*time.Second),
		web.RegisterTTL(30*time.Second),
	)

	// initialise service
	err := service.Init(
		web.Action(func(c *cli.Context) {
			job.Init()
		}),
	)
	if err != nil {
		logging.Error(err)
	}
	// register html handler
	// 路由的开头也必须是auth开头才可以
	service.Handle("/", router.NewRouter())

	// run service
	if err := service.Run(); err != nil {
		logging.Fatal(err)
	}
}
