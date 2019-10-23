package config

type Dbcfg struct {
	Drivename    string
	Dsn          string
	MaxIdle      int
	MaxConn      int
	MaxQueryTime int
}
