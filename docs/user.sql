DROP TABLE crocodile_user;
CREATE TABLE crocodile_user (
                                id VARCHAR(50) PRIMARY KEY NOT NULL ,
                                name VARCHAR(10) UNIQUE NOT NULL,
                                email VARCHAR(20) UNIQUE NOT NULL,
                                hashpassword VARCHAR(10) NOT NULL,
                                role INT(1) NOT NULL DEFAULT 1,
                                forbid INT(1) NOT NULL DEFAULT 1,
                                remark VARCHAR(100),
                                dingphone VARCHAR(20),
                                slack VARCHAR(20),
                                telegram VARCHAR(20)
                                wechat VARCHAR(20),
                                createTime INT NOT NULL,
                                updateTime INT NOT NULL
);
INSERT INTO crocodile_user (id, name, email, hashpassword, role, forbid, remark, createTime, updateTime) VALUES ('194503907731312640', 'admin', 'admin@admin.com', '$2a$10$.kP56QlS3AUGN3y/tFIL0en9ivnlG0h9LtxtgR6T9OKwSVvQtrjNu', 2, 0, '', 1572078878, 1572078878);
INSERT INTO crocodile_user (id, name, email, hashpassword, role, forbid, remark, createTime, updateTime) VALUES ('194507734782054400', 'normal', 'normal@admin.com', '$2a$10$GhuMUkhRv0IHC1bxKWBjZ.9LXNu6QIlc9KwL8cJf6.3rD80roHbeW', 1, 0, '', 1572079790, 1572079790);
-- 管理用户 user: admin password: admin
-- 普通用户 user: normal password: normal