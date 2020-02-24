DROP TABLE IF EXISTS  crocodile_host ;
CREATE TABLE IF NOT EXISTS crocodile_host (
        id VARCHAR(50) PRIMARY KEY NOT NULL ,-- "ID",
        addr VARCHAR(20) UNIQUE NOT NULL ,-- "地址",
        hostname VARCHAR(10) NOT NULL ,-- "主机名",
        runningTasks TEXT,-- "运行的任务",
        weight INT NOT NULL  DEFAULT 100 ,-- "权重",
        stop INT NOT NULL  DEFAULT 0 ,-- "暂停",
        version VARCHAR(10) NOT NULL ,-- "版本",
        lastUpdateTimeUnix INT NOT NULL DEFAULT 0,-- "更新时间",
        remark VARCHAR(1000) DEFAULT ""-- "备注"
)


