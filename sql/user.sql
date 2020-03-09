CREATE TABLE IF NOT EXISTS crocodile_user (
                                id VARCHAR(50) PRIMARY KEY NOT NULL ,-- "ID",
                                name VARCHAR(10) NOT NULL DEFAULT "" ,-- "用户名",
                                hashpassword VARCHAR(100) NOT NULL DEFAULT "" ,-- "加密后的密码",
                                role INT(1) NOT NULL DEFAULT 0 ,-- "用户类型",
                                forbid INT(1) NOT NULL DEFAULT 0 ,-- "禁止登陆",
                                remark VARCHAR(100) NOT NULL  DEFAULT "" DEFAULT "" ,-- "备注",
                                email VARCHAR(20) NOT NULL DEFAULT "" ,-- "邮箱",
                                dingphone VARCHAR(20) NOT NULL  DEFAULT "" ,-- "DingDing",
                                telegram VARCHAR(20) NOT NULL DEFAULT "" ,-- "Telegram",
                                wechat VARCHAR(20) NOT NULL DEFAULT "" ,-- "WeChat",
                                createTime INT NOT NULL DEFAULT 0 ,-- "创建时间",
                                updateTime INT NOT NULL DEFAULT 0-- "更新时间"
);