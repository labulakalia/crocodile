DROP TABLE  crocodile_log;
CREATE TABLE crocodile_log (
    id INTEGER  PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) NOT NULL,
    taskid VARCHAR(50) NOT NULL,
    starttime INT NOT NULL,
    endtime INT NOT NULL,
    totalruntime INT,
    status INT,
    taskresps TEXT,
    errcode INT,
    errmsg INT,
    errtasktype INT,
    errtaskid VARCHAR(50) ,
    errtask VARCHAR(50)
)