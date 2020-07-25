package define

const (
	// DefaultLimit set get total page
	DefaultLimit = 20
)

// Role Admin or Normal User
type Role uint8

const (
	// NormalUser define normal user
	NormalUser Role = iota + 1 // 普通用户 只对自已创建的主机或者主机组具有操作权限
	// AdminUser define admin user
	AdminUser // 管理员 具有所有操作
	// GuestUser only look
	GuestUser // 访客 只有查看的权限
)

func (r Role) String() string {
	switch r {
	case AdminUser:
		return "Admin"
	case NormalUser:
		return "Normal"
	case GuestUser:
		return "Guest"
	default:
		return "Unknown"
	}
}

// TaskType task type
// shell
// api
// TaskDataType
type TaskType uint8

const (
	// Code run code
	Code TaskType = iota + 1
	// API run http req
	API
)

func (tt TaskType) String() string {
	switch tt {
	case Code:
		return "code"
	case API:
		return "api"
	default:
		return "unknow"
	}
}

// RunMode crocodile run mode
// run crocodile as server or client
type RunMode uint8

const (
	// Server run crocodile as server
	Server RunMode = iota + 1
	// Client run crocodile as client
	Client
)

// TaskRespType task type (parent task,master task, child task)
// TODO Rename TaskRunType
type TaskRespType uint8

const (
	// MasterTask task as master run
	MasterTask TaskRespType = iota + 1
	// ParentTask task as a task's parent task run
	ParentTask
	// ChildTask task as a task's child task run
	ChildTask
)

func (tasktype TaskRespType) String() string {
	switch tasktype {
	case MasterTask:
		return "master"
	case ChildTask:
		return "child"
	case ParentTask:
		return "parent"
	default:
		return "unknown"
	}
}

// GetID get task id in post
type GetID struct {
	ID string `json:"id" form:"id" binding:"required,len=18"`
}

// GetName get task name in post
type GetName struct {
	Name string `json:"name" form:"name"  binding:"required,min=1,max=30"`
}

// Common struct
type Common struct {
	ID         string `json:"id" comment:"ID"`
	Name       string `json:"name,omitempty" comment:"名称"`
	CreateTime string `json:"create_time,omitempty" comment:"创建时间"` // 创建时间
	UpdateTime string `json:"update_time,omitempty" comment:"更新时间"` // 最后一次更新时间
	Remark     string `json:"remark" comment:"备注"`                  // 备注
}

// User Struct
type User struct {
	Role      Role     `json:"role"`                               // 用户类型: 1 普通用户 2 管理员 3访客
	Roles     []string `json:"roles"`                              // 管理员
	RoleStr   string   `json:"rolestr,omitempty" comment:"用户类型"`   // 用户类型
	Forbid    bool     `json:"forbid" comment:"禁止用户"`              // 禁止用户登陆
	Password  string   `json:"password,omitempty" comment:"密码"`    // 用户密码
	Email     string   `json:"email" binding:"email" comment:"邮箱"` // 用户邮箱 日后任务的通知信息会发送给此邮件
	WeChat    string   `json:"wechat" comment:"WeChat"`            // wechat id
	DingPhone string   `json:"dingphone" comment:"钉钉"`             // dingding phone
	Slack     string   `json:"slack" comment:"Slack"`              // slack user name
	Telegram  string   `json:"telegram" comment:"Telegram"`        // telegram bot chat id
	Common
}

// RegistryUser data
type RegistryUser struct {
	Name     string `json:"name" binding:"required,max=30"`      // 用户名
	Password string `json:"password" binding:"required,min=8"`   // 用户密码
	Role     Role   `json:"role" binding:"required,min=1,max=3"` // 用户类型: 1 普通用户 2 管理员
	Remark   string `json:"remark" binding:"max=100"`            // 备注
}

// CreateAdminUser first run must be create admin user
type CreateAdminUser struct {
	Name     string `json:"username" binding:"required,max=30"` // 用户名
	Password string `json:"password" binding:"required,min=8"`  // 用户密码
}

// AdminChangeUser struct
type AdminChangeUser struct {
	ID       string `json:"id"  binding:"required,len=18"`       // user id
	Role     Role   `json:"role" binding:"required,min=1,max=3"` // 用户类型: 1 普通用户 2 管理员
	Forbid   bool   `json:"forbid"`                              // 禁止用户: 1 未禁止 2 禁止登陆
	Password string `json:"password"`                            // 用户密码 Common
	Remark   string `json:"remark"`                              // 备注 Common
}

