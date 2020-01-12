DROP TABLE  crocodile_log;
CREATE TABLE crocodile_log (
    id INTEGER  PRIMARY KEY AUTOINCREMENT,
    taskid VARCHAR(50) NOT NULL,
    starttime INT NOT NULL,
    endtime INT NOT NULL,
    tasklog TEXT NOT NULL,
    splitpoint VARCHAR(300) NOT NULL
    -- 日志分割的间隔的点 会将每个任务的日志大小记录下来
)