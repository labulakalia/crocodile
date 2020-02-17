-- DROP TABLE crocodile_task;
CREATE TABLE crocodile_taskbakssss (
	id VARCHAR ( 50 ) PRIMARY KEY NOT NULL COMMENT "ID",
	name VARCHAR ( 10 ) NOT NULL COMMENT "名称",
	taskType INT NOT NULL DEFAULT 0 COMMENT "任务类型",
	taskData TEXT NOT NULL COMMENT "任务数据",
	run BOOL DEFAULT true COMMENT "运行",
	parentTaskIds TEXT COMMENT "父任务ID",
	parentRunParallel BOOL COMMENT "父任务并行运行",
	childTaskIds TEXT COMMENT "子任务ID",
	childRunParallel BOOL COMMENT "子任务并行运行",
	createByID VARCHAR ( 50 ) NOT NULL COMMENT "创建人ID",
	hostGroupID VARCHAR ( 50 ) NOT NULL COMMENT "主机组ID",
	cronExpr VARCHAR ( 20 ) NOT NULL COMMENT "CronExpr",
	timeout INT DEFAULT - 1 COMMENT "超时时间",
	alarmUserIds VARCHAR ( 1000 ) DEFAULT "" COMMENT "报警用户",
	routePolicy INT DEFAULT 0 COMMENT "路由策略",
	expectCode INT DEFAULT 0 COMMENT "期望返回码",
	expectContent TEXT DEFAULT "" COMMENT "期望返回内容",
	alarmStatus INT DEFAULT 0 COMMENT "报警策略",
	remark VARCHAR ( 50 ) DEFAULT "" COMMENT "备注",
	createTime INT NOT NULL COMMENT "创建时间",
updateTime INT NOT NULL COMMENT "更新时间" 
)