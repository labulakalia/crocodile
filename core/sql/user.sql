DROP TABLE crocodile_user;
CREATE TABLE crocodile_user (
                                id INT PRIMARY KEY NOT NULL ,
                                name VARCHAR(10) UNIQUE NOT NULL,
                                email VARCHAR(20) UNIQUE NOT NULL,
                                hashpassword VARCHAR(10) NOT NULL,
                                role INT(1) NOT NULL DEFAULT 0,
                                forbid INT(1) NOT NULL DEFAULT 1,
                                remark VARCHAR(50),
                                createTime INT NOT NULL,
                                updateTime INT NOT NULL
)