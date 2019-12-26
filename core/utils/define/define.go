package define

// Admin or Normal User
type Role uint8

const (
	NormalUser Role = iota + 1 // 普通用户 只对自已创建的主机或者主机组具有操作权限
	AdminUser                  // 管理员 具有所有操作
)

// task type
type TaskType uint8

const (
	Shell TaskType = iota + 1
	Api
)

// run crocodile as server or client
type RunMode uint8

const (
	Server RunMode = iota + 1
	Client
)

// task type (parent task,master task, child task)
type TaskRespType uint8

const (
	MasterTask = iota + 1
	ParentTask
	ChildTask
)

func GetUserByRole(r Role) string {
	switch r {
	case AdminUser:
		return "Admin"
	case NormalUser:
		return "Normal"
	default:
		return ""
	}
}

// 定义结构体
type common struct {
	Id         string `json:"id"`
	Name       string `json:"name,omitempty"`
	CreateTime string `json:"createTime,omitempty"` // 创建时间
	UpdateTime string `json:"updateTime,omitempty"` // 最后一次更新时间
	Remark     string `json:"remark"`               // 备注
}

// 用户
type User struct {
	Role     Role   `json:"role" binding:"gte=1,lte=2"`     // 用户类型: 2 管理员 1 普通用户
	Forbid   int    `json:"forbid" binding:"gte=0,lte=1"`   // 禁止用户: 0 未禁止 1 已禁止
	Email    string `json:"email" binding:"required,email"` // 用户邮箱 日后任务的通知信息会发送给此邮件
	Password string `json:"password,omitempty"`
	common
}

// 主机组
type HostGroup struct {
	HostsID []string `json:"addrs"`
	//Hosts       []*Host `json:"hosts,omitempty"`       // WorkerId
	CreateByUId string `json:"createByUId"` // 创建人ID
	CreateBy    string `json:"createBy"`    // 创建人ID
	common
}

// 主机信息
// 定时更新主机信息
type Host struct {
	common
	Addr               string `json:"addr"`     // 主机IP
	HostName           string `json:"hostname"` // 主机名
	Online             int    `json:"online"`   // 主机是否在线 0 not online,1 online
	Version            string `json:"version"`  // 版本号
	LastUpdateTimeUnix int64  `json:"lastUpdatetimeUnix"`
	LastUpdateTime     string `json:"lastUpdatetime"`
}

// 任务
type Task struct {
	TaskType          TaskType    `json:"taskType"`          // 任务类型
	TaskData          interface{} `json:"taskData"`          // 任务数据
	Run               int         `json:"run"`               // 0 为不能运行 1 为可以运行
	ParentTaskIds     []string    `json:"parentTaskIds"`     // 父任务 运行任务前先运行父任务 以父或子任务运行时 任务不会执行自已的父子任务，防止循环依赖
	ParentRunParallel int         `json:"parentRunParallel"` // 是否以并行运行父任务 0否 1是
	ChildTaskIds      []string    `json:"childTaskIds"`      // 子任务 运行结束后运行子任务
	ChildRunParallel  int         `json:"childRunParallel"`  // 是否以并行运行子任务 0否 1是
	//RunByTask         int               `json:"runByTask"`         // 通过其他的任务调用运行 如果是被其他任务依赖而运行就不会运行此任务的父子任务 1 为 true
	CreateBy    string `json:"createBy"`                       // 创建人
	CreateByUId string `json:"createByUId"`                    // 创建人ID
	HostGroup   string `json:"hostGroup"`                      // 执行计划
	HostGroupId string `json:"hostGroupID" binding:"required"` // 主机组ID

	CronExpr     string   `json:"cronExpr" binding:"required"` // 执行任务表达式
	Timeout      int      `json:"timeout"`                     // 任务超时时间
	AlarmUserIds []string `json:"alarmUserIds"`                // 报警用户 多个用户
	AutoSwitch   int      `json:"autoSwitch"`                  // 运行失败自动切换到其他主机上
	common
}

// 日志
type Log struct {
	RunByTaskId  string      `json:"runByTaskId"`
	TaskResps    []*TaskResp `json:"taskResps"`
	StartTimne   int64       `json:"startTime"`
	EndTime      int64       `json:"endTime"`
	TotalRunTime int         `json:"totalRunTime"`
}

// 返回值
type TaskResp struct {
	TaskType   TaskRespType `json:"taskType"` // 1 主任务 2 父任务 3 子任务
	TaskId     string       `json:"taskId"`
	Code       int32        `json:"code"`
	ErrMsg     string       `json:"errMsg"`
	RespData   string       `json:"respdata"`
	WorkerHost string       `json:"workerHost"`
}
