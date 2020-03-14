CREATE TABLE IF NOT EXISTS crocodile_log (
    id INTEGER  PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL DEFAULT "" ,-- "任务名称",
    taskid VARCHAR(50) NOT NULL DEFAULT "" ,-- "任务ID",
    starttime BIGINT NOT NULL DEFAULT 0  ,-- "开始时间，毫秒",
    endtime BIGINT NOT NULL DEFAULT 0  ,-- "结束时间，毫秒",
    totalruntime INT NOT NULL  DEFAULT 0 ,-- "运行时间",
    status INT NOT NULL  DEFAULT 0  ,-- "执行结果",
    taskresps TEXT, -- "任务日志",
    triggertype INT NOT NULL  DEFAULT 0,-- "触发方式",
    errcode INT NOT NULL  DEFAULT 0,-- "出错Code" ,
    errmsg TEXT,-- "出错信息",
    errtasktype INT NOT NULL  DEFAULT 0  ,-- "出错任务类型",
    errtaskid VARCHAR(50) NOT NULL  DEFAULT ""  ,-- "出错任务ID",
    errtask VARCHAR(50) NOT NULL  DEFAULT ""-- "出错任务名称"
)