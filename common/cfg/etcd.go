package cfg

type defaultEtcdConfig struct {
	Endpoints   []string `json:"endpoints"`
	DialTimeout int      `json:"dialTimeout"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
}
