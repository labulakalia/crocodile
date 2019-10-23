package resp

var msgcode = map[int]string{
	Success: "ok",

	ErrUnauthorized: "非法请求",
	ErrBadRequest:   "请求参数错误",

	ErrUserPassword:  "用户或者密码错误",
	ErrUserForbid:    "用户被禁止登陆",
	ErrUserNameExist: "用户名已经存在",
	ErrEmailExist:    "邮箱已经存在",
	ErrUserNotExist:  "用户不存在",

	ErrInternalServer: "服务端错误",
}

// 获取请求的消息
func GetMsg(code int) string {
	var (
		msg    string
		exists bool
	)

	if msg, exists = msgcode[code]; exists {
		return msg
	}
	return "unknown"
}
