CREATE TABLE crocodile_log
(
  id    INTEGER PRIMARY KEY AUTO_INCREMENT,
  taskname  VARCHAR(100)    NOT NULL               COMMENT '任务名称',
  command   VARCHAR(500)    NOT NULL               COMMENT '命令',
  cronexpr  VARCHAR(50)     NOT NULL               COMMENT 'cron表达式',
  createdby VARCHAR(100)    NOT NULL               COMMENT '创建人',
  timeout   INT UNSIGNED    NOT NULL               COMMENT '超时时间',
  actuator  VARCHAR(100)  NOT NULL               COMMENT '任务的执行器',
  runhost VARCHAR(100)      NOT NULL             COMMENT '执行主机地址',
  starttime TIMESTAMP       NOT NULL                 COMMENT '任务开始时间',
  endtime TIMESTAMP         NOT NULL                 COMMENT '任务结束时间',
  output TEXT       NOT NULL                 COMMENT '命令输出结果',
  err TEXT                        COMMENT '执行失败的输出'
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  COMMENT = '日志';