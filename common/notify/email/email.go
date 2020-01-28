package email

import (
	"crypto/tls"

	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/labulaka521/crocodile/common/log"
	"github.com/labulaka521/crocodile/common/notify"
	"github.com/labulaka521/crocodile/common/utils"
	"go.uber.org/zap"
)

// SMTP is email conf
type SMTP struct {
	SMTPHost   string
	Port       int
	Username   string
	Password   string
	From       string
	TLS        bool
	Anonymous  bool
	SkipVerify bool
}

// NewSMTP return a tls Smtp
func NewSMTP(smtphost string, port int, username, password, from string, tls, anonymous, skipVerify bool) notify.Sender {
	return &SMTP{
		SMTPHost:   smtphost,
		Username:   username,
		Password:   password,
		From:       from,
		TLS:        tls,
		Port:       port,
		Anonymous:  anonymous,
		SkipVerify: skipVerify,
	}
}

// Send send email to user
func (s *SMTP) Send(tos []string, title, content string) error {
	if s.SMTPHost == "" {
		return fmt.Errorf("address is necessary")
	}
	safetos := []string{}
	for _, to := range tos {
		err := utils.CheckEmail(to)
		if err != nil {
			log.Error("email check error", zap.Error(err))
			continue
		}
		safetos = append(safetos, to)
	}

	toaddr := strings.Join(safetos, ";")

	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	header := make(map[string]string)
	header["From"] = s.From
	header["To"] = toaddr
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte(title)))
	header["MIME-Version"] = "1.0"

	header["Content-Type"] = "text/plain"
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(content))

	var auth smtp.Auth = nil
	if !s.Anonymous {
		auth = smtp.PlainAuth("", s.Username, s.Password, s.SMTPHost)
	}
	return s.sendMail(auth, safetos, []byte(message))
}

// sendMail will send mail to user
func (s *SMTP) sendMail(auth smtp.Auth, to []string, msg []byte) (err error) {
	if err := validateLine(s.From); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}
	var client *smtp.Client
	addr := fmt.Sprintf("%s:%d", s.SMTPHost, s.Port)
	if s.TLS {
		// httpsProxyURI, _ := url.Parse("https://your https proxy:443")
		// httpsDialer, err := proxy.FromURL(httpsProxyURI, HttpsDialer)

		tlsconfig := &tls.Config{
			InsecureSkipVerify: s.SkipVerify,
			ServerName:         s.SMTPHost,
		}
		var c *tls.Conn
		c, err = tls.Dial("tcp", addr, tlsconfig)

		if err != nil {
			return err
		}

		// tls.DialWithDialer(dialer *net.Dialer, network string, addr string, config *tls.Config)
		client, err = smtp.NewClient(c, s.SMTPHost)
		if err != nil {
			return err
		}

		defer client.Close()
	} else {
		client, err = smtp.Dial(addr)
		if err != nil {
			return err
		}

		defer client.Close()

		if ok, _ := client.Extension("STARTTLS"); ok {
			config := &tls.Config{
				InsecureSkipVerify: s.SkipVerify,
				ServerName:         s.SMTPHost,
			}
			if err = client.StartTLS(config); err != nil {
				return err
			}
		}
	}
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}
	if err = client.Mail(s.From); err != nil {
		return err
	}
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()

}

func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return fmt.Errorf("smtp: A line must not contain CR or LF")
	}
	return nil
}
