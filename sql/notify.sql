DROP TABLE IF EXISTS crocodile_notify;
CREATE TABLE crocodile_notify (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    notyfytype INTEGER NOT NULL DEFAULT 0 COMMENT "通知类型",
    notifyuid VARCHAR(50) NOT NULL DEFAULT "" COMMENT "通知用户ID",
    notifytime INT DEFAULT 0 COMMENT "通知时间",
    title VARCHAR(30) NOT NULL DEFAULT "" COMMENT "标题",
    content VARCHAR(1000) COMMENT "内容",
    is_read BOOLEN NOT NULL DEFAULT false COMMENT "已读"
)