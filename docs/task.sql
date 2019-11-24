DROP TABLE crocodile_task;
CREATE TABLE crocodile_task (
        id VARCHAR(50) PRIMARY KEY NOT NULL ,
        name VARCHAR(10) UNIQUE NOT NULL,
        run INT DEFAULT 1,
        parentTaskIds TEXT,
        parentRunParallel INT DEFAULT 0,
        childTaskIds TEXT,
        childRunParallel INT DEFAULT 0,
        taskType INT NOT NULL,      -- request API,run shell
        taskData TEXT NOT NULL, -- 任务的数据等信息 json序列化后存储
            -- 任务的类型
            -- program 程序 必须在主机上已经存在 return code  repp content
            -- reqapi 请求一个http接口 return status code  resp context
        createByID VARCHAR(50) NOT NULL,
        hostGroupID VARCHAR(50) NOT NULL,
        cronExpr VARCHAR(20) NOT NULL,
        timeout INT DEFAULT 0,
        alarmUserIds VARCHAR(1000),
        autoSwitch INT DEFAULT 0,
        remark VARCHAR(50),
        createTime INT NOT NULL,
        updateTime INT NOT NULL
)