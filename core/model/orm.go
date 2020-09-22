package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

// Host orm model
type Host struct {
	Model
	Addr          string `gorm:"type:varchar(25);not null;uniqueindex" json:"addr"`
	HostName      string `gorm:"type:varchar(100);not null" json:"hostname"`
	CountRunTasks int    `gorm:"type:integer;not null" json:"count_run_tasks"`
	Online        bool   `gorm:"-" json:"online"`
	Weight        int    `gorm:"type:integer;not null;default:100" json:"weight"`
	Stop          bool   `gorm:"type:bool;not null;default:false" json:"stop"`
	Version       string `gorm:"type:varchar(10);size:10;not null;" json:"version"`
	Remark        string `gorm:"type:varchar(100);size:100;not null;default:''" json:"remark"`
}

// TableName custom Host table name
func (h Host) TableName() string {
	return dbPrefix + "host"
}

// HostGroup orm Model
type HostGroup struct {
	Model
	Name       string `gorm:"type:varchar(30);not null;uniqueindex" json:"name"`
	CreateID   string `gorm:"type:char(18);not null" json:"-"`
	CreateName string `gorm:"-" json:"create_name"`
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
	Name        string              `gorm:"type:varchar(30);not null" json:"name"`
	TaskID      string              `gorm:"type:char(18);not null" json:"task_id"`
	StartTime   time.Time           `gorm:"not null;index:idx_s_t" json:"start_time"`
	EndTime     time.Time           `gorm:"not null" json:"end_time"`
	Status      int                 `gorm:"type:tinyint;not null;default 0" json:"status"`
	TaskResps   TaskResps           `gorm:"type:mediumtext" json:"task_resps"`
	TriggerType define.Trigger      `gorm:"type:tinyint;not null;default 0" json:"trigger_type"`
	ErrCode     int                 `gorm:"type:integer;default 0;not null" json:"err_code"`
	ErrMsg      string              `gorm:"type:mediumtext;not null" json:"err_msg"`
	ErrTaskType define.TaskRespType `gorm:"type:integer;not null;default 0" json:"err_tasktype"`
	ErrTaskID   string              `gorm:"type:char(19);not null;default ''" json:"err_taskid"`
}

// TableName custom HostGroup table name
func (h Log) TableName() string {
	return dbPrefix + "log"
}

// TaskResps task resp log data
type TaskResps []TaskResp

// TaskResp task run log d
type TaskResp struct {
	TaskID   string              `json:"task_id"`
	Task     string              `json:"task"`
	LogData  string              `json:"resp_data"` // task run log data
	Code     int                 `json:"code"`      // return code
	TaskType define.TaskRespType `json:"task_type"` // 1 主任务 2 父任务 3 子任务
	RunHost  string              `json:"run_host"`  // task run host
	Status   string              `json:"status"`    // task status finish,fail, cancel
}

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
		return nil, fmt.Errorf("marshal taskresps failed %w", err)
	}
	return data, nil
}

// Notify orm model
type Notify struct {
	ID       uint              `gorm:"primarykey;autoIncrement"`
	Type     define.NotifyType `gorm:"type:integer;not null;default 0" json:"type"`
	UID      string            `gorm:"type:char(18);not null;default '';index" json:"uid"`
	CreateAt time.Time         `gorm:"not null" json:"create_at"`
	Title    string            `gorm:"type:varchar(30);not null;default ''" json:"title"`
	Content  string            `gorm:"type:varchar(500);not null;default ''" json:"content"`
	IsRead   bool              `gorm:"type:bool;not null;default false" json:"is_read"`
}

// TableName custom Notify table name
func (h Notify) TableName() string {
	return dbPrefix + "notify"
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
func (h Operate) TableName() string {
	return dbPrefix + "operate"
}

// User orm model
type User struct {
	Model
	Name         string      `gorm:"type:varchar(30);not null;default ''" json:"name"`
	HashPassword string      `gorm:"type:varchar(100);not null;default ''" json:"hash_passworf,omitempty"`
	Role         define.Role `gorm:"type:integer;not null;default 0" json:"role"`
	Forbid       bool        `gorm:"type:bool;not null;default false" json:"forbid"`
	Email        string      `gorm:"type:varchar(30)" json:"email"`
	DingPhone    string      `gorm:"type:varchar(12)" json:"dingphone"`
	Wechat       string      `gorm:"type:varchar(30)" json:"wechat"`
	WebHook      string      `gorm:"type:varchar(100)" json:"webhook"`
	Remark       string      `gorm:"type:varchar(100);not null;default ''" json:"remark"`
}

// TableName custom User table name
func (h User) TableName() string {
	return dbPrefix + "user"
}

// Task orm model
type Task struct {
	Model
	Name           string             `gorm:"type:varchar(30);not null" json:"name"`
	TaskType       define.TaskType    `gorm:"type:integer;not null" json:"task_type"`
	TaskData       string             `gorm:"type:mediumtext" json:"task_data"`
	Run            bool               `gorm:"type:bool;not null;default false" json:"run"`
	ParentTaskIDS  IDs                `gorm:"type:varchar(360);not null;default ''" json:"parent_task_ids"`
	ParentParallel bool               `gorm:"type:bool;not null;default false" json:"parent_parallel"`
	ChildTaskIDs   IDs                `gorm:"type:varchar(360);not null;default ''" json:"child_task_ids"`
	ChildParallel  bool               `gorm:"type:bool;not null;default false" json:"child_parallel"`
	CreateUID      string             `gorm:"type:char(18);not null;default '';index" json:"create_uid"`
	HostgroupID    string             `gorm:"type:char(18);not null;default ''" json:"hostgroup_id"`
	Cronexpr       string             `gorm:"type:varchar(200);not null;default ''" json:"cronexpr"`
	Timeout        int                `gorm:"type:integer;not null;default -1" json:"timeout"`
	AlarmUIDs      IDs                `gorm:"type:varchar(180);not null" json:"alarm_uids"`
	RoutePolicy    define.RoutePolicy `gorm:"type:integer;not null;default 1" json:"route_policy"`
	ExpectCode     int                `gorm:"type:integer;not null;default 0" json:"expect_code"`
	ExpectContent  string             `gorm:"type:varchar(100);not null;default ''" json:"expect_content"`
	AlarmPolicy    uint               `gorm:"type:integer;not null;default 2" json:"alarm_policy"`
	Remark         string             `gorm:"type:varchar(100);not null;default ''" json:"remark"`
}

// TableName custom Task table name
func (h Task) TableName() string {
	return dbPrefix + "task"
}
