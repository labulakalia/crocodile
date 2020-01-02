DROP TABLE  crocodile_log;
CREATE TABLE crocodile_log (
    id INTEGER  PRIMARY KEY AUTOINCREMENT,
    taskid VARCHAR(50) NOT NULL,
    starttime INT NOT NULL,
    endtime INT NOT NULL,
    taskresps TEXT NOT NULL
)