package logiacore

import (
	"fmt"
	"github.com/yusologia/go-core/v2/pkg"
	"gopkg.in/mail.v2"
	"os"
)

type MailMessage interface {
	Message() *mail.Message
}

type MailConf struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Mail struct {
	queue  string
	dialer *mail.Dialer
}

func (m *Mail) Dial(conf MailConf) *Mail {
	m.dialer = mail.NewDialer(conf.Host, logiapkg.ToInt(conf.Port), conf.Username, conf.Password)

	return m
}

func (m *Mail) Send(msg MailMessage) error {
	content := msg.Message()
	content.SetHeader("From", content.FormatAddress(os.Getenv("MAIL_FROM_ADDRESS"), os.Getenv("MAIL_FROM_NAME")))

	if err := m.dialer.DialAndSend(content); err != nil {
		logiapkg.LogError(fmt.Sprintf("Error sending email: %v", err))
		return err
	}

	return nil
}
