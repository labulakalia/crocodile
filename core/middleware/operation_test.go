package middleware

import (
	"context"
	"reflect"
	"testing"

	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"github.com/labulaka521/crocodile/core/utils/log"
)

func TestOprtation(t *testing.T) {
	// only save post, put, delete
	config.Init("/Users/labulakalia/workerspace/golang/crocodile/core/config/core.toml")
	config.CoreConf.Server.DB.Dsn = "/Users/labulakalia/workerspace/golang/crocodile/db/core.db"
	log.Init()
	model.InitDb()

	var (
		oldData interface{}
		newData interface{}
		err     error
	)

	user, err := model.GetUserByID(context.Background(), "194503907731312640")
	if err != nil {
		t.Fatal(err)
	}
	oldData = *user

	// user2, err := model.GetUserByID(context.Background(), "194503907731312640")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// user2.Telegram = "changetelegram"
	// user2.Remark = "tettete"
	// newData = *user2

	columns := make([]define.Column, 0, 50)
	parseColumn(oldData, newData, &columns,"")

	t.Logf("%+v", columns)

}

func TestO(t *testing.T) {
	type A struct {
		AA int `comment:"备注A"`
		CC int `comment:"备注C"`
	}
	type B struct {
		A  A
		BB int `comment:"备注B"`
	}

	// var c interface{}

	var b = B{
		BB: 1,
	}

	v := reflect.ValueOf(b)
	tt := reflect.TypeOf(b)
	// t.Log(v.NumField())
	t.Log(v.FieldByName("BB").Interface().(int) == 1)
	// t.Log(c == nil)
	t.Log(tt.NumField())
	t.Log(tt.Name())
	return

	for i := 0; i < v.NumField(); i++ {
		// if !tt.Field(i).Anonymous {
		// 	continue
		// }
		// t.Log(v.Field(i).Interface())
		// av := reflect.ValueOf(v.Field(i).Interface())
		// att := reflect.TypeOf(v.Field(i).Interface())

		// for i := 0; i < av.NumField(); i++ {
		// 	comment := att.Field(i).Tag.Get("comment")
		// 	comment = comment
		// 	t.Log(att.Field(i).Name)
		// }

		// t.Log(tt.Field(i).Anonymous)
		t.Log(tt.Field(i).Name)
	}
	// aa := make([]string, 0, 10)

	// a := func(p *[]string) {

	// 	*p = append(*p, "1", "111")
	// }
	// a(&aa)

	// fmt.Println(aa)
}