// ChangeUserSelf change self's config
type ChangeUserSelf struct {
	ID        string `json:"id"  binding:"required"`  // user id
	Name      string `json:"name" binding:"required"` // 用户名称
	Email     string `json:"email"`                   // 用户邮箱
	WeChat    string `json:"wechat"`                  // wechat id
	DingPhone string `json:"dingphone"`               // dingding phone
	Telegram  string `json:"telegram"`                // telegram bot chat id
	Password  string `json:"password"`
	Remark    string `json:"remark"`
}

// HostGroup define hostgroup
type HostGroup struct {
	HostsID     []string `json:"addrs" comment:"WorkerIDs"` // 主机host
	CreateByUID string   `json:"create_byuid"`              // 创建人ID
	CreateBy    string   `json:"create_by"`                 // 创建人ID
	Common
}

// CreateHostGroup new hostgroup
type CreateHostGroup struct {
	Name    string   `json:"name" binding:"required,max=30"`
	HostsID []string `json:"addrs"` // 主机host
	Remark  string   `json:"remark" binding:"max=100"`
}

// ChangeHostGroup new hostgroup
type ChangeHostGroup struct {
	ID      string   `json:"id" binding:"required"`
	HostsID []string `json:"addrs"` // 主机host
	Remark  string   `json:"remark" binding:"max=100"`
}

// Host worker host
type Host struct {
	ID                 string   `json:"id" comment:"ID"`
	Addr               string   `json:"addr" comment:"Worker地址"`
	HostName           string   `json:"hostname"`
	Online             bool     `json:"online"`
	Weight             int      `json:"weight"`
	RunningTasks       []string `json:"running_tasks"`
	Version            string   `json:"version"`
	Stop               bool     `json:"stop" comment:"暂停"`
	LastUpdateTimeUnix int64    `json:"last_updatetimeunix"`
	LastUpdateTime     string   `json:"last_updatetime" comment:"更新时间"`
	Remark             string   `json:"remark"`
}

// Task define Task
type Task struct {
	TaskType          TaskType    `json:"task_type" binding:"required"`                 // 任务类型
	TaskData          interface{} `json:"task_data" binding:"required"`                 // 任务数据
	Run               bool        `json:"run" `                                         // 是否可以自动调度  如果为false则只能手动或者被其他任务依赖运行
	ParentTaskIds     []string    `json:"parent_taskids" binding:"max=20"`              // 父任务 运行任务前先运行父任务 以父或子任务运行时 任务不会执行自已的父子任务，防止循环依赖
	ParentRunParallel bool        `json:"parent_runparallel"`                           // 是否以并行运行父任务 0否 1是
	ChildTaskIds      []string    `json:"child_taskids" binding:"max=20"`               // 子任务 运行结束后运行子任务
	ChildRunParallel  bool        `json:"child_runparallel"`                            // 是否以并行运行子任务 否 1是
	CreateBy          string      `json:"create_by"`                                    // 创建人
	CreateByUID       string      `json:"create_byuid"`                                 // 创建人ID
	HostGroup         string      `json:"host_group" `                                  // 主机组
	HostGroupID       string      `json:"host_groupid" binding:"required,len=18"`       // 主机组ID
	Cronexpr          string      `json:"cronexpr" binding:"required,max=1000"`         // 执行任务表达式
	Timeout           int         `json:"timeout" binding:"required,min=-1"`            // 任务超时时间 (s) -1 no limit
	AlarmUserIds      []string    `json:"alarm_userids" binding:"required,max=10"`      // 报警用户 最多十个多个用户
	RoutePolicy       RoutePolicy `json:"route_policy" binding:"required,min=1,max=4"`  // how to select a run worker from hostgroup
	ExpectCode        int         `json:"expect_code"`                                  // expect task return code. if not set 0 or 200
	ExpectContent     string      `json:"expect_content"`                               // expect task return content. if not set do not check
	AlarmStatus       AlarmStatus `json:"alarm_status" binding:"required,min=-2,max=1"` // alarm when task run success or fail or all all:-2 failed: -1 success: 1
	Remark            string      `json:"remark" binding:"max=100"`
}

// AlarmStatus task is alarm
type AlarmStatus int8

const (
	// All will alarm after task run
	All AlarmStatus = -2
	// Fail will alarm after task fail
	Fail AlarmStatus = -1
	// Success will alarm after task success
	Success AlarmStatus = 1
)

func (al AlarmStatus) String() string {
	switch al {
	case All:
		return "All"
	case Fail:
		return "Fail"
	case Success:
		return "Success"
	default:
		return "Unknown"
	}
}

// CreateTask struct
type CreateTask struct {
	GetName
	Task
}

// ChangeTask struct
type ChangeTask struct {
	IDName
	Task
}

// IDName struct
type IDName struct {
	GetID
	GetName
}

