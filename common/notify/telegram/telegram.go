package telegram

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labulaka521/crocodile/common/notify"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	once sync.Once
)

// Telegram conf
type Telegram struct {
	token string
	bot   *tb.Bot
}

// NewTelegram init telegram
func NewTelegram(token string) (notify.Sender, error) {
	telegram := &Telegram{
		token: token,
	}
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	bot, err := tb.NewBot(
		tb.Settings{
			Token:  token,
			Poller: &tb.LongPoller{Timeout: time.Second * 10},
			Client: client,
		})
	if err != nil {
		return nil, err
	}
	telegram.bot = bot

	go func() {
		once.Do(func() {
			bot.Start()
		})
	}()
	// wait bot start
	time.Sleep(time.Second)
	return telegram, nil
}

// Send will send notify to channel
func (t *Telegram) Send(tos []string, title string, content string) error {
	for _, id := range tos {
		uid, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		go func(uid int) {
			t.bot.Send(&tb.User{ID: uid}, title+"\n"+title)
		}(uid)
	}
	return nil
}
