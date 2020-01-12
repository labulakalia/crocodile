-- host 与 hostgroup 对应关系
CREATE TABLE crocodile_host_hostgroup (
        id INTEGER  PRIMARY KEY AUTOINCREMENT,
        hostid VARCHAR(50) NOT NULL ,
        hostgroupid VARCHAR(50) NOT NULL 
)