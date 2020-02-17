DROP TABLE IF EXISTS  crocodile_host ;
CREATE TABLE crocodile_host (
                                id VARCHAR(50) PRIMARY KEY NOT NULL COMMENT "ID",
                                addr VARCHAR(20) UNIQUE NOT NULL COMMENT "地址",
                                hostname VARCHAR(10) NOT NULL COMMENT "主机名",
                                runningTasks TEXT DEFAULT '' COMMENT "运行的任务",
                                weight INT DEFAULT 100 COMMENT "权重",
                                stop INT DEFAULT 0 COMMENT "暂停",
                                version VARCHAR(10) NOT NULL COMMENT "版本",
                                lastUpdateTimeUnix INT NOT NULL COMMENT "更新时间",
                                remark VARCHAR(1000) DEFAULT "" COMMENT "备注"
)
