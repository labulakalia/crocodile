DROP TABLE IF EXISTS crocodile_task;
CREATE TABLE IF NOT EXISTS crocodile_task (
	id VARCHAR ( 50 ) PRIMARY KEY NOT NULL ,-- "ID",
	name VARCHAR ( 10 ) NOT NULL ,-- "名称",
	taskType INT NOT NULL DEFAULT 0 ,-- "任务类型",
	taskData TEXT,-- "任务数据",
	run BOOL NOT NULL DEFAULT true,-- "运行",
	parentTaskIds TEXT,-- "父任务ID",
	parentRunParallel BOOL NOT NULL DEFAULT false,-- "父任务并行运行",
	childTaskIds TEXT ,-- "子任务ID",
	childRunParallel BOOL NOT NULL  DEFAULT false,-- "子任务并行运行",
	createByID VARCHAR ( 50 ) NOT NULL  DEFAULT "",-- "创建人ID",
	hostGroupID VARCHAR ( 50 ) NOT NULL  DEFAULT "",-- "主机组ID",
	cronExpr VARCHAR ( 20 ) NOT NULL  DEFAULT "",-- "CronExpr",
	timeout INT NOT NULL DEFAULT -1 ,-- "超时时间",
	alarmUserIds VARCHAR (1000) NOT NULL DEFAULT "" ,-- "报警用户",
	routePolicy INT NOT NULL DEFAULT 0 ,-- "路由策略",
	expectCode INT NOT NULL  DEFAULT 0 ,-- "期望返回码",
	expectContent TEXT,-- "期望返回内容",
	alarmStatus INT NOT NULL  DEFAULT 0 ,-- "报警策略",
	remark VARCHAR ( 50 ) NOT NULL DEFAULT "" ,-- "备注",
	createTime INT NOT NULL DEFAULT 0,-- "创建时间",
	updateTime INT NOT NULL DEFAULT 0-- "更新时间" 
)

SELECT count() FROM sqlite_master WHERE type="table" 
AND name="crocodile_host" AND name=? AND name=? AND name=? AND name=? AND name=? AND name=? 
[ crocodile_hostgroup crocodile_log crocodile_notify crocodile_operate crocodile_task crocodile_user]