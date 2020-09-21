package model

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGenerate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		// 	// DryRun: true,
	})

	db = db.Debug()
	
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
}