// GetTask get task
type GetTask struct {
	//
	TaskType          TaskType    `json:"task_type"`
	TaskTypeDesc      string      `json:"task_typedesc" comment:"任务类型"`
	TaskData          interface{} `json:"task_data" comment:"任务数据"`
	Run               bool        `json:"run" comment:"运行"`
	ParentTaskIds     []string    `json:"parent_taskids"`
	ParentTaskIdsDesc []string    `json:"parent_taskidsdesc" comment:"父任务"`
	ParentRunParallel bool        `json:"parent_runparallel" comment:"父任务运行策略"`
	ChildTaskIds      []string    `json:"child_taskids"`
	ChildTaskIdsDesc  []string    `json:"child_taskidsdesc"  comment:"子任务"`
	ChildRunParallel  bool        `json:"child_runparallel" comment:"子任务运行策略"`
	CreateBy          string      `json:"create_by"`
	CreateByUID       string      `json:"create_byuid"`
	HostGroup         string      `json:"host_group" comment:"主机组"`
	HostGroupID       string      `json:"host_groupid"`
	Cronexpr          string      `json:"cronexpr" comment:"CronExpr"`
	Timeout           int         `json:"timeout" comment:"超时时间"`
	AlarmUserIds      []string    `json:"alarm_userids"`
	AlarmUserIdsDesc  []string    `json:"alarm_useridsdesc" comment:"报警用户"`
	RoutePolicy       RoutePolicy `json:"route_policy"`
	RoutePolicyDesc   string      `json:"route_policydesc" comment:"路由策略"`
	ExpectCode        int         `json:"expect_code"  comment:"期望返回码"`
	ExpectContent     string      `json:"expect_content" comment:"期望返回内容"`
	AlarmStatus       AlarmStatus `json:"alarm_status"`
	AlarmStatusDesc   string      `json:"alarm_statusdesc" comment:"报警策略"`
	Common
}

// RoutePolicy set a task hot to select run worker
type RoutePolicy uint8

const (
	// Random get host by random
	Random RoutePolicy = iota + 1
	// RoundRobin get host by order
	RoundRobin
	// Weight get host by host weight
	Weight
	// LeastTask get host by host LeastTask
	LeastTask
)

func (r RoutePolicy) String() string {
	switch r {
	case Random:
		return "Random"
	case RoundRobin:
		return "RoundRobin"
	case Weight:
		return "Weight"
	case LeastTask:
		return "LeastTask"
	default:
		return "Unknown"
	}
}

// Trigger return how to trigger run task
type Trigger uint8

const (
	// Auto cron run task
	Auto Trigger = iota + 1
	// Manual trigger run task
	Manual
)

func (t Trigger) String() string {
	switch t {
	case Auto:
		return "自动触发"
	case Manual:
		return "手动触发"
	default:
		return "UnKnown"
	}
}

// RunTask running task message
type RunTask struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Cronexpr     string  `json:"cronexpr"`
	StartTimeStr string  `json:"start_timestr"`
	StartTime    int64   `json:"start_time"` // use ms,
	RunTime      int     `json:"run_time"`   // s
	Trigger      Trigger `json:"trigger"`
	TriggerStr   string  `json:"triggerstr"`
}

// TaskResp run task resp message
type TaskResp struct {
	TaskID      string       `json:"task_id"`
	Task        string       `json:"task"`
	LogData     string       `json:"resp_data"`    // task run log data
	Code        int          `json:"code"`         // return code
	TaskType    TaskRespType `json:"task_type"`    // 1 主任务 2 父任务 3 子任务
	TaskTypeStr string       `json:"task_typestr"` // 1 主任务 2 父任务 3 子任务
	RunHost     string       `json:"run_host"`     // task run host
	Status      string       `json:"status"`       // task status finish,fail, cancel
}

// Log task log
type Log struct {
	Name           string       `json:"name"`                 // task log
	RunByTaskID    string       `json:"runby_taskid"`         // run taskid
	StartTime      int64        `json:"start_time"`           // ms
	StartTimeStr   string       `json:"start_timestr"`        //
	EndTime        int64        `json:"end_time"`             // ms
	EndTimeStr     string       `json:"end_timestr"`          //
	TotalRunTime   int          `json:"total_runtime"`        // ms
	Status         int          `json:"status"`               // 任务运行结果 -1 失败 1 成功
	TaskResps      []*TaskResp  `json:"task_resps,omitempty"` // 任务执行过程日志
	Trigger        Trigger      `json:"trigger"`              // 任务触发
	Triggerstr     string       `json:"trigger_str"`          // 任务触发
	ErrCode        int          `json:"err_code"`             // err code
	ErrMsg         string       `json:"err_msg"`              // 错误原因
	ErrTasktype    TaskRespType `json:"err_tasktype"`         // err task type
	ErrTaskTypeStr string       `json:"err_tasktypestr"`      // 1 主任务 2 父任务 3 子任务
	ErrTaskID      string       `json:"err_taskid"`           // task failed id
	ErrTask        string       `json:"err_task"`             // task failed id
}

