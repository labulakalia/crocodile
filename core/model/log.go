package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// SaveLog  save task reps log
func SaveLog(ctx context.Context, l *define.Log) error {
	log.Info("start savelog", zap.Any("tasklog", l))
	savesql := `INSERT INTO crocodile_log
				(name,
				taskid,
				starttime,
				endtime,
				totalruntime,
				status,
				taskresps,
				triggertype,
				errcode,
				errmsg,
				errtasktype,
				errtaskid,
				errtask
			)
			VALUES
			(?,?,?,?,?,?,?,?,?,?,?,?,?)`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, savesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	taskresps, err := json.Marshal(l.TaskResps)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %w", err)
	}
	_, err = stmt.ExecContext(ctx, l.Name, l.RunByTaskID,
		l.StartTime, l.EndTime, l.TotalRunTime,
		l.Status, taskresps, l.Trigger, l.ErrCode, l.ErrMsg,
		l.ErrTasktype, l.ErrTaskID, l.ErrTask)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// GetLog get task resp log by taskid
func GetLog(ctx context.Context, taskname string, status int, offset, limit int) ([]*define.Log, int, error) {
	logs := []*define.Log{}
	getsql := `SELECT 
					name,
					taskid,
					starttime,
					endtime,
					totalruntime,
					status,
					triggertype,
					errcode,
					errmsg,
					errtasktype,
					errtaskid,
					errtask
				FROM 
					crocodile_log`
	args := []interface{}{}
	if taskname != "" {
		args = append(args, taskname)
		getsql += ` WHERE name=?`
	}

	if status != 0 {
		getsql += ` AND status=?`
		args = append(args, status)
	}
	count, err := countColums(ctx, getsql, args...)
	if err != nil {
		return logs, 0, fmt.Errorf("countColums failed: %w", err)
	}
	getsql += ` ORDER BY id DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	conn, err := db.GetConn(ctx)
	if err != nil {
		return logs, 0, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return logs, 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return logs, 0, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}
	for rows.Next() {
		getlog := define.Log{}
		taskrepos := []*define.TaskResp{}
		err = rows.Scan(
			&getlog.Name,
			&getlog.RunByTaskID,
			&getlog.StartTime,
			&getlog.EndTime,
			&getlog.TotalRunTime,
			&getlog.Status,
			&getlog.Trigger,
			&getlog.ErrCode,
			&getlog.ErrMsg,
			&getlog.ErrTasktype,
			&getlog.ErrTaskID,
			&getlog.ErrTask,
		)
		if err != nil {
			log.Error("rows.Scan failed", zap.Error(err))
			continue
		}
		getlog.ErrTaskTypeStr = getlog.ErrTasktype.String()
		getlog.TaskResps = taskrepos
		getlog.StartTimeStr = utils.UnixToStr(getlog.StartTime / 1e3)
		getlog.EndTimeStr = utils.UnixToStr(getlog.EndTime / 1e3)
		getlog.Triggerstr = getlog.Trigger.String()
		logs = append(logs, &getlog)
	}
	return logs, count, nil
}

// GetTreeLog get tree log data
func GetTreeLog(ctx context.Context, id string, startTime int64) ([]*define.TaskStatusTree, error) {
	sqlget := `SELECT taskresps FROM crocodile_log WHERE taskid=? AND starttime=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, sqlget)
	if err != nil {
		return nil, err
	}

	var taskreposbyte []byte
	err = stmt.QueryRowContext(ctx, id, startTime).Scan(&taskreposbyte)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*define.TaskStatusTree, 0), nil
		}
		return nil, err
	}
	taskrepos := []*define.TaskResp{}
	err = json.Unmarshal(taskreposbyte, &taskrepos)
	if err != nil {
		return nil, err
	}
	retTasksStatus := define.GetTasksTreeStatus()
	task, err := GetTaskByID(ctx, id)
	switch err.(type) {
	case nil:
		goto Next
	case define.ErrNotExist:
		return retTasksStatus, nil
	default:
		return nil, err
	}
