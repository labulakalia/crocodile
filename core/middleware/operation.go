package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"go.uber.org/zap"
)

const (
	task      = "task"
	hostgroup = "hostgroup"
	host      = "host"
	user      = "user"
)

var moduleMap = map[string]string{
	task:      "任务",
	hostgroup: "主机组",
	host:      "主机",
	user:      "用户",
}

// Oprtation save all user operate log
func Oprtation() func(c *gin.Context) {
	return func(c *gin.Context) {

		for _, url := range excludepath {
			if strings.Contains(c.Request.RequestURI, url) {
				c.Next()
				return
			}
		}

		// save operate on request method [POST,PUT,DELETE]
		if !(c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodDelete) {
			c.Next()
			return
		}

		var (
			module       string
			oldData      interface{}
			newData      interface{}
			name         string
			id           string
			exist        bool
			modulename   string
			operatetimne int64
		)

		operatetimne = time.Now().Unix()

		if len(strings.Split(c.Request.URL.Path, "/")) > 3 {
			module = strings.Split(c.Request.URL.Path, "/")[3]
		}
		if _, exist = moduleMap[module]; !exist {
			c.Next()
			return
		}

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error("ioutil.ReadAll failed", zap.Error(err))
			c.Next()
			return
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		ctx, cancel := context.WithTimeout(context.Background(),
			config.CoreConf.Server.DB.MaxQueryTime.Duration)
		defer cancel()

		uid := c.GetString("uid")
		username := c.GetString("username")
		// 获取用户的类型
		var role define.Role
		if v, ok := c.Get("role"); ok {
			role = v.(define.Role)
		}

		columns := make([]define.Column, 0, 50)

		// 一些不能使用修改前后两次数据对比的审计
		// 复制任务
		// /api/v1/task/clone post
		if c.Request.Method == http.MethodPost && c.Request.URL.Path == "/api/v1/task/clone" {
			idname := define.IDName{}
			err = json.Unmarshal(body, &idname)
			if err != nil {
				log.Error("json.Unmarshal failed", zap.Error(err))
				c.Next()
				return
			}
			clonetask, err := model.GetTaskByID(ctx, idname.ID)
			if err != nil {
				log.Error("model.GetTaskByID failed", zap.Error(err))
				c.Next()
				return
			}
			c.Next() // 为了获取状态码
			// 任务ID不存在时，返回值为nil
			if c.GetInt("statuscode") != 0 {
				log.Error("req status code is not 0, do not save", zap.Int("statuscode", c.GetInt("statuscode")))
				return
			}
			err = model.SaveOperateLog(ctx, c,
				uid,
				username,
				role,
				c.Request.Method,
				module,
				idname.Name,
				operatetimne,
				fmt.Sprintf("通过任务 %s 克隆新的任务 %s", clonetask.Name, idname.Name), columns)
			if err != nil {
				log.Error("model.SaveOperateLog failed", zap.Error(err))
			}
			return
		}
		// 删除日志
		// /api/v1/task/log delete
		if c.Request.Method == http.MethodDelete && c.Request.URL.Path == "/api/v1/task/log" {
			cleanlog := define.Cleanlog{}
			err = json.Unmarshal(body, &cleanlog)
			if err != nil {
				log.Error("json.Unmarshal failed", zap.Error(err))
				c.Next()
				return
			}
			var desc string
			if cleanlog.PreDay > 0 {
				desc = fmt.Sprintf("清除任务%s %d天前的任务日志", cleanlog.Name, cleanlog.PreDay)
			} else {
				desc = "清除全部日志"
			}
			c.Next() // 为了获取状态码
			model.SaveOperateLog(ctx, c,
				uid,
				username,
				role,
				c.Request.Method,
				module,
				cleanlog.Name,
				operatetimne,
				desc, columns)
			if err != nil {
				log.Error("model.SaveOperateLog failed", zap.Error(err))
			}
			return
		}
		// 运行任务
		// /api/v1/task/run PUT
		if c.Request.Method == http.MethodPut && c.Request.URL.Path == "/api/v1/task/run" {
			getid := define.GetID{}
			err = json.Unmarshal(body, &getid)
			if err != nil {
				log.Error("json.Unmarshal failed", zap.Error(err))
				c.Next()
				return
			}
			task, err := model.GetTaskByID(ctx, getid.ID)
			if err != nil {
				log.Error("model.GetTaskByID failed", zap.Error(err))
				c.Next()
				return
			}
			c.Next()
			model.SaveOperateLog(ctx, c,
				uid,
				username,
				role,
				c.Request.Method,
				module,
				task.Name,
				operatetimne,
				fmt.Sprintf("触发运行任务%s", task.Name), columns)
			if err != nil {
				log.Error("model.SaveOperateLog failed", zap.Error(err))
			}
			return

		}
		// 杀死任务
		// /api/v1/task/kill PUT
		if c.Request.Method == http.MethodPut && c.Request.URL.Path == "/api/v1/task/kill" {
			getid := define.GetID{}
			err = json.Unmarshal(body, &getid)
			if err != nil {
				log.Error("json.Unmarshal failed", zap.Error(err))
				c.Next()
				return
			}
			task, err := model.GetTaskByID(ctx, getid.ID)
			if err != nil {
				log.Error("model.GetTaskByID failed", zap.Error(err))
				c.Next()
				return
			}
			c.Next() // 为了获取状态码
			model.SaveOperateLog(ctx, c,
				uid,
				username,
				role,
				c.Request.Method,
				module,
				task.Name,
				operatetimne,
				fmt.Sprintf("终止任务%s", task.Name), columns)
			if err != nil {
				log.Error("model.SaveOperateLog failed", zap.Error(err))
			}
			return
		}

		// get old data
		switch c.Request.Method {
		case http.MethodPost:
			getname := define.GetName{}
			err := json.Unmarshal(body, &getname)
			if err != nil {
				log.Error("json.Unmarshal failed", zap.Error(err))
				c.Next()
				return
			}
			name = getname.Name

		case http.MethodPut, http.MethodDelete:
			getid := define.GetID{}
			err := json.Unmarshal(body, &getid)
			if err != nil {
				log.Error("json.Unmarshal failed", zap.Error(err))
				c.Next()
				return
			}
			id = getid.ID
			switch module {
			case user:
				userData, err := model.GetUserByID(ctx, id)
				if err != nil {
					log.Error("model.GetUserByID", zap.String("uid", id), zap.Error(err))
					c.Next()
					return
				}
				modulename = userData.Name
				oldData = *userData
			case hostgroup:
				hostgroupData, err := model.GetHostGroupByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostGroupByID", zap.String("id", id), zap.Error(err))
					c.Next()
					return
				}
				modulename = hostgroupData.Name
				oldData = *hostgroupData
			case host:
				hostData, err := model.GetHostByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					c.Next()
					return
				}
				modulename = hostData.Addr
				oldData = *hostData
			case task:
				taskData, err := model.GetTaskByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					c.Next()
					return
				}
				modulename = taskData.Name
				oldData = *taskData
			default:
				log.Debug("can get module name from url", zap.String("url", c.Request.URL.Path))
				c.Next()
				return
			}
		default:
			c.Next()
			return
		}

		c.Next()

		if c.GetInt("statuscode") != 0 {
			log.Error("req status code is not 0", zap.Int("statuscode", c.GetInt("statuscode")))
			return
		}

		// get new data
		switch c.Request.Method {
		case http.MethodPost:
			switch module {
			case user:
				userData, err := model.GetUserByName(ctx, name)
				if err != nil {
					log.Error("model.GetUserByID", zap.String("name", name), zap.Error(err))
					c.Next()
					return
				}
				modulename = userData.Name
				newData = *userData
			case hostgroup:
				hostgroupData, err := model.GetHostGroupByName(ctx, name)
				if err != nil {
					log.Error("model.GetHostGroupByName", zap.String("name", name), zap.Error(err))
					c.Next()
					return
				}
				modulename = hostgroupData.Name
				newData = *hostgroupData
			case task:
				taskData, err := model.GetTaskByName(ctx, name)
				if err != nil  {
					log.Error("model.GetTaskByName", zap.String("name", name), zap.Error(err))
					c.Next()
					return
				}
				modulename = taskData.Name
				newData = *taskData
			default:
				log.Debug("can get module name from url", zap.String("url", c.Request.URL.Path))
				c.Next()
				return
			}
		case http.MethodPut:
			switch module {
			case user:
				userData, err := model.GetUserByID(ctx, id)
				if err != nil {
					log.Error("model.GetUserByID", zap.String("uid", id), zap.Error(err))
					c.Next()
					return
				}
				newData = *userData
			case hostgroup:
				hostgroupData, err := model.GetHostGroupByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostGroupByID", zap.String("id", id), zap.Error(err))
					c.Next()
					return
				}
				newData = *hostgroupData
			case host:
				hostData, err := model.GetHostByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					c.Next()
					return
				}
				newData = *hostData
			case task:
				taskData, err := model.GetTaskByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					c.Next()
					return
				}
				newData = *taskData
			default:
				log.Debug("can get module name from url", zap.String("url", c.Request.URL.Path))
				c.Next()
				return
			}
		}

		parseColumn(oldData, newData, &columns, "")

		if len(columns) == 0 {
			return
		}

		err = model.SaveOperateLog(ctx, c,
			uid,
			username,
			role,
			c.Request.Method,
			module,
			modulename,
			operatetimne,
			"",
			columns)
		if err != nil {
			log.Error("model.SaveOperateLog failed", zap.Error(err))
		}
	}
}

