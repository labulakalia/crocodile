DROP TABLE IF EXISTS  crocodile_host ;
CREATE TABLE crocodile_host (
                                hostname VARCHAR(10) UNIQUE NOT NULL,
                                ip VARCHAR(20) NOT NULL,
                                port INT NOT NULL,
                                online int NOT NULL,
                                runingTasks TEXT,
                                version VARCHAR(10) NOT NULL,
                                lastUpdateTime INT NOT NULL
)
