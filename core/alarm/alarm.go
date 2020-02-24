package alarm

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"
	"strconv"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/notify"
	"github.com/labulaka521/crocodile/common/notify/dingding"
	"github.com/labulaka521/crocodile/common/notify/email"
	"github.com/labulaka521/crocodile/common/notify/slack"
	"github.com/labulaka521/crocodile/common/notify/telegram"
	"github.com/labulaka521/crocodile/common/notify/wechat"
	"github.com/labulaka521/crocodile/common/utils"
	"github.com/labulaka521/crocodile/core/config"
	"github.com/labulaka521/crocodile/core/model"
	"github.com/labulaka521/crocodile/core/utils/define"
	"go.uber.org/zap"
)

type sendNotify struct {
	dingding, email, slack,
	telegram, wechat notify.Sender
}

var send sendNotify

// maintask name[id]
// taskrunstatus fail/success
// starttime
// endtime
// totalruntime
//
// errMsg
// errTaskTypeStr
// errTask name[id]

const (
	title     = "任务通知 {{ .TaskName }}[{{ .TaskID }}]"
	alarmtmpl = `任务名称    : {{ .TaskName }}
任务ID      : {{ .TaskID }}
开始时间   : {{ .StartTime }}
结束时间     : {{ .EndTime }}
总运行时间: {{ .TotalRuntime }}
状态      : {{ .Status -}}
{{if eq .Status "fail" }}
错误任务名称 : {{ .ErrTaskName }} {{ .ErrTasktypestr }}
错误任务ID   : {{ .ErrTaskID }}
错误信息      : {{ .ErrMsg }}
{{- end }}`
)

// InitAlarm init alarm notify
func InitAlarm() {
	log.Info("start init alarm")
	send = sendNotify{}
	notifycfg := config.CoreConf.Notify
	if notifycfg.DingDing.Enable {
		log.Debug("init dingding alarm")
		send.dingding = dingding.NewDing(notifycfg.DingDing.WebHook,
			dingding.Secrue(notifycfg.DingDing.SecureLevel),
			notifycfg.DingDing.Secret)
	}

	if notifycfg.Email.Enable {
		log.Debug("init email alarm")
		send.email = email.NewSMTP(notifycfg.Email.SMTPHost,
			notifycfg.Email.Port,
			notifycfg.Email.UserName,
			notifycfg.Email.Password,
			notifycfg.Email.From,
			notifycfg.Email.TLS,
			notifycfg.Email.Anonymous,
			notifycfg.Email.SkipVerify,
		)
	}

	if notifycfg.Slack.Enable {
		log.Debug("init slack alarm")
		send.slack = slack.NewSlack(notifycfg.Slack.WebHook)
	}

	if notifycfg.Telegram.Enable {
		log.Debug("init telegram alarm")
		var err error
		send.telegram, err = telegram.NewTelegram(notifycfg.Telegram.BotToken)
		if err != nil {
			log.Error("New telegram Bot failed", zap.Error(err))
		}
	}

	if notifycfg.WeChat.Enable {
		log.Debug("init wechat alarm")
		send.wechat = wechat.NewWeChat(notifycfg.WeChat.CropID, notifycfg.WeChat.AgentID, notifycfg.WeChat.AgentSecret)
	}
}

// JudgeNotify send notify to user
func JudgeNotify(tasklog *define.Log) {
	taskdata, err := model.GetTaskByID(context.Background(), tasklog.RunByTaskID)
	if err != nil {
		log.Error("get task failed", zap.String("taskid", tasklog.RunByTaskID), zap.Error(err))
		return
	}
	// Check this task alarm
	// if task alarmsttaus equal -2,it will alarm when task run finish
	// whether alarm when task alarmstatus equal task run resp statsu

	var status string
	if tasklog.Status == -1 {
		status = "fail"
	} else if tasklog.Status == 1 {
		status = "success"
	} else {
		status = "unknow"
	}
	totalruntime := strconv.Itoa(tasklog.TotalRunTime) + "ms"
	if tasklog.TotalRunTime > 1000 {
		totalruntime = strconv.Itoa(tasklog.TotalRunTime/1000) + "s"
	}

	if int(taskdata.AlarmStatus) == -2 || int(taskdata.AlarmStatus) == tasklog.Status {
		err := sendalarm(taskdata.AlarmUserIds,
			taskdata.Name,
			tasklog.RunByTaskID,
			utils.UnixToStr(tasklog.StartTime/1e3),
			utils.UnixToStr(tasklog.EndTime/1e3),
			status,
			totalruntime,
			tasklog.ErrMsg,
			tasklog.ErrTaskTypeStr,
			tasklog.ErrTask,
			tasklog.ErrTaskID,
		)
		if err != nil {
			log.Error("sendalarm failed", zap.Error(err))
		}
	}
}

type notifymsg struct {
	TaskID         string   `json:"task_id"`
	TaskName       string   `json:"task_name"`
	StartTime      string   `json:"start_time"`
	EndTime        string   `json:"end_time"`
	Status         string   `json:"status"`
	TotalRuntime   string   `json:"total_runtime"`
	AlarmUsers     []string `json:"alarm_users"`
	ErrTaskID      string   `json:"err_taskid,omitempty"`
	ErrTaskName    string   `json:"err_taskname,omitempty"`
	ErrTasktypestr string   `json:"err_tasktype,omitempty"`
	ErrMsg         string   `json:"err_msg,omitempty"`
}

