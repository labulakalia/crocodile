package model

import (
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"go.uber.org/zap"
)

func InitDb() {
	var (
		err error
	)
	dbcfg := config.CoreConf.Db
	err = db.NewDb(db.Drivename(dbcfg.Drivename),
		db.Dsn(dbcfg.Dsn),
		db.MaxIdleConnection(dbcfg.MaxIdle),
		db.MaxOpenConnection(dbcfg.MaxConn),
		db.MaxQueryTime(dbcfg.MaxQueryTime),
	)
	if err != nil {
		log.Fatal("InitDb failed", zap.Error(err))
	}
}
