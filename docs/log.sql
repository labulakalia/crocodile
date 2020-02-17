DROP TABLE  crocodile_log;
CREATE TABLE crocodile_log (
    id INTEGER  PRIMARY KEY AUTOINCREMENT COMMENT "ID",
    name VARCHAR(50) NOT NULL COMMENT "名称",
    taskid VARCHAR(50) NOT NULL COMMENT "任务ID",
    starttime INT NOT NULL COMMENT "开始时间",
    endtime INT NOT NULL  COMMENT "结束时间",
    totalruntime INT  COMMENT "运行时长",
    status INT  COMMENT "执行状态",
    taskresps TEXT  COMMENT "任务输出",
    trigger INT  COMMENT DEFAULT 0 "触发方式",
    errcode INT  COMMENT DEFAULT 0  "任务返回码",
    errmsg INT  COMMENT DEFAULT "" "错误任务信息",
    errtasktype INT DEFAULT 0   COMMENT "错误任务类型",
    errtaskid VARCHAR(50) DEFAULT ""  COMMENT "错误任务ID",
    errtask VARCHAR(50) DEFAULT "" COMMENT "错误任务""
)