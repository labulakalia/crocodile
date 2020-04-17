-- 任务名 30
-- 子任务 父任务 最多20个
-- 报警用户数量，最多可以设置10个
-- 备注 最多100字节
CREATE TABLE IF NOT EXISTS `crocodile_task` (
	`id` CHAR(18) NOT NULL COMMENT "ID",
	`name` VARCHAR (30) NOT NULL COMMENT "任务名称",
	`taskType` INT NOT NULL DEFAULT 0 COMMENT "任务类型 1:Code 2:HTTP" ,
	`taskData` MEDIUMTEXT COMMENT "任务数据",
	`run` BOOL NOT NULL DEFAULT true COMMENT "是否自动调度运行",
	`parentTaskIds` VARCHAR (380) COMMENT "父任务ID，最多20个",
	`parentRunParallel` BOOL NOT NULL DEFAULT false COMMENT "父任务是否并行运行",
	`childTaskIds` VARCHAR (380) COMMENT "子任务ID，最多20个",
	`childRunParallel` BOOL NOT NULL  DEFAULT false  COMMENT "子任务是否并行运行",
	`createByID` CHAR (18) NOT NULL  DEFAULT "" COMMENT "创建人ID",
	`hostGroupID` CHAR (18) NOT NULL  DEFAULT "" COMMENT "主机组ID",
	`cronExpr` VARCHAR (1000) NOT NULL  DEFAULT "" COMMENT "定时任务表达式,共7位 秒、分、时、日、月、周、年",
	`timeout` INT NOT NULL DEFAULT -1 COMMENT "任务超时时间，默认-1即不设置超时时间",
	`alarmUserIds` VARCHAR (200) NOT NULL DEFAULT "" COMMENT "报警用户 最多设置10个",
	`routePolicy` INT NOT NULL DEFAULT 0 COMMENT "路由策略 1:Random 2:RoundRobin 3:Weight 4:LeastTask",
	`expectCode` INT NOT NULL  DEFAULT 0 COMMENT "期望返回码 CODE默认为0 HTTP默认为200",
	`expectContent` TEXT COMMENT "期望返回部分",
	`alarmStatus` INT NOT NULL  DEFAULT 0 COMMENT "报警策略 1:任务运行结束 2:任务运行失败 3:任务运行成功",
	`remark` VARCHAR (100) NOT NULL DEFAULT "" COMMENT "备注",
	`createTime` INT NOT NULL DEFAULT 0 COMMENT "任务创建时间 时间戳(秒)",
	`updateTime` INT NOT NULL DEFAULT 0 COMMENT "任务上次修改时间 时间戳(秒)",
	PRIMARY KEY(`id`),
	KEY `idx_name`(`name`),
	KEY `idx_cbi` (`createByID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;