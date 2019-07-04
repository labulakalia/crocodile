package event

const (
	Run_Task  int32 = 1
	Kill_Task int32 = 2
)

var MsgData = map[int32]string{
	Run_Task:  "Run Task",
	Kill_Task: "Kill Task ",
}

// 获取请求的消息
func GetEvent(code int32) string {
	var (
		msg    string
		exists bool
	)
	if msg, exists = MsgData[code]; exists {
		return msg
	}
	return MsgData[code]
}
