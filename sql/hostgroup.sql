CREATE TABLE IF NOT EXISTS crocodile_hostgroup (
                                id VARCHAR(50) PRIMARY KEY NOT NULL,-- "ID",
                                name VARCHAR(50) NOT NULL DEFAULT "",-- "名称",
                                remark VARCHAR(50) NOT NULL  DEFAULT "" ,-- "备注",
                                createByID VARCHAR(50) NOT NULL DEFAULT "",-- "创建人ID",
                                hostIDs TEXT ,--  "Worker IDs",
                                createTime INT NOT NULL DEFAULT 0,-- "创建时间",
                                updateTime INT NOT NULL DEFAULT 0-- "更新时间"
)