// sendalarm will send notify to task's alarm users
func sendalarm(notifyuids []string, taskname, taskid, starttime, endtime, status, totalruntime,
	errmsg, errtasktypestr, errtaskname, errtaskid string) error {
	log.Info("start send alarm", zap.Strings("uids", notifyuids), zap.String("task", taskname))
	if len(notifyuids) == 0 {
		err := fmt.Errorf("task %s[%s] not exist alarm users", taskname, taskid)
		return err
	}
	alarmusers, _, err := model.GetUsers(context.Background(), notifyuids, 0, 0)
	if err != nil {
		log.Error("get alarm user info failed", zap.Error(err))
		return err
	}
	alarmUsernNames := make([]string, 0, len(alarmusers))
	alarmEmail := make([]string, 0, len(alarmusers))
	alarmWeChat := make([]string, 0, len(alarmusers))
	alarmDingDing := make([]string, 0, len(alarmusers))
	alarmSlack := make([]string, 0, len(alarmusers)) // tg
	alarmTelegram := make([]string, 0, len(alarmusers))

	for _, user := range alarmusers {
		if user.Email != "" {
			alarmEmail = append(alarmEmail, user.Email)
		}
		if user.WeChat != "" {
			alarmWeChat = append(alarmWeChat, user.WeChat)
		}
		if user.DingPhone != "" {
			alarmDingDing = append(alarmDingDing, user.DingPhone)
		}
		if user.Slack != "" {
			alarmSlack = append(alarmSlack, user.Slack)
		}
		if user.Telegram != "" {
			alarmTelegram = append(alarmTelegram, user.Telegram)
		}
	}

	notifymsg := notifymsg{
		TaskName:       taskname,
		TaskID:         taskid,
		StartTime:      starttime,
		EndTime:        endtime,
		Status:         status,
		TotalRuntime:   totalruntime,
		AlarmUsers:     alarmUsernNames,
		ErrMsg:         errmsg,
		ErrTasktypestr: errtasktypestr,
		ErrTaskName:    errtaskname,
		ErrTaskID:      errtaskid,
	}

	// TODO problem
	// send webhook
	if config.CoreConf.Notify.WebHook.Enable {
		_, err := notify.JSONPost(http.MethodPost, config.CoreConf.Notify.WebHook.WebHookURL, notifymsg, http.DefaultClient)
		if err != nil {
			log.Error("send webhook failed",
				zap.String("webhookurl", config.CoreConf.Notify.WebHook.WebHookURL),
				zap.Error(err),
			)
		}
	}

	alarmtitle, err := template.New("title").Parse(title)
	if err != nil {
		log.Error("template new title failed", zap.Error(err))
		return err
	}

	var titlebuf = &bytes.Buffer{}
	err = alarmtitle.Execute(titlebuf, notifymsg)
	if err != nil {
		log.Error("title template new title failed", zap.Error(err))
		return err
	}

	alarmcontent, err := template.New("content").Parse(alarmtmpl)
	if err != nil {
		log.Error("template new title failed", zap.Error(err))
		return err
	}

	var contentbuf = &bytes.Buffer{}
	err = alarmcontent.Execute(contentbuf, notifymsg)
	if err != nil {
		log.Error("title template new title failed", zap.Error(err))
		return err
	}

	for _, uid := range notifyuids {
		notify := define.Notify{
			NotifyType: define.TaskNotify,
			NotifyUID:  uid,
			Title:      taskname,
			Content:    contentbuf.String(),
			NotifyTime: time.Now().Unix(),
		}
		err = model.SaveNewNotify(context.Background(), notify)
		if err != nil {
			log.Error("model.SaveNewNotify", zap.Error(err))
		}
	}

	// send notify alarm
	if send.dingding != nil && len(alarmDingDing) != 0 {
		err := send.dingding.Send(alarmDingDing, titlebuf.String(), contentbuf.String())
		if err != nil {
			log.Error("send dingding notify failed", zap.Error(err))
		}
	}

	if send.email != nil && len(alarmEmail) != 0 {
		err := send.email.Send(alarmEmail, titlebuf.String(), contentbuf.String())
		if err != nil {
			log.Error("send email notify failed", zap.Error(err))
		}
	}

	if send.wechat != nil && len(alarmWeChat) != 0 {
		err := send.wechat.Send(alarmWeChat, titlebuf.String(), contentbuf.String())
		if err != nil {
			log.Error("send wechat notify failed", zap.Error(err))
		}
	}

	if send.slack != nil && len(alarmSlack) != 0 {
		err := send.slack.Send(alarmSlack, titlebuf.String(), contentbuf.String())
		if err != nil {
			log.Error("send slack notify failed", zap.Error(err))
		}
	}

	if send.telegram != nil && len(alarmTelegram) != 0 {
		err := send.telegram.Send(alarmTelegram, titlebuf.String(), contentbuf.String())
		if err != nil {
			log.Error("send telegram notify failed", zap.Error(err))
		}
	}
	return nil
}
