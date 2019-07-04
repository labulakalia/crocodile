CREATE DATABASE IF NOT EXISTS crocodile CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE TABLE IF NOT EXISTS crocodile_user
(
  id           VARCHAR(100) NOT NULL PRIMARY KEY ,
  username     VARCHAR(100) NOT NULL,
  hashpassword VARCHAR(200) NOT NULL,
  email        VARCHAR(100) NOT NULL,
  avatar       VARCHAR(200),
  forbid       BOOLEAN DEFAULT FALSE,
  super        BOOLEAN DEFAULT FALSE,

  UNIQUE INDEX ix_id (id ASC),
  UNIQUE INDEX ix_username (username ASC)
)