DROP TABLE IF EXISTS  crocodile_operate ;
CREATE TABLE IF NOT EXISTS crocodile_operate (
        id INTEGER  PRIMARY KEY AUTOINCREMENT,
        uid VARCHAR(50) NOT NULL DEFAULT "",-- "操作用户ID",
        username VARCHAR(50) NOT NULL DEFAULT "",-- "操作用户名",
        role INTEGER NOT NULL DEFAULT 0 ,-- "操作用户类型",
        method VARCHAR(10) NOT NULL DEFAULT "" ,-- "操作类型",
        module VARCHAR(10) NOT NULL DEFAULT "" ,-- "操作模块",
        modulename VARCHAR(10) NOT NULL DEFAULT "" ,-- "操作模块名称",
        operatetime INTEGER NOT NULL DEFAULT 0 ,-- "操作时间",
        desc VARCHAR(200) NOT NULL  DEFAULT "" ,-- "描述",
        columns TEXT-- "操作字段"
)