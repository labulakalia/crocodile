CREATE TABLE IF NOT EXISTS `crocodile_hostgroup` (
    `id` CHAR(18) NOT NULL COMMENT "ID",
    `name` VARCHAR(30) NOT NULL DEFAULT "" COMMENT "主机组名称",
    `remark` VARCHAR(100) NOT NULL  DEFAULT "" COMMENT "备注",
    `createByID` CHAR(18) NOT NULL DEFAULT "" COMMENT "创建人ID",
    `hostIDs` TEXT COMMENT "主机ID",
    `createTime` INT NOT NULL DEFAULT 0 COMMENT "创建时间",
    `updateTime` INT NOT NULL DEFAULT 0 COMMENT "更新时间",
    PRIMARY KEY (`id`),
    KEY `idx_name` (`name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
