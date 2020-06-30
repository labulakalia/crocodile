CREATE TABLE IF NOT EXISTS `crocodile_log` (
    `id` INT AUTO_INCREMENT COMMENT "ID",
    `name` VARCHAR(30) NOT NULL DEFAULT "" COMMENT "任务名称",
    `taskid` CHAR(18) NOT NULL DEFAULT "" COMMENT "任务ID",
    `starttime` BIGINT NOT NULL DEFAULT 0  COMMENT "开始时间毫秒",
    `endtime` BIGINT NOT NULL DEFAULT 0 COMMENT "结束时间 毫秒",
    `totalruntime` INT NOT NULL  DEFAULT 0 COMMENT "总共运行时间",
    `status` INT NOT NULL  DEFAULT 0  COMMENT "执行结束 1:成功 -1:失败",
    `taskresps` MEDIUMTEXT COMMENT "任务日志",
    `triggertype` INT NOT NULL  DEFAULT 0 COMMENT "触发方式",
    `errcode` INT NOT NULL  DEFAULT 0 COMMENT "错误返回码",
    `errmsg` TEXT COMMENT "错误信息",-- "出错信息",
    `errtasktype` INT NOT NULL  DEFAULT 0 COMMENT "出错任务类型",
    `errtaskid` CHAR(18) NOT NULL  DEFAULT ""  COMMENT "出错任务ID",
    `errtask` CHAR(30) NOT NULL  DEFAULT "" COMMENT "出错任务名称",
     PRIMARY KEY (`id`),
     KEY `idx_name` (`name`),
     KEY `idx_s_t` (`starttime`,`taskid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
