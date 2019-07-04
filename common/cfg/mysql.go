package cfg

type defaultMysqlConfig struct {
	DSN               string `json:"dsn"`
	MaxIdleConnection int    `json:"maxIdleConnection"`
	MaxOpenConnection int    `json:"maxOpenConnection"`
	MaxQueryTime      int    `json:"maxQueryTime"`
}
