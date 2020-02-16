DROP TABLE IF EXISTS  crocodile_host ;
CREATE TABLE crocodile_host (
                                id VARCHAR(
                                    
                                    
                                ) PRIMARY KEY NOT NULL ,
                                addr VARCHAR(20) UNIQUE NOT NULL,
                                hostname VARCHAR(10) NOT NULL,
                                runningTasks TEXT DEFAULT '',
                                weight INT DEFAULT 100,
                                stop INT DEFAULT 0,
                                version VARCHAR(10) NOT NULL,
                                lastUpdateTimeUnix INT NOT NULL,
                                remark VARCHAR(1000)
)
