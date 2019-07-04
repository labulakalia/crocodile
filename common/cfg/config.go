package cfg

import (
	"github.com/micro/go-micro/config"
)

var (
	MysqlConfig   defaultMysqlConfig
	JwtConfig     defaultJwtConfig
	EtcdConfig    defaultEtcdConfig
	ExecuteConfig defaultExecuteConfig
	LogConfig     defaultLog
)

func Init() {
	var (
		err      error
		filepath string
	)
	filepath = "../../conf/config.yaml"
	if err = config.LoadFile(filepath); err != nil {
		panic(err)
	}
	if err = config.Get("app", "mysql").Scan(&MysqlConfig); err != nil {
		panic(err)
	}
	if err = config.Get("app", "jwt").Scan(&JwtConfig); err != nil {
		panic(err)
	}
	if err = config.Get("app", "etcd").Scan(&EtcdConfig); err != nil {
		panic(err)
	}
	if err = config.Get("app", "executor").Scan(&ExecuteConfig); err != nil {
		panic(err)
	}
	if err = config.Get("app", "log").Scan(&LogConfig); err != nil {
		panic(err)
	}
}
