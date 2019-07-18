CREATE DATABASE IF NOT EXISTS crocodile CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS crocodile_actuator
                   (
                     id         INT AUTO_INCREMENT,
                     name       VARCHAR(100)  NOT NULL    COMMENT '执行器名称',
                     address    VARCHAR(500)              COMMENT '执行器地址列表',
                     createdby  VARCHAR(50)   NOT NULL    COMMENT '创建人',
                     PRIMARY KEY (id),
                     UNIQUE INDEX i_n (name)
                   ) ENGINE=InnoDB
                     DEFAULT CHARSET = utf8mb4
                     COLLATE = utf8mb4_bin
                     COMMENT = '执行器表';
INSERT INTO crocodile_actuator  (`name`,`address`,`createdby`) VALUE ('actuator','','labulaka');

CREATE TABLE IF NOT EXISTS crocodile_task
(
  id        INT AUTO_INCREMENT,
  taskname  VARCHAR(100) NOT NULL               COMMENT '任务名称',
  command   VARCHAR(500) NOT NULL               COMMENT '命令',
  cronexpr  VARCHAR(50)  NOT NULL               COMMENT 'cron表达式',
  createdby VARCHAR(100) NOT NULL               COMMENT '创建人',
  remark    VARCHAR(500)                        COMMENT '备注',
  stop      BOOLEAN               DEFAULT FALSE COMMENT '停止调度',
  timeout   INT UNSIGNED NOT NULL DEFAULT 0     COMMENT '超时时间',
  nexttime  TIMESTAMP   NOT NULL                COMMENT '下次运行时间',
  actuator VARCHAR(100) NOT NULL                COMMENT '任务的执行器',
  PRIMARY KEY (id),
  UNIQUE INDEX i_t (taskname) # 查询任务名称
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  COMMENT = '任务表';
INSERT INTO crocodile_task (taskname, command, cronexpr, createdby, remark,nexttime, actuator)
VALUES ('task1','echo test1_task','*/10 * * * * * *','labulaka','',NOW(),'actuator'),
       ('task2','echo test1_task','*/11 * * * * * *','labulaka','',NOW(),'actuator'),
       ('task3','echo test1_task','*/12 * * * * * *','labulaka','',NOW(),'actuator'),
       ('task4','echo test1_task','*/13 * * * * * *','labulaka','',NOW(),'actuator'),
       ('task5','echo test1_task','*/14 * * * * * *','labulaka','',NOW(),'actuator');

CREATE TABLE IF NOT EXISTS crocodile_log
(
  id        INTEGER         AUTO_INCREMENT,
  taskname  VARCHAR(100)    NOT NULL                COMMENT '任务名称',
  command   VARCHAR(500)    NOT NULL                COMMENT '命令',
  cronexpr  VARCHAR(50)     NOT NULL                COMMENT 'cron表达式',
  createdby VARCHAR(100)    NOT NULL                COMMENT '创建人',
  timeout   INT UNSIGNED    NOT NULL                COMMENT '超时时间',
  actuator  VARCHAR(100)    NOT NULL                COMMENT '任务的执行器',
  runhost   VARCHAR(100)    NOT NULL                COMMENT '执行主机地址',
  starttime TIMESTAMP       NOT NULL                COMMENT '任务开始时间',
  endtime   TIMESTAMP       NOT NULL                COMMENT '任务结束时间',
  output    VARCHAR(500)    NOT NULL DEFAULT ''     COMMENT '命令输出结果',
  err       VARCHAR(500)    NOT NULL DEFAULT ''     COMMENT '执行失败的输出',
  PRIMARY KEY (id),
  INDEX i_t_s (taskname,starttime)
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  COMMENT = '任务日志';

CREATE TABLE IF NOT EXISTS crocodile_user
(
  id           INT AUTO_INCREMENT,
  username     VARCHAR(100) NOT NULL COMMENT '用户名',
  hashpassword VARCHAR(200) NOT NULL COMMENT '哈希密码',
  email        VARCHAR(100) NOT NULL COMMENT '用户邮箱',
  avatar       VARCHAR(200) NOT NULL COMMENT '用户图像',
  forbid       BOOLEAN DEFAULT FALSE COMMENT '禁止登录',
  super        BOOLEAN DEFAULT FALSE COMMENT '用户类型',
  PRIMARY KEY (id),
  UNIQUE INDEX i_u (username)
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  COMMENT = '用户表';

INSERT INTO crocodile_user (username, hashpassword, email, avatar,forbid,super)
VALUE ('labulaka',
       '$2a$10$c/oWAYCrUFz.LDYMXINi2ei.NN9.Q81Afd5LQKnAypKOA0YdE1jE2',
       'labulaka@qq.com','https://www.gravatar.com/avatar/d41d8cd98f00b204e9800998ecf8427e?d=identicon&s=128',
       FALSE,TRUE);