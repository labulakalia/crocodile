CREATE TABLE IF NOT EXISTS crocodile_actuator
                   (
                     id    INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
                     name  VARCHAR(100) NOT NULL   COMMENT '执行器名称',
                     address TEXT COMMENT '执行器地址列表',
                     createdby VARCHAR(50) NOT NULL COMMENT '创建人',
                     UNIQUE INDEX ix_id (id ASC),
                     UNIQUE INDEX ix_name (name ASC)
                   ) ENGINE=InnoDB
                     DEFAULT CHARSET = utf8mb4
                     COLLATE = utf8mb4_bin
  COMMENT = '执行器表';