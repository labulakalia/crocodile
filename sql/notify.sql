CREATE TABLE IF NOT EXISTS `crocodile_notify` (
    `id` INT AUTO_INCREMENT COMMENT "ID",
    `notyfytype` INTEGER NOT NULL DEFAULT 0 COMMENT "通知类型",
    `notifyuid` CHAR(18) NOT NULL DEFAULT "" COMMENT "通知用户ID",
    `notifytime` INT NOT NULL  DEFAULT 0 COMMENT "通知时间",
    `title` VARCHAR(30) NOT NULL DEFAULT "" COMMENT "标题",
    `content` VARCHAR(1000) NOT NULL DEFAULT "" COMMENT "内容",
    `is_read` BOOL NOT NULL DEFAULT false COMMENT "是否已读",
     PRIMARY KEY (`id`),
     KEY `idx_nu` (`notifyuid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;