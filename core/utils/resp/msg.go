package resp

import "errors"

var msgcode = map[int]string{
	Success: "ok",

	ErrUnauthorized: "非法请求",
	ErrBadRequest:   "请求参数错误",

	ErrUserPassword:  "用户名或者密码错误",
	ErrUserForbid:    "禁止登陆",
	ErrEmailExist:    "邮箱已经存在",
	ErrUserNameExist: "用户名已存在",
	ErrUserNotExist:  "用户不存在",

	ErrTaskExist:    "任务名已存在",
	ErrTaskNotExist: "任务不存在",

	ErrHostgroupExist:    "主机组已存在",
	ErrHostgroupNotExist: "主机组不存在",

	ErrExecPlanExist:    "执行计划已存在",
	ErrExecPlanNotExist: "执行计划不存在",

	ErrInternalServer: "服务端错误",

	ErrRpcDeadlineExceeded: "调用超时",
	ErrRpcCanceled:         "取消调用",
	ErrRpcUnauthenticated:  "密钥认证失败",
	ErrRpcUnavailable:      "调用对端不可用",
	ErrRpcUnknow:           "调用未知错误",
	ErrRpcNotValidHost:     "未发现worker",
	ErrRpcNotConn:          "连接目标主机失败",
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

func GetMsgErr(code int) error {
	msg := GetMsg(code)
	return errors.New(msg)
}
