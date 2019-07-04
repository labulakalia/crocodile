CREATE TABLE IF NOT EXISTS crocodile_task
(
  id        INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
  taskname  VARCHAR(100) NOT NULL               COMMENT '任务名称',
  command   VARCHAR(500) NOT NULL               COMMENT '命令',
  cronexpr  VARCHAR(50)  NOT NULL               COMMENT 'cron表达式',
  createdby VARCHAR(100) NOT NULL               COMMENT '创建人',
  remark    VARCHAR(500)                        COMMENT '备注',
  stop      BOOLEAN               DEFAULT FALSE COMMENT '停止调度',
  timeout   INT UNSIGNED NOT NULL DEFAULT 0     COMMENT '超时时间',
  nexttime  TIMESTAMP   NOT NULL                   COMMENT '下次运行时间',
  actuator VARCHAR(100) NOT NULL             COMMENT '任务的执行器',
  UNIQUE INDEX ix_id (id ASC),
  UNIQUE INDEX ix_taskname (taskname ASC)
) ENGINE=InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_bin
  COMMENT = '任务表';