func parseColumn(oldData, newData interface{}, columns *[]define.Column, precm string) {
	oldt := reflect.TypeOf(oldData)

	newt := reflect.TypeOf(newData)

	oldv := reflect.ValueOf(oldData)

	newv := reflect.ValueOf(newData)

	var total int

	// 特殊情况 结构体不相同 任务中TaskData
	if oldv.Kind() != reflect.Invalid &&
		newv.Kind() != reflect.Invalid &&
		oldt.Name() != newt.Name() {
		parsestruct(oldData, true, columns, precm)
		parsestruct(newData, false, columns, precm)
		return
	}

	// 默认为结构体相同

	if oldv.Kind() != reflect.Invalid {
		total = oldv.NumField() // 字段数量
	} else if newv.Kind() != reflect.Invalid {
		total = newv.NumField()
	}

	for i := 0; i < total; i++ {
		var (
			comment    string
			oldvalue   interface{}
			newvalue   interface{}
			anonStruct bool
		)

		if oldv.Kind() != reflect.Invalid {
			ft := oldt.Field(i)
			if ft.Anonymous {
				anonStruct = true
			} else {
				cm := ft.Tag.Get("comment") // get tag comment
				if cm == "" {
					continue
				}
				comment = cm
			}
			oldvalue = oldv.Field(i).Interface() // get value
		}

		if newv.Kind() != reflect.Invalid {
			ft := newt.Field(i)
			if ft.Anonymous {
				anonStruct = true
			} else {
				cm := ft.Tag.Get("comment")
				if cm == "" {
					continue
				}
				comment = cm
			}
			newvalue = newv.Field(i).Interface()
		}
		// 当字段是匿名结构体
		if anonStruct {
			parseColumn(oldvalue, newvalue, columns, "")
			continue
		}
		if reflect.DeepEqual(oldvalue, newvalue) { // check is equal
			continue
		}

		// 任务数据是结构体数据类型
		if comment == "任务数据" {
			parseColumn(oldvalue, newvalue, columns, comment)
			continue
		}
		// 取消密码回显
		if comment == "密码" {
			oldvalue = "-"
			newvalue = "-"
		}
		if precm != "" {
			var name string
			if oldv.Kind() != reflect.Invalid {
				name = oldt.Name()
			} else if newv.Kind() != reflect.Invalid {
				name = newt.Name()
			}
			comment = fmt.Sprintf("%s-%s-%s", precm, name, comment)
		}
		c := define.Column{
			Name:     comment,
			OldValue: oldvalue,
			NewValue: newvalue,
		}
		*columns = append(*columns, c)
	}
}

// parsestruct get struct tag comment and struct's value
func parsestruct(structdata interface{},
	isold bool, /*judge value is old or new */
	columns *[]define.Column,
	precm string) {

	t := reflect.TypeOf(structdata)
	v := reflect.ValueOf(structdata)

	for i := 0; i < v.NumField(); i++ {
		var (
			oldvalue interface{}
			newvalue interface{}
			comment  string
			column   define.Column
		)
		ft := t.Field(i)
		cm := ft.Tag.Get("comment")
		if cm == "" {
			continue
		}

		comment = fmt.Sprintf("%s-%s-%s", precm, t.Name(), cm)
		if isold {
			oldvalue = v.Field(i).Interface()
		} else {
			newvalue = v.Field(i).Interface()
		}
		column = define.Column{
			Name:     comment,
			OldValue: oldvalue,
			NewValue: newvalue,
		}
		*columns = append(*columns, column)
	}

}
