package main

import (
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/router"
	"github.com/labulaka521/crocodile/core/utils/log"
)

func main() {
	config.InitConf()
	log.InitLog()
	model.InitDb()
	model.InitRabc()
	router.InitRouter()
}
