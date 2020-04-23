CREATE TABLE IF NOT EXISTS `crocodile_operate` (
        `id` INT  AUTO_INCREMENT COMMENT "ID",
        `uid` CHAR(18) NOT NULL DEFAULT "" COMMENT "操作用户ID",
        `username` VARCHAR(50) NOT NULL DEFAULT "" COMMENT "操作用户名",
        `role` INT NOT NULL DEFAULT 0 COMMENT "操作用户类型",
        `method` VARCHAR(6) NOT NULL DEFAULT "" COMMENT "操作类型",
        `module` VARCHAR(10) NOT NULL DEFAULT "" COMMENT "操作模块",
        `modulename` VARCHAR(30) NOT NULL DEFAULT "" COMMENT "操作模块名称 例如任务名称",
        `operatetime` INTEGER NOT NULL DEFAULT 0 COMMENT "操作时间",
        `description` VARCHAR(200) COMMENT "操作说明，一般用户用户操作未直接改变数据库变化的操作，例如运行任务",-- "描述"
        `columns` MEDIUMTEXT COMMENT "修改的字段",
         PRIMARY KEY (`id`),
         KEY `idx_username` (`username`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;