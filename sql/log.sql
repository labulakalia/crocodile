DROP TABLE  crocodile_log;
CREATE TABLE crocodile_log (
    id INTEGER  PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL DEFAULT "" COMMENT "任务名称",
    taskid VARCHAR(50) NOT NULL DEFAULT "" COMMENT "任务ID",
    starttime INT NOT NULL DEFAULT 0  COMMENT "开始时间",
    endtime INT NOT NULL DEFAULT 0  COMMENT "结束时间",
    totalruntime INT  DEFAULT 0 COMMENT "运行时间",
    status INT  DEFAULT 0  COMMENT "执行结果",
    taskresps TEXT  COMMENT "任务日志",
    trigger INT DEFAULT 0  COMMENT "触发方式",
    errcode INT DEFAULT 0  COMMENT "出错Code" ,
    errmsg INT DEFAULT ""  COMMENT "出错信息",
    errtasktype INT DEFAULT 0  COMMENT "出错任务类型",
    errtaskid VARCHAR(50) DEFAULT ""  COMMENT "出错任务ID",
    errtask VARCHAR(50) DEFAULT ""  COMMENT "出错任务名称"
)