Next:
	// TODO 优化 只循环taskresps就可以取出
	if len(task.ParentTaskIds) != 0 {
		var isSet bool
		for _, taskid := range task.ParentTaskIds {
			var taskresp *define.TaskResp
			for _, task := range taskrepos {
				if taskid == task.TaskID && task.TaskType == define.ParentTask {
					taskresp = task
					break
				}
			}
			if taskresp == nil {
				continue
			}

			tasktreestatus := define.TaskStatusTree{
				Status:       taskresp.Status,
				ID:           taskid,
				Name:         taskresp.Task,
				TaskType:     define.ParentTask,
				TaskRespData: taskresp.LogData,
			}
			retTasksStatus[0].Children = append(retTasksStatus[0].Children, &tasktreestatus)

			if !isSet {
				// 如果存在fail那么节点的状态就是fail
				if taskresp.Status == define.TsFail.String() {
					retTasksStatus[0].Status = taskresp.Status
					isSet = true
				} else {
					retTasksStatus[0].Status = taskresp.Status
				}
			}
		}
		retTasksStatus[0].TaskType = define.ParentTask
	}

	var taskresp *define.TaskResp
	for _, task := range taskrepos {
		if id == task.TaskID && task.TaskType == define.MasterTask {
			taskresp = task
			break
		}
	}
	retTasksStatus[1].ID = taskresp.TaskID
	retTasksStatus[1].Name = taskresp.Task
	retTasksStatus[1].Status = taskresp.Status
	retTasksStatus[1].TaskRespData = taskresp.LogData
	retTasksStatus[1].TaskType = define.MasterTask

	if len(task.ChildTaskIds) != 0 {
		var isSet bool
		for _, id := range task.ChildTaskIds {
			var taskresp *define.TaskResp
			for _, task := range taskrepos {
				if id == task.TaskID && task.TaskType == define.ChildTask {
					taskresp = task
					break
				}
			}
			if taskresp == nil {
				continue
			}

			tasktreestatus := define.TaskStatusTree{
				Status:       taskresp.Status,
				ID:           id,
				Name:         taskresp.Task,
				TaskType:     define.ParentTask,
				TaskRespData: taskresp.LogData,
			}
			retTasksStatus[2].Children = append(retTasksStatus[2].Children, &tasktreestatus)

			if !isSet {
				// 如果存在fail那么节点的状态就是fail
				if taskresp.Status == define.TsFail.String() {
					retTasksStatus[2].Status = taskresp.Status
					isSet = true
				} else {
					retTasksStatus[2].Status = taskresp.Status
				}
			}
		}
		retTasksStatus[2].TaskType = define.ChildTask
	}
	return retTasksStatus, nil
}

// CleanTaskLog clean old task from time ago
func CleanTaskLog(ctx context.Context, name, taskid string, deletetime int64) (int64, error) {
	delsql := `DELETE FROM crocodile_log WHERE starttime < ?`
	args := []interface{}{deletetime}
	if name != "" {
		delsql += " AND name=? "
		args = append(args, name)
	} else if taskid != "" {
		delsql += " AND taskid=? "
		args = append(args, taskid)
	} else {
		return 0, errors.New("no delete id or name")
	}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return 0, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, delsql)
	if err != nil {
		return 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, fmt.Errorf("stmt.ExecContext failed: %w", err)
	}

	delcount, _ := res.RowsAffected()

	return delcount, nil
}

// SaveOperateLog save all user change operate
func SaveOperateLog(ctx context.Context,
	c *gin.Context, uid, username string,
	role define.Role, method, module, modulename string,
	operatetime int64, desc string, columns []define.Column) error {
	if c.GetInt("statuscode") != 0 {
		log.Error("req status code is not 0, do not save", zap.Int("statuscode", c.GetInt("statuscode")))
		return errors.New("return code is not equal 0")
	}
	log.Debug("start save operate", zap.String("username", username))
	operatesql := `INSERT INTO crocodile_operate
			(uid,
			username,
			role,
			method,
			module,
			modulename,
			operatetime,
			description,
			columns)
			VALUES
			(
				?,?,?,?,?,?,?,?,?
			)
		`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx, operatesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	defer stmt.Close()
	columnsdata, err := json.Marshal(columns)
	_, err = stmt.ExecContext(ctx, uid, username, role, method, module, modulename, operatetime, desc, columnsdata)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}

