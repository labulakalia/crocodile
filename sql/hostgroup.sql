DROP TABLE crocodile_hostgroup;
CREATE TABLE crocodile_hostgroup (
                                id VARCHAR(50) PRIMARY KEY NOT NULL COMMENT "ID",
                                name VARCHAR(10) NOT NULL COMMENT "名称",
                                remark VARCHAR(50) DEFAULT "" COMMENT "备注",
                                createByID VARCHAR(50) NOT NULL COMMENT "创建人ID",
                                hostIDs TEXT COMMENT  "Worker IDs",
                                createTime INT NOT NULL COMMENT "创建时间",
                                updateTime INT NOT NULL COMMENT "更新时间"
)
