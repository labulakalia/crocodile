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
					alarmUserIds,
					autoSwitch,
					createByID,
					hostGroupID,
					remark,
					createTime,
					updateTime)
				VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
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
		strings.Join(t.AlarmUserIds, ","),
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
						alarmUser=?,
						autoSwitch=?,
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
		strings.Join(t.AlarmUserIds, ","),
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
					t.alarmUserIds,
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
		getsql += " AND t.id = ?"
		args = append(args, id)
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
			alarmUserids                string
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
			&alarmUserids,
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
			log.Error("GetDataRun failed", zap.Any("type", t.TaskType))
			continue
		}
		res = append(res, t)
	}
	return res, nil
}
