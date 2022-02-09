package model

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/validator.v8"
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

type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

type User1 struct {
	FirstName      string `validate:"required"`
	LastName       string `validate:"required"`
	Age            uint8  `validate:"gte=0,lte=130"`
	Email          string `validate:"required,email"`
	FavouriteColor string `validate:"hexcolor|rgb|rgba"`
	// Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

func TestJsonMarshak(t *testing.T) {
	// address := &Address{
	// 	Street: "Eavesdown Docks",
	// 	Planet: "Persphone",
	// 	Phone:  "none",
	// }

	user := &User1{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000",
		// Addresses:      []*Address{address},
	}

	config := &validator.Config{TagName: "validate"}

	validate := validator.New(config)

	// returns nil or ValidationErrors ( map[string]*FieldError )
	errs := validate.Struct(user)

	if errs != nil {

		fmt.Println(errs) // output: Key: "User.Age" Error:Field validation for "Age" failed on the "lte" tag
		//	                         Key: "User.Addresses[0].City" Error:Field validation for "City" failed on the "required" tag
		merr, ok := errs.(validator.ValidationErrors)
		if !ok {
			t.Log("not ok")
			return
		}
		for k, v := range merr {
			t.Log("k", k, "v", v)
		}

		// from here you can create your own error messages in whatever language you wish

	}
}
