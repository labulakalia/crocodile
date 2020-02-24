package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/labulaka521/crocodile/common/db"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// 查询表是否已经创建
//

// QueryIsInstall check table is create
func QueryIsInstall(ctx context.Context) (bool, error) {
	var querytable string
	needtables := []interface{}{}

	for tbname := range crcocodileTables {
		needtables = append(needtables, tbname)
	}
	var queryname string

	drivename := config.CoreConf.Server.DB.Drivename
	if drivename == "sqlite3" {
		querytable = `SELECT count() FROM sqlite_master WHERE type="table" AND (`
		queryname = "name"
	} else if drivename == "mysql" {
		querytable = `SELECT count() FROM information_schema.TABLES WHERE (`
		queryname = "table_name"
	} else {
		return false, fmt.Errorf("unsupport drive type %s, only support sqlite3 or mysql", drivename)
	}
	params := []string{}

	for i := 0; i < len(crcocodileTables); i++ {
		params = append(params, queryname+"=?")
	}
	querytable += strings.Join(params, " OR ")
	querytable += ")"
	var count int
	conn, err := db.GetConn(ctx)
	if err != nil {
		return false, errors.Wrap(err, "db.GetConn")
	}
	err = conn.QueryRowContext(ctx, querytable, needtables...).Scan(&count)
	if err != nil {
		log.Error("msg string", zap.Error(err))
		return false, errors.Wrap(err, "Scan")
	}

	if count != len(crcocodileTables) {
		return false, nil
	}
	return true, nil
}

