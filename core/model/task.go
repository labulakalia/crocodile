package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/schedule"
	"github.com/labulaka521/crocodile/core/tasktype"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

// 执行计划：
//

func CreateTask(ctx context.Context, t *define.Task) error {
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
					runtime,
					alarmTotal,
					alarmUser,
					autoSwitch,
					createByID,
					hostGroupID,
					remark,
					createTime,
					updateTime)
				VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}

	stmt, err := conn.PrepareContext(ctx, createsql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	createTime := time.Now().Unix()
	taskdata, _ := json.Marshal(t.TaskData)
	_, err = stmt.ExecContext(ctx,
		t.Id,
		t.Name,
		t.TaskType,
		fmt.Sprintf("%s", taskdata),
		t.Run,
		strings.Join(t.ParentTaskIds, ","),
		t.ParentRunParallel,
		strings.Join(t.ChildTaskIds, ","),
		t.ChildRunParallel,
		t.CronExpr,
		t.Timeout,
		t.RunTime,
		t.AlarmTotal,
		strings.Join(t.AlarmUser, ","),
		t.AutoSwitch,
		t.CreateByUId,
		t.HostGroupId,
		t.Remark,
		createTime,
		createTime,
	)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}

	return nil
}

func ChangeTask(ctx context.Context, t *define.Task) error {
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
						runtime=?,
						alarmTotal=?,
						alarmUser=?,
						autoSwitch=?,
						remark=?,
						updateTime=?
					WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	stmt, err := conn.PrepareContext(ctx, changesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	updateTime := time.Now().Unix()
	taskdata, _ := json.Marshal(t.TaskData)
	_, err = stmt.ExecContext(ctx,
		t.HostGroupId,
		t.Run,
		t.TaskType,
		fmt.Sprintf("%s", taskdata),
		strings.Join(t.ParentTaskIds, ","),
		t.ParentRunParallel,
		strings.Join(t.ChildTaskIds, ","),
		t.ChildRunParallel,
		t.CronExpr,
		t.Timeout,
		t.RunTime,
		t.AlarmTotal,
		strings.Join(t.AlarmUser, ","),
		t.AutoSwitch,
		t.Remark,
		updateTime,
		t.Id,
	)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func DeleteTask(ctx context.Context, id string) error {
	deletesql := `DELETE FROM crocodile_task WHERE id=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}
	stmt, err := conn.PrepareContext(ctx, deletesql)
	if err != nil {
		return errors.Wrap(err, "conn.PrepareContext")
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return errors.Wrap(err, "stmt.ExecContext")
	}
	return nil
}

func GetTasks(ctx context.Context) ([]define.Task, error) {
	return getTasks(ctx, "")
}

func GetTaskByID(ctx context.Context, id string) (*define.Task, error) {
	tasks, err := getTasks(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(tasks) != 1 {
		return nil, errors.Errorf("getTasks id failed: %v", err)
	}
	return &tasks[0], nil
}

// 获取任务
func getTasks(ctx context.Context, id string) ([]define.Task, error) {
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
					t.runtime,
					t.alarmTotal,
					t.alarmUser,
					t.autoSwitch,
					u.name,
					t.createByID,
					hg.name,
					t.hostGroupID,
					t.remark,
					t.createTime,
					t.updateTime 
				FROM crocodile_task as t,crocodile_user as u,crocodile_hostgroup as hg
				WHERE t.createByID == u.id AND t.hostGroupID = hg.id`
	args := []interface{}{}
	if id != "" {
		getsql += " t.id = ?"
		args = append(args, id)
	}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "db.GetConn")
	}
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return nil, errors.Wrap(err, "conn.PrepareContext")
	}
	res := []define.Task{}
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "stmt.QueryContext")
	}
	for rows.Next() {
		t := define.Task{}
		var (
			parentTaskIds, childTaskIds string
			createTime, updateTime      int64
			taskdata                    string
			alarmUser                   string
		)

		err = rows.Scan(&t.Id,
			&t.Name,
			&t.TaskType,
			&taskdata,
			&t.Run,
			&parentTaskIds,
			&t.ParentRunParallel,
			&childTaskIds,
			&t.ChildRunParallel,
			&t.CronExpr,
			&t.Timeout,
			&t.RunTime,
			&t.AlarmTotal,
			&alarmUser,
			&t.AutoSwitch,
			&t.CreateBy,
			&t.CreateByUId,
			&t.HostGroup,
			&t.HostGroupId,
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
		t.AlarmUser = []string{}
		if alarmUser != "" {
			t.AlarmUser = append(t.AlarmUser, strings.Split(alarmUser, ",")...)
		}
		switch t.TaskType {
		case define.Shell:
			pro := tasktype.DataShell{}
			err = json.Unmarshal([]byte(taskdata), &pro)
			if err != nil {
				continue
			}
			t.TaskData = pro
		case define.Api:
			api := tasktype.DataApi{}
			err = json.Unmarshal([]byte(taskdata), &api)
			if err != nil {
				continue
			}
			t.TaskData = api
		}

		res = append(res, t)
	}
	return res, nil
}
