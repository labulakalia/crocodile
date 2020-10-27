package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
	"gorm.io/gorm"
)

const dbPrefix = "test_crocodile_"

// Model custom common model
type Model struct {
	ID        string         `gorm:"type:CHAR(18);primaryKey;index" json:"id"`
	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook generate snake id
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.ID = utils.GetID()
	return nil
}

// Host orm model
type Host struct {
	Model
	Addr          string `gorm:"type:varchar(25);not null;index" json:"addr"`
	HostName      string `gorm:"type:varchar(100);not null" json:"hostname"`
	CountRunTasks int    `gorm:"type:integer;not null" json:"count_run_tasks"`
	Online        bool   `gorm:"-" json:"online"`
	Weight        int32  `gorm:"type:integer;not null;default:100" json:"weight"`
	Stop          bool   `gorm:"type:bool;not null;default:false" json:"stop"`
	Version       string `gorm:"type:varchar(10);size:10;not null;" json:"version"`
	Remark        string `gorm:"type:varchar(100);size:100;not null;default:''" json:"remark"`
}

var maxworklive time.Duration = 20 * time.Second

// AfterFind change host online status
func (h *Host) AfterFind(tx *gorm.DB) error {
	h.Online = time.Now().Sub(h.UpdatedAt) < maxworklive
	return nil
}

// TableName custom Host table name
func (h Host) TableName() string {
	return dbPrefix + "host"
}

// HostGroup orm Model
type HostGroup struct {
	Model
	Name       string `gorm:"type:varchar(30);not null;uniqueindex" json:"name"`
	CreateID   string `gorm:"type:char(18);not null" json:"create_id"`
	CreateName string `gorm:"type:varchar(30);not null" json:"create_name"`
	Hosts      IDs    `gorm:"type:varchar(360);not null;default ''" json:"hosts"`
	Remark     string `gorm:"type:varchar(100);not null;default:''" json:"remark"`
}

// IDs custom gorm type
type IDs []string

// Scan impl sql.Scanner interface
func (hids *IDs) Scan(value interface{}) error {
	ids, ok := value.(string)
	if !ok {
		return fmt.Errorf("Scan value must be string, but get type %T", value)
	}
	a := IDs{}
	if ids == "" {
		hids = &a
		return nil
	}
	for i := 0; i < len(ids); i += 18 {
		*hids = append(*hids, ids[i:i+18])
	}
	return nil
}

// Value impl driver.Valuer interface
func (hids IDs) Value() (driver.Value, error) {
	if len(hids) == 0 {
		return "", nil
	}
	for _, id := range hids {
		if len(id) != 18 {
			return nil, fmt.Errorf("%s is not valid id", id)
		}
	}
	return strings.Join(hids, ""), nil
}

// TableName custom HostGroup table name
func (h HostGroup) TableName() string {
	return dbPrefix + "hostgroup"
}

// Log orm model
type Log struct {
	ID          int                 `gorm:"primarykey;autoIncrement"`
	TaskName    string              `gorm:"type:varchar(30);not null;index" json:"task_name"`
	TaskID      string              `gorm:"type:char(18);not null" json:"task_id"`
	StartTime   time.Time           `gorm:"not null;index:idx_s_t" json:"start_time"`
	EndTime     time.Time           `gorm:"not null" json:"end_time"`
	Status      int                 `gorm:"type:tinyint;not null;default 0" json:"status"`
	TaskResps   TaskResps           `gorm:"type:mediumtext" json:"task_resps"`
	TriggerType define.Trigger      `gorm:"type:tinyint;not null;default 0" json:"trigger_type"`
	ErrCode     int                 `gorm:"type:integer;default 0;not null" json:"err_code"`
	ErrMsg      string              `gorm:"type:mediumtext;not null" json:"err_msg"`
	ErrTaskType define.TaskRespType `gorm:"type:integer;not null;default 0" json:"err_tasktype"`
	ErrTaskID   string              `gorm:"type:varchar(19);not null;default ''" json:"err_taskid"`
}

// TableName custom HostGroup table name
func (h Log) TableName() string {
	return dbPrefix + "log"
}

// TaskResp task run log d
type TaskResp struct {
	TaskID   string              `json:"task_id"`
	TaskName string              `json:"task_name"`
	LogData  string              `json:"resp_data"` // task run log data
	Code     int                 `json:"code"`      // return code
	TaskType define.TaskRespType `json:"task_type"` // 1 主任务 2 父任务 3 子任务
	RunHost  string              `json:"run_host"`  // task run host
	Status   string              `json:"status"`    // task status finish,fail, cancel
}

// TaskResps task resp log data
type TaskResps []TaskResp

// Scan impl sql.Scanner interface
func (t *TaskResps) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan value must be []byte, but get type %T", value)
	}
	err := json.Unmarshal(data, t)
	if err != nil {
		return fmt.Errorf("can unmarshal to task resps %w", err)
	}
	return nil
}

// Value impl driver.Valuer interface
func (t TaskResps) Value() (driver.Value, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("marshal taskresps failed: %w", err)
	}
	return data, nil
}

// Notify orm model
type Notify struct {
	ID       uint              `gorm:"primarykey;autoIncrement"`
	Type     define.NotifyType `gorm:"type:integer;not null;default 0" json:"type"`
	TypeDesc string            `gorm:"-" json:"notify_typedesc"`
	UID      string            `gorm:"type:char(18);not null;default '';index" json:"uid"`
	CreateAt time.Time         `gorm:"not null" json:"create_at"`
	Title    string            `gorm:"type:varchar(30);not null;default ''" json:"title"`
	Content  string            `gorm:"type:varchar(500);not null;default ''" json:"content"`
}

