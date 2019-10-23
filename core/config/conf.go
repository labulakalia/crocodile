package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

var (
	CoreConf *coreConf
)

type coreConf struct {
	Log LogConfig `toml:"log"`
	Db  Dbcfg     `toml:"db"`
}

func InitConf() {
	var (
		cfgpath string
		err     error
	)

	flag.StringVar(&cfgpath, "c", "core.toml", "core config")
	flag.Parse()
	_, err = toml.DecodeFile(cfgpath, &CoreConf)
	if err != nil {
		panic(fmt.Sprintf("Can Read Toml File %s", CoreConf))
	}
}
