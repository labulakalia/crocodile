DROP TABLE IF EXISTS  crocodile_operate ;
CREATE TABLE crocodile_operate (
        id INTEGER  PRIMARY KEY AUTOINCREMENT,
        uid VARCHAR(50) NOT NULL DEFAULT "",
        username VARCHAR(50) NOT NULL DEFAULT "",
        role INTEGER NOT NULL DEFAULT 0,
        method VARCHAR(10) NOT NULL DEFAULT "",
        module VARCHAR(10) NOT NULL DEFAULT "",
        modulename VARCHAR(10) NOT NULL DEFAULT "",
        operatetime INTEGER NOT NULL DEFAULT 0,
        columns TEXT NOT NULL
)