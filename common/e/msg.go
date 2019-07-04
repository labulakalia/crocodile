package e

var MsgData = map[int32]string{
	SUCCESS: "success",

	FORBIDDEN: "forbid",

	ERR_CREATE_USER_FAIL: "创建用户失败",
	ERR_CHANGE_USER_FAIL: "修改用户失败",
	ERR_GET_USER_FAIL:    "获取用户失败",
	ERR_USER_PASS_FAIL:   "用户密码错误",
	ERR_USER_NOT_EXIST:   "用户不存在",
	ERR_USER_EXIST:       "用户已经存在",
	ERR_DELETE_USER_FAIL: "删除用户失败",
	ERR_LOGIN_USER_FAIL:  "用户登录失败",
	ERR_LOGOUT_USER_FAIL: "用户注销失败",
	ERR_NOT_ALLOW_LOGIN:  "用户禁止登录",

	ERR_NOT_PERMISSION: "没有权限",
	ERR_TOKEN_INVALID:  "Token无效",

	ERR_BAD_REQUEST:            "请求参数错误",
	ERR_GENERATE_HASHPASS_FAIL: "生成哈希密码失败",
	ERR_GENERATE_TOKEN_FAIL:    "生成TOKEN失败",

	ERR_SQL_FAIl: "SQL 错误",

	ERR_CREATE_JOB_FAIL:      "创建任务失败",
	ERR_DELETE_JOB_FAIL:      "删除任务失败",
	ERR_CHANGE_JOB_FAIL:      "修改任务失败",
	ERR_GET_JOB_FAIL:         "获取任务失败",
	ERR_KILL_JOB_FAIL:        "强杀任务失败",
	ERR_RUN_JOB_FAIL:         "运行任务失败",
	ERR_JOB_NOT_EXITS:        "任务不存在",
	ERR_JOB_EXITS:            "任务存在",
	ERR_CREATE_ACTUAT_FAIL:   "创建执行器失败",
	ERR_DELETE_ACTUAT_FAIL:   "删除执行器失败",
	ERR_CHANGE_ACTUAT_FAIL:   "修改执行器失败",
	ERR_GET_ACTUAT_FAIL:      "获取执行器失败",
	ERR_GET_EXECUTOR_IP_FAIL: "获取执行主机IP失败",
	ERR_GET_JOB_LOG_FAIL:     "获取任务日志失败",
	ERR_ACTUATOR_NOT_EXITS:   "执行器不存在",
}

// 获取请求的消息
func GetMsg(code int32) string {
	var (
		msg    string
		exists bool
	)
	if msg, exists = MsgData[code]; exists {
		return msg
	}
	return MsgData[code]
}