// GetOperate get operate log
func GetOperate(ctx context.Context, uid, username, method, module string, limit, offset int) ([]define.OperateLog, int, error) {
	getsql := `SELECT 
					uid,username,role,method,module,modulename, operatetime,description,columns
			   FROM 
					crocodile_operate`
	query := []string{}
	args := []interface{}{}
	var count int

	if uid != "" {
		query = append(query, " uid=? ")
		args = append(args, uid)
	}
	if username != "" {
		query = append(query, " username=? ")
		args = append(args, username)
	}
	if method != "" {
		query = append(query, " method=? ")
		args = append(args, method)
	}
	if module != "" {
		query = append(query, " module=? ")
		args = append(args, module)
	}

	if len(query) > 0 {
		getsql += "WHERE"
		getsql += strings.Join(query, "AND")
	}
	oplogs := make([]define.OperateLog, 0, limit)

	if limit > 0 {
		var err error
		count, err = countColums(ctx, getsql, args...)
		if err != nil {
			return oplogs, 0, fmt.Errorf("countColums failed: %w", err)
		}
		getsql += ` ORDER BY id DESC LIMIT ? OFFSET ?`
		args = append(args, limit, offset)
	}

	conn, err := db.GetConn(ctx)
	if err != nil {
		return oplogs, 0, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()
	log.Debug("sql", zap.String("sql", getsql))
	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return oplogs, 0, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return oplogs, 0, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}

	for rows.Next() {
		var (
			err         error
			columnsdata []byte
			oplog       define.OperateLog
			operatetime int64
		)
		// uid,username,role,method,module,modulename, operatetime,columns
		err = rows.Scan(&oplog.UID,
			&oplog.UserName,
			&oplog.Role,
			&oplog.Method,
			&oplog.Module,
			&oplog.ModuleName,
			&operatetime,
			&oplog.Desc,
			&columnsdata,
		)
		if err != nil {
			log.Error("rows.Scan failed", zap.Error(err))
			continue
		}

		oplog.OperateTime = utils.UnixToStr(operatetime)
		var columns []define.Column
		err = json.Unmarshal(columnsdata, &columns)
		if err != nil {
			log.Error("json.Unmarshal failed", zap.Error(err))
			continue
		}
		oplog.Columns = columns

		oplogs = append(oplogs, oplog)
	}
	return oplogs, count, nil
}

// SaveNewNotify save new notify
func SaveNewNotify(ctx context.Context, notify define.Notify) error {
	savesql := `INSERT INTO crocodile_notify
				(
					notyfytype,
					notifyuid,
					title,
					content,
					is_read,
					notifytime
				)
			  VALUES 
				(
					?,?,?,?,?,?
				)
					`
	log.Debug("start save new notify", zap.Any("notify", notify))
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, savesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}

	_, err = stmt.ExecContext(ctx,
		notify.NotifyType,
		notify.NotifyUID,
		notify.Title,
		notify.Content,
		false,
		notify.NotifyTime,
	)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}

	return nil
}

// GetNotifyByUID get user's notify
func GetNotifyByUID(ctx context.Context, uid string) ([]define.Notify, error) {
	getsql := `SELECT 
					id,notyfytype,title,content,notifytime
			   FROM 
					crocodile_notify 
			   WHERE 
			   		is_read=? AND notifyuid=?`
	notifys := []define.Notify{}
	conn, err := db.GetConn(ctx)
	if err != nil {
		return notifys, fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, getsql)
	if err != nil {
		return notifys, fmt.Errorf("conn.PrepareContext failed: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, false, uid)
	if err != nil {
		return notifys, fmt.Errorf("stmt.QueryContext failed: %w", err)
	}

	for rows.Next() {
		var notify define.Notify
		err = rows.Scan(&notify.ID, &notify.NotifyType, &notify.Title, &notify.Content, &notify.NotifyTime)
		if err != nil {
			log.Error("rows.Scan failed", zap.Error(err))
			continue
		}
		notify.NotifyTypeDesc = notify.NotifyType.String()
		notify.NotifyTimeDesc = utils.UnixToStr(notify.NotifyTime)
		notifys = append(notifys, notify)
	}
	return notifys, nil
}

// NotifyRead make notify status is readed
func NotifyRead(ctx context.Context, id int, uid string) error {
	updatesql := `UPDATE crocodile_notify SET is_read=? WHERE id=? AND notifyuid=?`
	conn, err := db.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("db.GetConn failed: %w", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, updatesql)
	if err != nil {
		return fmt.Errorf("conn.PrepareContext failed: %w", err)
	}
	_, err = stmt.ExecContext(ctx, true, id, uid)
	if err != nil {
		return fmt.Errorf("stmt.ExecContext failed: %w", err)
	}
	return nil
}
