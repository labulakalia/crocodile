DROP TABLE IF EXISTS crocodile_notify;
CREATE TABLE IF NOT EXISTS crocodile_notify (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    notyfytype INTEGER NOT NULL DEFAULT 0 ,-- "通知类型",
    notifyuid VARCHAR(50) NOT NULL DEFAULT "" ,-- "通知用户ID",
    notifytime INT NOT NULL  DEFAULT 0 ,-- "通知时间",
    title VARCHAR(30) NOT NULL DEFAULT "" ,-- "标题",
    content VARCHAR(1000) NOT NULL DEFAULT "" ,-- "内容",
    is_read BOOL NOT NULL DEFAULT false-- "已读"
)