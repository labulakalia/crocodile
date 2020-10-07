package model

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGenerateMysql(t *testing.T) {
	dsn := "root:crocodile@tcp(127.0.0.1:3306)/crocodile?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// logger.New()
		// 	// DryRun: true,
	})

	db = db.Debug()
	// db.Table(name string, args ...interface{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(
		&Host{},
		&HostGroup{},
		&Log{},
		&Notify{},
		&Operate{},
		&Task{},
		&User{},
	)
	
	// CREATE TABLE `test_crocodile_host` (`id` CHAR(18),`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`created` datetime(3) NULL,`addr` varchar(25) NOT NULL,`host_name` varchar(100) NOT NULL,`count_run_tasks` bigint NOT NULL,`weight` bigint NOT NULL DEFAULT 100,`stop` boolean NOT NULL DEFAULT false,`version` varchar(10) NOT NULL,`remark` varchar(100) NOT NULL DEFAULT "",PRIMARY KEY (`id`),INDEX idx_test_crocodile_host_deleted_at (`deleted_at`,`deleted_at`),INDEX idx_test_crocodile_host_id (`id`),UNIQUE INDEX idx_test_crocodile_host_addr (`addr`))
}
