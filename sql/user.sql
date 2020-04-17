-- 用户名 30
-- email 30
-- dingdingphone 11
-- telegram 20
-- wechat 20
CREATE TABLE IF NOT EXISTS `crocodile_user` (
    `id` CHAR(18)  NOT NULL COMMENT "用户ID",
    `name` VARCHAR(30) NOT NULL DEFAULT "" COMMENT "用户名",
    `hashpassword` VARCHAR(100) NOT NULL DEFAULT "" COMMENT "加密密码",
    `role` INT(1) NOT NULL DEFAULT 0 COMMENT "用户类型 1:普通用户 2:管理员 3:访客",
    `forbid` INT(1) NOT NULL DEFAULT 0 COMMENT "用户是否可以登陆 ",
    `remark` VARCHAR(100) NOT NULL  DEFAULT "" COMMENT "备注" ,
    `email` VARCHAR(30) NOT NULL DEFAULT "" COMMENT "邮箱",
    `dingphone` CHAR(11) NOT NULL  DEFAULT "" COMMENT "钉钉绑定的手机号",
    `telegram` VARCHAR(20) NOT NULL DEFAULT "" COMMENT "TelegramBot ID",
    `wechat` VARCHAR(20) NOT NULL DEFAULT "" COMMENT "企业微信ID",
    `createTime` INT NOT NULL DEFAULT 0 COMMENT "创建时间",
    `updateTime` INT NOT NULL DEFAULT 0 COMMENT "创建时间",
    PRIMARY KEY (`id`),
    KEY `idx_name`(`name`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;