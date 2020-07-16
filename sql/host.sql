CREATE TABLE IF NOT EXISTS `crocodile_host` (
        `id` CHAR(18) NOT NULL COMMENT "主机ID",
        `addr` VARCHAR(25) NOT NULL COMMENT "Host地址",
        `hostname` VARCHAR(100) NOT NULL COMMENT "主机名",
        `runningTasks` TEXT COMMENT "运行的任务",
        `weight` INT NOT NULL DEFAULT 100 COMMENT "权重",
        `stop` INT NOT NULL  DEFAULT 0 COMMENT "主机暂停执行任务",
        `version` VARCHAR(10) NOT NULL COMMENT "版本号",
        `lastUpdateTimeUnix` INT NOT NULL DEFAULT 0 COMMENT "更新时间",
        `remark` VARCHAR(100) DEFAULT "" COMMENT "备注",
        PRIMARY KEY(`id`),
        UNIQUE KEY `idx_addr` (`addr`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


