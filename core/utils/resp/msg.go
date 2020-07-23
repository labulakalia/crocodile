package resp

import "errors"

var msgcode = map[int]string{
	Success: "ok",

	ErrUnauthorized: "非法请求",
	ErrBadRequest:   "请求参数错误",

	ErrUserPassword:  "用户或密码错误",
	ErrUserForbid:    "禁止登陆",
	ErrEmailExist:    "邮箱已经存在",
	ErrUserNameExist: "用户名已存在",
	ErrUserNotExist:  "用户不存在",

	ErrTaskExist:    "任务名已存在",
	ErrTaskNotExist: "任务不存在",

	ErrHostgroupExist:      "主机组已存在",
	ErrHostgroupNotExist:   "主机组不存在",
	ErrDelHostUseByOtherHG: "正在被其他的主机组使用，不能删除",
	ErrHostNotExist:        "主机不存在",

	ErrCronExpr: "CronExpr表达式不规范",

	ErrTaskUseByOtherTask:    "存在任务依赖此任务，请先在其他的任务的父子任务中移除此任务",
	ErrDelHostGroupUseByTask: "正在被其他的任务使用，不能删除",
	ErrDelUserUseByOther:     "请先删除此用户创建的主机组或者任务后再删除",

	ErrInternalServer: "服务端错误",

	ErrCtxDeadlineExceeded: "调用超时",
	ErrCtxCanceled:         "取消调用",
	ErrRPCUnauthenticated:  "密钥认证失败",
	ErrRPCUnavailable:      "调用对端不可用",
	ErrRPCUnknow:           "调用未知错误",
	ErrRPCNotValidHost:     "未发现worker",
	ErrRPCNotConnHost:      "未找到存活的worker",

	NeedInstall:   "系统还未安装，请等待安装后再进行操作",
	IsInstall:     "系统已经安装完成，请勿再次执行安装操作",
	ErrInstall:    "安装失败",
	ErrDBConnFail: "数据库连接失败",
}

// GetMsg get msg by code
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

// GetMsgErr get error msg by code
func GetMsgErr(code int) error {
	msg := GetMsg(code)
	return errors.New(msg)
}