// TableName custom Notify table name
func (n Notify) TableName() string {
	return dbPrefix + "notify"
}

// AfterFind change host online status
func (n *Notify) AfterFind(tx *gorm.DB) error {
	n.TypeDesc = n.Type.String()
	return nil
}

// Operate orm model
type Operate struct {
	ID          uint        `gorm:"primarykey;autoIncrement" json:"id"`
	UID         string      `gorm:"type:char(18);not null;index" json:"uid"`
	Role        define.Role `gorm:"type:integer;not null;default 0" json:"role"`
	Method      string      `gorm:"type:varchar(7);not null;default ''" json:"method"`
	Module      string      `gorm:"type:varchar(10);not null;default ''" json:"module"`
	ModuleName  string      `gorm:"type:varchar(30);not null;default ''" json:"module_name"`
	OperateTime time.Time   `json:"operate_time"`
	Description string      `gorm:"type:varchar(200);" json:"description"`
	Columns     string      `gorm:"type:mediumtext" json:"columns"`
}

// TableName custom Operate table name
func (o Operate) TableName() string {
	return dbPrefix + "operate"
}

// User orm model
type User struct {
	Model
	Name         string      `gorm:"type:varchar(30);not null;default ''" json:"name"`
	HashPassword string      `gorm:"type:varchar(100);not null;default ''" json:"hash_password,omitempty"`
	Role         define.Role `gorm:"type:integer;not null;default 0" json:"role"`
	Forbid       bool        `gorm:"type:bool;not null;default false" json:"forbid"`
	Email        string      `gorm:"type:varchar(30)" json:"email"`
	DingPhone    string      `gorm:"type:varchar(12)" json:"dingphone"`
	Wechat       string      `gorm:"type:varchar(30)" json:"wechat"`
	WechatBot    string      `gorm:"type:varchar(30)" json:"wechat_bot"`
	Telegram     string      `gorm:"type:varchar(100)" json:"telegram"`
	WebHook      string      `gorm:"type:varchar(100)" json:"webhook"`
	Env          Env         `gorm:"type:text" json:"env"`                // 用户的环境变量 用于替换任务数据的隐密字段
	AlartTmpl    string      `gorm:"type:varchar(100)" json:"alarm_tmpl"` // 报警模版
	Remark       string      `gorm:"type:varchar(100);not null;default ''" json:"remark"`
}

// TableName custom User table name
func (u User) TableName() string {
	return dbPrefix + "user"
}

// AfterFind query after change password to empty
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.HashPassword = ""
	return nil
}

// Env task env
type Env map[string]string

// Scan impl sql.Scanner interface
func (e *Env) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan value must be []byte, but get type %T", value)
	}
	err := json.Unmarshal(data, e)
	if err != nil {
		return fmt.Errorf("can unmarshal to task resps %w", err)
	}
	return nil
}

// Value impl driver.Valuer interface
func (e Env) Value() (driver.Value, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("marshal taskresps failed: %w", err)
	}
	return data, nil
}

// Task orm model
type Task struct {
	Model
	Name           string          `gorm:"type:varchar(30);not null" json:"name"`
	TaskType       define.TaskType `gorm:"type:integer;not null" json:"task_type"`
	TaskData       string          `gorm:"type:mediumtext" json:"task_data"`
	Run            bool            `gorm:"type:bool;not null;default false" json:"run"`
	ParentTaskIDs  IDs             `gorm:"type:varchar(360);not null;default ''" json:"parent_task_ids"`
	ParentParallel bool            `gorm:"type:bool;not null;default false" json:"parent_parallel"`
	ChildTaskIDs   IDs             `gorm:"type:varchar(360);not null;default ''" json:"child_task_ids"`
	ChildParallel  bool            `gorm:"type:bool;not null;default false" json:"child_parallel"`
	CreateUID      string          `gorm:"type:char(18);not null;default '';index" json:"create_uid"`
	// CreateName     string             `gorm:"type:varchar(30);not null;default '';index" json:"create_name"`
	HostgroupID string `gorm:"type:char(18);not null;default ''" json:"hostgroup_id"`
	// HostgroupName  string             `gorm:"type:varchar(30);not null;default ''" json:"hostgroup_name"`
	Cronexpr      string             `gorm:"type:varchar(200);not null;default ''" json:"cronexpr"`
	Timeout       int                `gorm:"type:integer;not null;default -1" json:"timeout"`
	AlarmUIDs     IDs                `gorm:"type:varchar(180);not null" json:"alarm_uids"`
	RoutePolicy   define.RoutePolicy `gorm:"type:integer;not null;default 1" json:"route_policy"`
	ExpectCode    int                `gorm:"type:integer;not null;default 0" json:"expect_code"`
	ExpectContent string             `gorm:"type:varchar(500);not null;default ''" json:"expect_content"`
	AlarmStatus   define.AlarmStatus `gorm:"type:integer;not null;default -1" json:"alarm_status"`
	Remark        string             `gorm:"type:varchar(100);not null;default ''" json:"remark"`
}

// TableName custom Task table name
func (t Task) TableName() string {
	return dbPrefix + "task"
}

// TDO
