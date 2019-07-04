package cfg

type defaultLog struct {
	Level string `json:"level"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
}
