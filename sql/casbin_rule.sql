BEGIN;
CREATE TABLE IF NOT EXISTS `casbin_rule` (
  `p_type` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/hostgroup*','(GET)|(POST)|(DELETE)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Normal','/api/v1/hostgroup*','(GET)|(POST)|(DELETE)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Guest','/api/v1/hostgroup*','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/task*','(GET)|(POST)|(DELETE)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Normal','/api/v1/task*','(GET)|(POST)|(DELETE)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Guest','/api/v1/task*','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/host*','(GET)|(POST)|(DELETE)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Normal','/api/v1/host*','(GET)|(POST)|(DELETE)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Guest','/api/v1/host*','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/user/info','(GET)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Normal','/api/v1/user/info','(GET)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Guest','/api/v1/user/info','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/user/select','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Normal','/api/v1/user/select','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Guest','/api/v1/user/select','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/user/registry','(POST)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/user/all','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/user/admin','(PUT)|(DELETE)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/user/alarmstatus','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Normal','/api/v1/user/alarmstatus','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Guest','/api/v1/user/alarmstatus','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/user/operate','(GET)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Admin','/api/v1/notify','(GET)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Normal','/api/v1/notify','(GET)|(PUT)','','','');
INSERT INTO casbin_rule (p_type,v0,v1,v2,v3,v4,v5) VALUES ('p','Guest','/api/v1/notify','(GET)','','','');
COMMIT;
