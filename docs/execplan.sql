DROP TABLE crocodile_execplan;
CREATE TABLE crocodile_execplan (
                                    id VARCHAR(50) PRIMARY KEY NOT NULL ,
                                     name VARCHAR(10) UNIQUE NOT NULL,
                                     cronExpr VARCHAR(20) NOT NULL,
                                     timeout INT DEFAULT 0,
                                     runTime INT DEFAULT 0,
                                     alarmTotal INT DEFAULT 0,
                                     alarmUser VARCHAR(200),
                                     autoSwitch INT DEFAULT 0,
                                    createByID VARCHAR(50) NOT NULL,
                                    hostGroupID VARCHAR(50) NOT NULL,
                                     remark VARCHAR(50),
                                     createTime INT NOT NULL,
                                     updateTime INT NOT NULL
)
-- 执行计划 可以将一个任务绑定在这个平台
