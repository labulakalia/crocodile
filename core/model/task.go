package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	pb "github.com/labulaka521/crocodile/core/proto"
	"github.com/labulaka521/crocodile/core/tasktype"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

// CreateTask create task
func CreateTask(ctx context.Context, id, name string, tasktype define.TaskType, taskData interface{},
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
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, createsql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	createTime := time.Now().Unix()
	taskdata, _ := json.Marshal(taskData)
	_, err = stmt.ExecContext(ctx,
		id,
		name,
		tasktype,
		fmt.Sprintf("%s", taskdata),
		true,
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
		return errors.Wrap(err, "stmt.ExecContext")
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
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, changesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
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
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

// DeleteTask delete task
func DeleteTask(ctx context.Context, id string) error {
	deletesql := `DELETE FROM crocodile_task WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, deletesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

// GetTasks get all tasks
func GetTasks(ctx context.Context, offset, limit int) ([]define.GetTask, error) {
	return getTasks(ctx, "", offset, limit)
}

// GetTaskByID get task by id
func GetTaskByID(ctx context.Context, id string) (*define.GetTask, error) {
	tasks, err := getTasks(ctx, id, 0, 0)
	if err != nil {
		return nil, err
	}
	if len(tasks) != 1 {
		err := fmt.Errorf("can find task %s, map be it has deleted", id)
		return nil, errors.Errorf("getTasks id failed: %v", err)
	}
	return &tasks[0], nil
}

// getTasks get takls by id
func getTasks(ctx context.Context, id string, offset, limit int) ([]define.GetTask, error) {
	getsql := `SELECT 
					t.id,
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
				WHERE t.createByID == u.id AND t.hostGroupID = hg.id`
	args := []interface{}{}
	if id != "" {
		getsql += " AND t.id = ?"
		args = append(args, id)
	}
	if limit > 0 {
		getsql += " LIMIT ? OFFSET ?"
		args = append(args, limit, offset)
	}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "db.GetConn")
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return nil, errors.Wrap(err, "conn.PrepareContext")
	}
	defer stmt.Close()
	res := []define.GetTask{}
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "stmt.QueryContext")
	}
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
		}
		t.ParentTaskIds = []string{}
		if parentTaskIds != "" {
			t.ParentTaskIds = append(t.ParentTaskIds, strings.Split(parentTaskIds, ",")...)
		}
		t.ChildTaskIds = []string{}
		if childTaskIds != "" {
			t.ChildTaskIds = append(t.ChildTaskIds, strings.Split(childTaskIds, ",")...)
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
		
		res = append(res, t)
	}
	return res, nil
}