var crcocodileTables = map[string]string{
	TBHost: `CREATE TABLE IF NOT EXISTS crocodile_host (
		id VARCHAR(50) PRIMARY KEY NOT NULL,-- "ID",
		addr VARCHAR(20) UNIQUE NOT NULL,-- "地址",
		hostname VARCHAR(10) NOT NULL,-- "主机名",
		runningTasks TEXT,-- "运行的任务",
		weight INT NOT NULL  DEFAULT 100,-- "权重",
		stop INT NOT NULL  DEFAULT 0,-- "暂停",
		version VARCHAR(10) NOT NULL,-- "版本",
		lastUpdateTimeUnix INT NOT NULL DEFAULT 0,-- "更新时间",
		remark VARCHAR(1000) DEFAULT ""-- "备注"
)`,
	TBHostgroup: `CREATE TABLE IF NOT EXISTS crocodile_hostgroup (
		id VARCHAR(50) PRIMARY KEY NOT NULL,-- "ID",
		name VARCHAR(10) NOT NULL DEFAULT "",-- "名称",
		remark VARCHAR(50) NOT NULL  DEFAULT "",-- "备注",
		createByID VARCHAR(50) NOT NULL DEFAULT "",-- "创建人ID",
		hostIDs TEXT,--  "Worker IDs",
		createTime INT NOT NULL DEFAULT 0,-- "创建时间",
		updateTime INT NOT NULL DEFAULT 0-- "更新时间"
)`,
	TBLog: `CREATE TABLE IF NOT EXISTS crocodile_log (
		id INTEGER  PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(50) NOT NULL DEFAULT "",-- "任务名称",
		taskid VARCHAR(50) NOT NULL DEFAULT "",-- "任务ID",
		starttime INT NOT NULL DEFAULT 0,-- "开始时间",
		endtime INT NOT NULL DEFAULT 0,-- "结束时间",
		totalruntime INT NOT NULL  DEFAULT 0,-- "运行时间",
		status INT NOT NULL  DEFAULT 0,-- "执行结果",
		taskresps TEXT,-- "任务日志",
		trigger INT NOT NULL  DEFAULT 0 ,-- "触发方式",
		errcode INT NOT NULL  DEFAULT 0 ,-- "出错Code"
		errmsg INT NOT NULL  DEFAULT "" ,-- "出错信息",
		errtasktype INT NOT NULL  DEFAULT 0 ,-- "出错任务类型",
		errtaskid VARCHAR(50) NOT NULL  DEFAULT "" ,-- "出错任务ID"
		errtask VARCHAR(50) NOT NULL  DEFAULT ""-- "出错任务名称"
	)`,
	TBNotify: `CREATE TABLE IF NOT EXISTS crocodile_notify (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		notyfytype INTEGER NOT NULL DEFAULT 0,-- "通知类型",
		notifyuid VARCHAR(50) NOT NULL DEFAULT "",-- "通知用户ID",
		notifytime INT NOT NULL  DEFAULT 0,-- "通知时间",
		title VARCHAR(30) NOT NULL DEFAULT "",-- "标题",
		content VARCHAR(1000) NOT NULL DEFAULT "",-- "内容",
		is_read BOOL NOT NULL DEFAULT false-- "已读"
	)`,
	TBOperate: `CREATE TABLE IF NOT EXISTS crocodile_operate (
        id INTEGER  PRIMARY KEY AUTOINCREMENT,
        uid VARCHAR(50) NOT NULL DEFAULT "",-- "操作用户ID",
        username VARCHAR(50) NOT NULL DEFAULT "",-- "操作用户名",
        role INTEGER NOT NULL DEFAULT 0 ,-- "操作用户类型",
        method VARCHAR(10) NOT NULL DEFAULT "" ,-- "操作类型",
        module VARCHAR(10) NOT NULL DEFAULT "" ,-- "操作模块",
        modulename VARCHAR(10) NOT NULL DEFAULT "" ,-- "操作模块名称",
        operatetime INTEGER NOT NULL DEFAULT 0 ,-- "操作时间",
        desc VARCHAR(200) NOT NULL  DEFAULT "" ,-- "描述",
        columns TEXT-- "操作字段"
)`,
	TBTask: `CREATE TABLE IF NOT EXISTS crocodile_task (
		id VARCHAR ( 50 ) PRIMARY KEY NOT NULL,-- "ID",
		name VARCHAR ( 10 ) NOT NULL,-- "名称",
		taskType INT NOT NULL DEFAULT 0,-- "任务类型",
		taskData TEXT,-- "任务数据",
		run BOOL NOT NULL DEFAULT true,-- "运行",
		parentTaskIds TEXT,-- "父任务ID",
		parentRunParallel BOOL NOT NULL DEFAULT false,-- "父任务并行运行",
		childTaskIds TEXT,-- "子任务ID",
		childRunParallel BOOL NOT NULL  DEFAULT false,-- "子任务并行运行",
		createByID VARCHAR ( 50 ) NOT NULL  DEFAULT "",-- "创建人ID",
		hostGroupID VARCHAR ( 50 ) NOT NULL  DEFAULT "",-- "主机组ID",
		cronExpr VARCHAR ( 20 ) NOT NULL  DEFAULT "",-- "CronExpr",
		timeout INT NOT NULL DEFAULT -1,-- "超时时间",
		alarmUserIds VARCHAR (1000) NOT NULL DEFAULT "",-- "报警用户",
		routePolicy INT NOT NULL DEFAULT 0,-- "路由策略",
		expectCode INT NOT NULL  DEFAULT 0,-- "期望返回码",
		expectContent TEXT,-- "期望返回内容",
		alarmStatus INT NOT NULL  DEFAULT 0,-- "报警策略",
		remark VARCHAR ( 50 ) NOT NULL DEFAULT "",-- "备注",
		createTime INT NOT NULL DEFAULT 0,-- "创建时间",
		updateTime INT NOT NULL DEFAULT 0-- "更新时间" 
	)`,
	TBUser: `CREATE TABLE IF NOT EXISTS crocodile_user (
		id VARCHAR(50) PRIMARY KEY NOT NULL,-- "ID",
		name VARCHAR(10) NOT NULL DEFAULT "",-- "用户名",
		hashpassword VARCHAR(10) NOT NULL DEFAULT "",-- "加密后的密码",
		role INT(1) NOT NULL DEFAULT 0,-- "用户类型",
		forbid INT(1) NOT NULL DEFAULT 0,-- "禁止登陆",
		remark VARCHAR(100) NOT NULL  DEFAULT "" DEFAULT "",-- "备注",
		email VARCHAR(20) NOT NULL DEFAULT "",-- "邮箱",
		dingphone VARCHAR(20) NOT NULL  DEFAULT "",-- "DingDing",
		slack VARCHAR(20) NOT NULL DEFAULT "",-- "Slack",
		telegram VARCHAR(20) NOT NULL DEFAULT "",-- "Telegram",
		wechat VARCHAR(20) NOT NULL DEFAULT "",-- "WeChat",
		createTime INT NOT NULL DEFAULT 0,-- "创建时间",
		updateTime INT NOT NULL DEFAULT 0-- "更新时间"
);`,
	TBCasbin: `BEGIN;
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/hostgroup*', '(GET)|(POST)|(DELETE)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Normal', '/api/v1/hostgroup*', '(GET)|(POST)|(DELETE)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Guest', '/api/v1/hostgroup*', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/task*', '(GET)|(POST)|(DELETE)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Normal', '/api/v1/task*', '(GET)|(POST)|(DELETE)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Guest', '/api/v1/task*', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/host*', '(GET)|(POST)|(DELETE)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Normal', '/api/v1/host*', '(GET)|(POST)|(DELETE)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Guest', '/api/v1/host*', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/user/info', '(GET)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Normal', '/api/v1/user/info', '(GET)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Guest', '/api/v1/user/info', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/user/select', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Normal', '/api/v1/user/select', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Guest', '/api/v1/user/select', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/user/registry', '(POST)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/user/all', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/user/admin', '(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/user/alarmstatus', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Normal', '/api/v1/user/alarmstatus', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Guest', '/api/v1/user/alarmstatus', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/user/operate', '(GET)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Admin', '/api/v1/notify', '(GET)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Normal', '/api/v1/notify', '(GET)|(PUT)', '', '', '');
	INSERT INTO "casbin_rule" VALUES ('p', 'Guest', '/api/v1/notify', '(GET)', '', '', '');
	COMMIT;
	`,
}

// StartInstall start install system
func StartInstall(ctx context.Context, username, password string) error {
	// create table
	conn, err := db.GetConn(ctx)
	if err != nil {
		return errors.Wrap(err, "db.GetConn")
	}

	for tbname, tbsql := range crcocodileTables {
		_, err = conn.ExecContext(ctx, tbsql)
		if err != nil {
			log.Error("conn.ExecContext failed", zap.Error(err), zap.String("tbname", tbname))
			return errors.Wrap(err, "conn.ExecContext")
		}
		// wait second
		time.Sleep(time.Second / 2)
	}
	log.Debug("Success Run All crocodile Sql")

	// create admin user
	hashpassword, err := utils.GenerateHashPass(password)
	if err != nil {
		return errors.Wrap(err, "utils.GenerateHashPass")
	}
	err = AddUser(ctx, username, hashpassword, define.AdminUser)
	if err != nil {
		return errors.Wrap(err, "AddUser")
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Error("enforcer.LoadPolicy failed", zap.Error(err))
		return errors.Wrap(err, "enforcer.LoadPolicy")
	}

	log.Debug("Success Install Crocodile")
	return nil
}
