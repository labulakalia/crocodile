DROP TABLE IF EXISTS  crocodile_host ;
CREATE TABLE crocodile_host (
                                id VARCHAR(50) PRIMARY KEY NOT NULL ,
                                addr VARCHAR(20) UNIQUE NOT NULL,
                                hostname VARCHAR(10) NOT NULL,
                                runingTasks TEXT DEFAULT '',
                                version VARCHAR(10) NOT NULL,
                                lastUpdateTimeUnix INT NOT NULL
)
