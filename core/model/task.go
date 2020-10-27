package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/tasktype"
	"github.com/labulaka521/crocodile/core/utils/define"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateTask create task
func CreateTask(ctx context.Context, id, name string, tasktype define.TaskType, taskData interface{}, run bool,
	parentTaskIds []string, parentRunParallel bool, childTaskIds []string, childRunParallel bool,
	cronExpr string, timeout int, alarmUserIds []string, routePolicy define.RoutePolicy, expectCode int,
	expectContent string, alarmStatus define.AlarmStatus, createByID, hostGroupID, remark string) error {
	createsql := `INSERT INTO crocodile_task 
					(id,
					name,
					taskType,
					taskData,
					run,
					parentTaskIds,
					parentRunParallel,
					childTaskIds,
					childRunParallel,
					cronExpr,
					timeout,
					alarmUserIds,
					routePolicy,
					expectCode,
					expectContent,
					alarmStatus,
					createByID,
					hostGroupID,
					remark,
					createTime,
					updateTime)
				VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, createsql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	createTime := time.Now().Unix()
	taskdata, _ := json.Marshal(taskData)
	_, err = stmt.ExecContext(ctx,
		id,
		name,
		tasktype,
		fmt.Sprintf("%s", taskdata),
		run,
		strings.Join(parentTaskIds, ","),
		parentRunParallel,
		strings.Join(childTaskIds, ","),
		childRunParallel,
		cronExpr,
		timeout,
		strings.Join(alarmUserIds, ","),
		routePolicy,
		expectCode,
		expectContent,
		alarmStatus,
		createByID,
		hostGroupID,
		remark,
		createTime,
		createTime,
	)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}

	return nil
}

// ChangeTask change task
func ChangeTask(ctx context.Context, id string, run bool, tasktype define.TaskType, taskData interface{},
	parentTaskIds []string, parentRunParallel bool, childTaskIds []string, childRunParallel bool,
	cronExpr string, timeout int, alarmUserIds []string, routePolicy define.RoutePolicy, expectCode int,
	expectContent string, alarmStatus define.AlarmStatus, hostGroupID, remark string) error {
	changesql := `UPDATE crocodile_task 
					SET hostGroupID=?,
						run=?,
						taskType=?,
						taskData=?,
						parentTaskIds=?,
						parentRunParallel=?,
						childTaskIds=?,
						childRunParallel=?,
						cronExpr=?,
						timeout=?,
						alarmUserIds=?,
						routePolicy=?,
						expectCode=?,
						expectContent=?,
						alarmStatus=?,
						remark=?,
						updateTime=?
					WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, changesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	updateTime := time.Now().Unix()
	taskdata, _ := json.Marshal(taskData)

	_, err = stmt.ExecContext(ctx,
		hostGroupID,
		run,
		tasktype,
		fmt.Sprintf("%s", taskdata),
		strings.Join(parentTaskIds, ","),
		parentRunParallel,
		strings.Join(childTaskIds, ","),
		childRunParallel,
		cronExpr,
		timeout,
		strings.Join(alarmUserIds, ","),
		routePolicy,
		expectCode,
		expectContent,
		alarmStatus,
		remark,
		updateTime,
		id,
	)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// DeleteTask delete task
func DeleteTask(ctx context.Context, id string) error {
	deletesql := `DELETE FROM crocodile_task WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, deletesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// TaskIsUse check a task is other task's parent task ids or child task
func TaskIsUse(ctx context.Context, taskid string) (int, error) {
	querysql := `select count(*) from crocodile_task WHERE id!=? AND (parentTaskIds LIKE ? OR childTaskIds LIKE ?) `
	conn, err := db.GetConn(ctx)
	if err != nil {
		return 0, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, querysql)
	if err != nil {
		return 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	var count int
	likequery := "%" + taskid + "%"
	err = stmt.QueryRowContext(ctx, taskid, likequery, likequery).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("stmt.QueryRowContext failed: %w", err)
	}
	return count, nil
}

// GetTasks get all tasks
func GetTasks(ctx context.Context, offset, limit int, name, presearchname, createby string) ([]define.GetTask, int, error) {
	return getTasks(ctx, nil, name, offset, limit, true, presearchname, createby)
}

// GetTaskByID get task by id
func GetTaskByID(ctx context.Context, id string) (*define.GetTask, error) {
	tasks, _, err := getTasks(ctx, []string{id}, "", 0, 0, true, "", "")
	if err != nil {
		return nil, err
	}
	if len(tasks) != 1 {
		err = define.ErrNotExist{Value: id}
		log.Error("get taskid failed", zap.Error(err))
		return nil, err
	}
	return &tasks[0], nil
}

// GetTaskByName get task by id
func GetTaskByName(ctx context.Context, name string) (*define.GetTask, error) {
	tasks, _, err := getTasks(ctx, nil, name, 0, 0, true, "", "")
	if err != nil {
		return nil, err
	}
	if len(tasks) != 1 {
		err = define.ErrNotExist{Value: name}
		log.Error("get taskname failed", zap.Error(err))
		return nil, err
	}
	return &tasks[0], nil
}

// getTasks get takls by id
func getTasks(ctx context.Context,
	ids []string,
	name string,
	offset,
	limit int,
	first bool, /*Preventing endless loops*/
	presearchname,
	createbyid string) ([]define.GetTask, int, error) {
	getsql := `SELECT t.id,
					t.name,
					t.tasktype,
					t.taskdata,
					t.run,
					t.parentTaskIds,
					t.parentRunParallel,
					t.childTaskIds,
					t.childRunParallel,
					t.cronExpr,
					t.timeout,
					t.alarmUserIds,
					t.routePolicy,
					t.expectCode,
					t.expectContent,
					t.alarmStatus,
					u.name,
					t.createByID,
					hg.name,
					t.hostGroupID,
					t.remark,
					t.createTime,
					t.updateTime
				FROM 
					crocodile_task as t,crocodile_user as u,crocodile_hostgroup as hg 
				WHERE
					t.createByID = u.id AND t.hostGroupID = hg.id`
	args := []interface{}{}
	var count int
	if len(ids) != 0 {
		getsql += " AND ("
		querys := []string{}
		for _, id := range ids {
			querys = append(querys, "t.id=?")
			args = append(args, id)
		}
		getsql += strings.Join(querys, " OR ")
		getsql += ")"
	}
	if name != "" {
		getsql += " AND t.name=?"
		args = append(args, name)
	}
	if presearchname != "" {
		getsql += " AND t.name LIKE ?"
		args = append(args, presearchname+"%")
	}
	if createbyid != "" {
		getsql += " AND t.createByID=?"
		args = append(args, createbyid)
	}
	tasks := []define.GetTask{}
	if limit > 0 {
		var err error
		count, err = countColums(ctx, getsql, args...)
		if err != nil {
			return tasks, 0, fmt.Errorf("countColums failed: %w", err)
		}
		getsql += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return tasks, 0, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return tasks, 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return tasks, 0, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		t := define.GetTask{}
		var (
			parentTaskIds, childTaskIds string
			createTime, updateTime      int64
			taskdata                    string
			alarmUserids                string
		)

		err = rows.Scan(&t.ID,
			&t.Name,
			&t.TaskType,
			&taskdata,
			&t.Run,
			&parentTaskIds,
			&t.ParentRunParallel,
			&childTaskIds,
			&t.ChildRunParallel,
			&t.Cronexpr,
			&t.Timeout,
			&alarmUserids,
			&t.RoutePolicy,
			&t.ExpectCode,
			&t.ExpectContent,
			&t.AlarmStatus,
			&t.CreateBy,
			&t.CreateByUID,
			&t.HostGroup,
			&t.HostGroupID,
			&t.Remark,
			&createTime,
			&updateTime,
		)
		if err != nil {
			log.Error("rows.Scan ", zap.Error(err))
			continue
		}
		t.CreateTime = utils.UnixToStr(createTime)
		t.UpdateTime = utils.UnixToStr(updateTime)
		t.AlarmUserIds = []string{}
		if alarmUserids != "" {
			t.AlarmUserIds = append(t.AlarmUserIds, strings.Split(alarmUserids, ",")...)
			users, _, err := GetUsers(ctx, t.AlarmUserIds, 0, 0)
			if err != nil {
				log.Error("GetUsers ids failed", zap.Strings("uids", t.AlarmUserIds))
			}
			t.AlarmUserIdsDesc = make([]string, 0, len(t.AlarmUserIds))
			for _, user := range users {
				t.AlarmUserIdsDesc = append(t.AlarmUserIdsDesc, user.Name)
			}
		}
		t.ParentTaskIds = []string{}
		t.ParentTaskIdsDesc = []string{}
		if parentTaskIds != "" {
			t.ParentTaskIds = append(t.ParentTaskIds, strings.Split(parentTaskIds, ",")...)
			if first {
				ptasks, _, err := getTasks(ctx, t.ParentTaskIds, "", 0, 0, false, "", "")
				if err != nil {
					log.Error("getTasks failed", zap.Error(err))
				}
				for _, task := range ptasks {
					t.ParentTaskIdsDesc = append(t.ParentTaskIdsDesc, task.Name)
				}
			}

		}
		t.ChildTaskIds = []string{}
		t.ChildTaskIdsDesc = []string{}
		if childTaskIds != "" {
			t.ChildTaskIds = append(t.ChildTaskIds, strings.Split(childTaskIds, ",")...)
			if first {
				ctasks, _, err := getTasks(ctx, t.ChildTaskIds, "", 0, 0, false, "", "")
				if err != nil {
					log.Error("getTasks failed", zap.Error(err))
				}
				for _, task := range ctasks {
					t.ChildTaskIdsDesc = append(t.ChildTaskIdsDesc, task.Name)
				}
			}

		}
		req := pb.TaskReq{
			TaskType: int32(t.TaskType),
			TaskData: []byte(taskdata),
		}
		t.TaskData, err = tasktype.GetDataRun(&req)
		if err != nil {
			log.Error("GetDataRun failed", zap.Any("type", t.TaskType), zap.Error(err))
			continue
		}
		t.RoutePolicyDesc = t.RoutePolicy.String()
		t.TaskTypeDesc = t.TaskType.String()
		t.AlarmStatusDesc = t.AlarmStatus.String()
		t.TaskTypeDesc = t.TaskType.String()

		tasks = append(tasks, t)
	}
	return tasks, count, nil
}

// CreateTaskv2 create task v2
func CreateTaskv2(ctx context.Context, task *Task) (string, error) {
	err := gormdb.WithContext(ctx).Create(task).Error
	if err != nil {
		return "", fmt.Errorf("create task failed: %w", err)
	}
	return task.ID, nil
}

// ChangeTaskv2 change task v2
func ChangeTaskv2(ctx context.Context, task *Task) error {
	result := gormdb.WithContext(ctx).Model(&Task{}).Where("id = ?", task.ID).Updates(task)
	if result.Error != nil {
		return fmt.Errorf("update task %s failed: %w", task.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return define.ErrNotExist{Type: "task id", Value: task.ID}
	}
	return nil
}

// DeleteTaskv2 delete task
func DeleteTaskv2(ctx context.Context, taskid string) error {
	result := gormdb.WithContext(ctx).Model(&Task{}).Delete("id = ?", taskid)
	if result.Error != nil {
		return fmt.Errorf("delete task %s failed: %w", taskid, result.Error)
	}
	if result.RowsAffected == 0 {
		return define.ErrNotExist{Type: "task id", Value: taskid}
	}
	return nil
}

// TaskIsUsev2 check a task is other task's parent task ids or child task when delete task
func TaskIsUsev2(ctx context.Context, taskid string) (int64, error) {
	var count int64
	gormdb = gormdb.Debug()
	err := gormdb.WithContext(ctx).Model(&Task{}).Where("child_task_ids LIKE @likeid OR parent_task_ids LIKE @likeid", sql.Named("likeid", "%"+taskid+"%")).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("find taskid %s failed: %w", taskid, err)
	}
	return count, nil
}

// GetTasksv2 get all tasks
func GetTasksv2(ctx context.Context, offset, limit int, presearchtaskname, createUID string) ([]*Task, int64, error) {
	var (
		count int64
		tasks = []*Task{}
	)
	db := gormdb.WithContext(ctx).Where(&Task{CreateUID: createUID})
	if presearchtaskname != "" {
		db = db.Where("name LIKE ?", presearchtaskname+"%")
	}
	err := db.Find(&tasks).Count(&count).Limit(limit).Offset(offset).Error
	if err != nil {
		return nil, 0, fmt.Errorf("get tasks failed: %w", err)
	}
	return tasks, count, nil
}

// GetTaskByIDv2 get task by id
func GetTaskByIDv2(ctx context.Context, taskid string) (*Task, error) {
	var task = Task{}
	result := gormdb.WithContext(ctx).Where("id = ?", taskid).First(&task)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("find taskkid %s failed: %w", taskid, result.Error)
	}
	if result.Error == gorm.ErrRecordNotFound {
		return nil, define.ErrNotExist{Type: "task id", Value: taskid}
	}
	return &task, nil
}

// GetTaskByNamev2 get task by name
func GetTaskByNamev2(ctx context.Context, taskname string) (*Task, error) {
	var task = Task{}
	result := gormdb.WithContext(ctx).Where("name = ?", taskname).First(&task)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("find taskkname %s failed: %w", taskname, result.Error)
	}
	if result.Error == gorm.ErrRecordNotFound {
		return nil, define.ErrNotExist{Type: "task name", Value: taskname}
	}
	return &task, nil
}
