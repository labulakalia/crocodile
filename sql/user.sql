DROP TABLE IF EXISTS crocodile_user;
CREATE TABLE IF NOT EXISTS crocodile_user (
                                id VARCHAR(50) PRIMARY KEY NOT NULL ,-- "ID",
                                name VARCHAR(10) NOT NULL DEFAULT "" ,-- "用户名",
                                hashpassword VARCHAR(10) NOT NULL DEFAULT "" ,-- "加密后的密码",
                                role INT(1) NOT NULL DEFAULT 0 ,-- "用户类型",
                                forbid INT(1) NOT NULL DEFAULT 0 ,-- "禁止登陆",
                                remark VARCHAR(100) NOT NULL  DEFAULT "" DEFAULT "" ,-- "备注",
                                email VARCHAR(20) NOT NULL DEFAULT "" ,-- "邮箱",
                                dingphone VARCHAR(20) NOT NULL  DEFAULT "" ,-- "DingDing",
                                slack VARCHAR(20) NOT NULL DEFAULT "" ,-- "Slack",
                                telegram VARCHAR(20) NOT NULL DEFAULT "" ,-- "Telegram",
                                wechat VARCHAR(20) NOT NULL DEFAULT "" ,-- "WeChat",
                                createTime INT NOT NULL DEFAULT 0 ,-- "创建时间",
                                updateTime INT NOT NULL DEFAULT 0-- "更新时间"
);
INSERT INTO crocodile_user (id, name, email, hashpassword, role, forbid, remark, createTime, updateTime) VALUES ('194503907731312640', 'admin', 'admin@admin.com', '$2a$10$.kP56QlS3AUGN3y/tFIL0en9ivnlG0h9LtxtgR6T9OKwSVvQtrjNu', 2, 0, '', 1572078878, 1572078878);
INSERT INTO crocodile_user (id, name, email, hashpassword, role, forbid, remark, createTime, updateTime) VALUES ('194507734782054400', 'normal', 'normal@admin.com', '$2a$10$GhuMUkhRv0IHC1bxKWBjZ.9LXNu6QIlc9KwL8cJf6.3rD80roHbeW', 1, 0, '', 1572079790, 1572079790);
-- 管理用户 user: admin password: admin
-- 普通用户 user: normal password: normal