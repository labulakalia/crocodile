package tasklog

import (
	"context"
	pbtasklog "crocodile/service/tasklog/proto/tasklog"
	"database/sql"
	"github.com/golang/protobuf/ptypes"
	"github.com/labulaka521/logging"
	"time"
)

type Servicer interface {
	CreateLog(ctx context.Context, simplelog *pbtasklog.SimpleLog) (err error)
	GetLog(ctx context.Context, querylog *pbtasklog.QueryLog) (resplog *pbtasklog.RespLog, err error)
}

type Service struct {
	DB *sql.DB
}

var _ Servicer = &Service{}

// 	新建日志
func (s *Service) CreateLog(ctx context.Context, log *pbtasklog.SimpleLog) (err error) {
	var (
		createlog_sql string
		stmt          *sql.Stmt
		starttime     time.Time
		endtime       time.Time
	)

	starttime, _ = ptypes.Timestamp(log.Starttime)
	endtime, _ = ptypes.Timestamp(log.Endtime)
	createlog_sql = `INSERT INTO crocodile_log 
					(taskname,command,cronexpr,createdby,timeout,actuator,runhost,starttime,endtime,output,err)
					VALUE(?,?,?,?,?,?,?,?,?,?,?)`
	if stmt, err = s.DB.PrepareContext(ctx, createlog_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", createlog_sql, err)
		return
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, log.Taskname, log.Command, log.Cronexpr, log.Createdby,
		log.Timeout, log.Actuator, log.Runhost, starttime, endtime, log.Output, log.Err)
	if err != nil {
		logging.Errorf("SQL %s Exec Err: %v", createlog_sql, err)
	}
	return
}

// 查询日志
func (s *Service) GetLog(ctx context.Context, querylog *pbtasklog.QueryLog) (resplog *pbtasklog.RespLog, err error) {
	var (
		getlog_sql string
		count_sql  string
		stmt       *sql.Stmt
		rows       *sql.Rows
		starttime  time.Time
		endtime    time.Time
		fromtime   time.Time
		totime     time.Time
	)
	resplog = &pbtasklog.RespLog{}
	getlog_sql = `SELECT *
				 FROM crocodile_log 
				 WHERE taskname=?
				 AND starttime BETWEEN ? AND ?
				 ORDER BY starttime DESC
				 LIMIT ? OFFSET ?
				`
	count_sql = `SELECT COUNT(id) 
				 FROM crocodile_log 
				 WHERE taskname=?
				 AND starttime BETWEEN ? AND ?`

	fromtime, _ = ptypes.Timestamp(querylog.Fromtime)
	totime, _ = ptypes.Timestamp(querylog.Totime)
	if stmt, err = s.DB.PrepareContext(ctx, getlog_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", getlog_sql, err)
		return
	}
	defer stmt.Close()

	if rows, err = stmt.QueryContext(ctx, querylog.Taskname, fromtime.Format(time.RFC3339), totime.Format(time.RFC3339), querylog.Limit, querylog.Offset); err != nil {
		logging.Errorf("SQL %s Exec Err: %v", getlog_sql, err)
		return
	}

	for rows.Next() {
		simplelog := pbtasklog.SimpleLog{}
		err = rows.Scan(&simplelog.Id, &simplelog.Taskname, &simplelog.Command, &simplelog.Cronexpr, &simplelog.Createdby,
			&simplelog.Timeout, &simplelog.Actuator, &simplelog.Runhost, &starttime, &endtime, &simplelog.Output, &simplelog.Err)
		if err != nil {
			continue
		}
		simplelog.Starttime, _ = ptypes.TimestampProto(starttime)
		simplelog.Endtime, _ = ptypes.TimestampProto(endtime)
		resplog.Logs = append(resplog.Logs, &simplelog)
	}

	if stmt, err = s.DB.PrepareContext(ctx, count_sql); err != nil {
		logging.Errorf("Prepare SQL %s Err: %v", getlog_sql, err)
		return
	}
	defer stmt.Close()
	if err = stmt.QueryRowContext(ctx, querylog.Taskname, fromtime.Format(time.RFC3339), totime.Format(time.RFC3339)).Scan(&resplog.Count); err != nil {
		logging.Errorf("SQL %s Exec Err: %v", getlog_sql, err)
		return
	}

	return
}
