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
	"github.com/labulaka521/crocodile/core/utils/resp"
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
		// operate id
		// 操作人id，姓名，用户类型
		// Method 创建，修改，删除
		// 模块: 用户，主机组，主机，任务，用户
		// 操作对象名称
		// 修改时间
		// 	修改字段
		// 	旧值
		// 	新值

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
			log.Error("ioutil.ReadAll", zap.Error(err))
			resp.JSON(c, resp.ErrInternalServer, nil)
			return
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
		ctx, cancel := context.WithTimeout(context.Background(),
			config.CoreConf.Server.DB.MaxQueryTime.Duration)
		defer cancel()

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
					resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				modulename = userData.Name
				oldData = *userData
			case hostgroup:
				hostgroupData, err := model.GetHostGroupByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostGroupByID", zap.String("id", id), zap.Error(err))
					resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				modulename = hostgroupData.Name
				oldData = *hostgroupData
			case host:
				hostData, err := model.GetHostByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				modulename = hostData.Addr
				oldData = *hostData
			case task:
				taskData, err := model.GetTaskByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					resp.JSON(c, resp.ErrInternalServer, nil)
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
					// resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				modulename = userData.Name
				newData = *userData
			case hostgroup:
				hostgroupData, err := model.GetHostGroupByName(ctx, name)
				if err != nil {
					log.Error("model.GetHostGroupByName", zap.String("name", name), zap.Error(err))
					// resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				modulename = hostgroupData.Name
				newData = *hostgroupData
			case task:
				taskData, err := model.GetTaskByName(ctx, name)
				if err != nil {
					log.Error("model.GetTaskByName", zap.String("name", name), zap.Error(err))
					// resp.JSON(c, resp.ErrInternalServer, nil)
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
					// resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				newData = *userData
			case hostgroup:
				hostgroupData, err := model.GetHostGroupByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostGroupByID", zap.String("id", id), zap.Error(err))
					// resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				newData = *hostgroupData
			case host:
				hostData, err := model.GetHostByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					// resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				newData = *hostData
			case task:
				taskData, err := model.GetTaskByID(ctx, id)
				if err != nil {
					log.Error("model.GetHostByID", zap.String("id", id), zap.Error(err))
					// resp.JSON(c, resp.ErrInternalServer, nil)
					return
				}
				newData = *taskData
			default:
				log.Debug("can get module name from url", zap.String("url", c.Request.URL.Path))
				c.Next()
				return
			}
		}

		columns := make([]define.Column, 0, 50)

		parseColumn(oldData, newData, &columns, "")

		if len(columns) == 0 {
			return
		}

		uid := c.GetString("uid")
		username := c.GetString("username")
		// 获取用户的类型
		var role define.Role
		if v, ok := c.Get("role"); ok {
			role = v.(define.Role)
		}

		err = model.SaveOperateLog(ctx,
			uid,
			username,
			role,
			c.Request.Method,
			module,
			modulename,
			operatetimne,
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
			comment = fmt.Sprintf("%s-%s-%s", precm, newt.Name(), comment)
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