// Cleanlog data
type Cleanlog struct {
	GetName
	PreDay int64 `json:"preday"` // preday几天前的日志 0 为全部日志
}

// Query recv url query params
type Query struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

// KlOption vue el-select
type KlOption struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Online int    `json:"online,omitempty"` // online: 1 offline: -1
}

// TaskStatusTree real task tree
type TaskStatusTree struct {
	Name         string       `json:"name"`
	ID           string       `json:"id,omitempty"`
	Status       string       `json:"status"`
	TaskType     TaskRespType `json:"tasktype"`
	TaskRespData string       `json:"taskresp_data,omitempty"`
	// RunHost      string            `json:"runhost,omitempty"`
	Children []*TaskStatusTree `json:"children,omitempty"`
}

// GetTasksTreeStatus return a slice
func GetTasksTreeStatus() []*TaskStatusTree {
	retTasksStatus := make([]*TaskStatusTree, 0, 3)
	parentTasksStatus := &TaskStatusTree{
		Name:     "ParentTasks",
		Status:   TsNoData.String(),
		Children: make([]*TaskStatusTree, 0),
	}

	mainTaskStatus := &TaskStatusTree{
		// Name:   task.name,
		// ID:     taskid,
		TaskType: MasterTask,
		Status:   TsNoData.String(),
	}

	childTasksStatus := &TaskStatusTree{
		Name:     "ChildTasks",
		Status:   TsNoData.String(),
		Children: make([]*TaskStatusTree, 0),
	}

	retTasksStatus = append(retTasksStatus,
		parentTasksStatus,
		mainTaskStatus,
		childTasksStatus)
	return retTasksStatus
}

// TaskStatus task run status
type TaskStatus uint

const (
	// TsWait task is waiting pre task is running
	TsWait TaskStatus = iota + 1
	// TsRun tassk is running
	TsRun
	// TsFinish task is run finish
	TsFinish
	// TsFail task run fail
	TsFail
	// TsCancel task is cancel ,because pre task is run fail
	TsCancel
	// TsNoData parenttasks or childtasks no task
	TsNoData
)

func (t TaskStatus) String() string {
	switch t {
	case TsWait:
		return "wait"
	case TsRun:
		return "run"
	case TsFinish:
		return "finish"
	case TsFail:
		return "fail"
	case TsCancel:
		return "cancel"
	case TsNoData:
		return "nodata"
	default:
		return "unknown"
	}
}

// OperateLog openrate log
type OperateLog struct {
	UID         string   `json:"user_id"`      // 修改人ID
	UserName    string   `json:"user_name"`    // 修改人姓名
	Role        Role     `json:"user_role"`    // 用户类型
	Method      string   `json:"method"`       // 新增 修改删除
	Module      string   `json:"module"`       // 修改模块 任务 主机组 主机 用户
	ModuleName  string   `json:"module_name"`  // 修改的对象名称
	OperateTime string   `json:"operate_time"` // 修改时间
	Desc        string   `json:"desc"`         // 描述
	Columns     []Column `json:"columns"`      // 修改的字段及新旧值
}

// Column change column old and new value
type Column struct {
	Name     string      `json:"name"`      // 修改的字段
	OldValue interface{} `json:"old_value"` // 修改前的旧值
	NewValue interface{} `json:"new_value"` // 修改后的新值
}

// NotifyType notify type
type NotifyType uint8

const (
	// TaskNotify 任务通知
	TaskNotify NotifyType = iota + 1
	// UpgradeNotify 升级提醒
	UpgradeNotify
	// ReviewReq 审核请求
	ReviewReq
)

func (nt NotifyType) String() string {
	switch nt {
	case TaskNotify:
		return "任务通知"
	case UpgradeNotify:
		return "新版本发布"
	// case ReviewReq:
	// 	return "审核请求" // zaicontent中点击url到任务列表
	default:
		return "Unknow"
	}
}

// Notify notify msg
type Notify struct {
	ID             int        `json:"id"`
	NotifyType     NotifyType `json:"notify_type"` // 通知类型
	NotifyTypeDesc string     `json:"notify_typedesc"`
	NotifyUID      string     `json:"notify_uid,omitempty"` // 通知用户
	Title          string     `json:"title"`                // 标题
	Content        string     `json:"content"`              // 通知内容
	NotifyTime     int64      `json:"notify_time"`
	NotifyTimeDesc string     `json:"notify_timedesc"`
}
