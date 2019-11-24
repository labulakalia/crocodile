DROP TABLE IF EXISTS  crocodile_host ;
CREATE TABLE crocodile_host (
                                addr VARCHAR(20) PRIMARY KEY NOT NULL,
                                hostname VARCHAR(10) NOT NULL,
                                runingTasks TEXT DEFAULT '',
                                version VARCHAR(10) NOT NULL,
                                lastUpdateTime INT NOT NULL
)
