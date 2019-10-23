package define

type Role uint8

const (
	GuestUser  Role = iota // 访客 没有任和操作权限，并且不可以看见用户管理的界面
	NormalUser             // 普通用户 只对自已创建的主机或者主机组具有操作权限
	AdminUser              // 管理员 具有所有操作
)

func GetUserByRole(r Role) string {
	switch r {
	case AdminUser:
		return "Admin"
	case NormalUser:
		return "Normal"
	case GuestUser:
		return "Guest"
	default:
		return ""
	}
}

// 定义结构体
type common struct {
	Id         int64  `json:"id"`
	Name       string `json:"name" binding:"required"`
	CreateTime string `json:"create_time"` // 创建时间
	UpdateTime string `json:"update_time"` // 最后一次更新时间
	Remark     string `json:"remark"`      // 备注
}

// 用户
type User struct {
	Role     Role   `json:"role" binding:"gte=0,lte=2"`     // 用户类型: 2 管理员 1 普通用户 0 访客
	Forbid   int    `json:"forbid"`                         // 禁止用户
	Email    string `json:"email" binding:"required,email"` // 用户邮箱 日后任务的通知信息会发送给此邮件
	Password string `json:"password,omitempty"`
	common
}

// 主机组
type HostGroup struct {
	Hosts    []int64 `json:"hosts" binding:"required"` // 主机
	CreateBy int     `json:"create_by"`                // 创建人ID
	common
}

// 主机信息
// 定时更新主机信息
type Host struct {
	common
	IP       string  `json:"ip"`      // 主机IP
	Port     int     `json:"port"`    // 运行端口
	Priority string  `json:"priorty"` // 主机执行优先级
	Online   int     `json:"online"`  // 主机是否在线
	Tasks    []int64 `json:"tasks"`   // 运行的任务  会依照优先级和任务数的多少来给执行端分配worker
}

// 任务
type Task struct {
	Timeout           int   `json:"timeout"`             // 0 为不启动超时控制 单位秒
	Run               int   `json:"run"`                 // 0 为不能运行 1 为可以运行
	ParentTaskIds     []int `json:"parent_task_ids"`     // 父任务 运行任务前先运行父任务 以父或子任务运行时 任务不会执行自已的父子任务，防止循环依赖
	ParentRunParallel int   `json:"parent_run_parallel"` // 是否以并行运行父任务 0否 1是
	ChildTaskIds      []int `json:"child_task_ids"`      // 子任务 运行结束后运行子任务
	ChildRunParallel  int   `json:"child_run_parallel"`  // 是否以并行运行子任务 0否 1是
	RunByTask         int   `json:"run_by_task"`         // 通过其他的任务调用运行 如果是被其他任务依赖而运行就不会运行此任务的父子任务 1 为 true
	RunTime           int   `json:"run_time"`            // 运行次数 默认为0则不限制，如果设置大于0，等成功调度这么多次数后，任务会停止
	AlarmTotal        int   `json:"alarm_total"`         // 任务失败报警次数，默认为0则不限制，如果设置大于0，等报警这么多次数后，任务会自动停止
	CreateBy          int   `json:"create_by"`           // 创建人ID
	common
}

// 日志
type Log struct {
	TaskId int `json:""`
}
