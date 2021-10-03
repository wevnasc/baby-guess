package email

import (
	"fmt"
	"net/smtp"

	"github.com/wevnasc/baby-guess/config"
)

type SmtpClient struct {
	Host string
	Port string
	User string
	From string
	auth smtp.Auth
}

func NewSmtpClient(config *config.Config) *SmtpClient {
	return &SmtpClient{
		Host: config.SMTPHost,
		Port: config.SMTPPort,
		User: config.SMTPUser,
		From: config.SMTPFrom,
		auth: smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, config.SMTPHost),
	}
}

func (c *SmtpClient) address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *SmtpClient) Send(template Template, to []string, meta map[string]string) error {
	message, err := newMessage(c.From, template)

	if err != nil {
		return err
	}

	if err := smtp.SendMail(
		c.address(),
		c.auth,
		c.From,
		to,
		message.get(meta),
	); err != nil {
		return fmt.Errorf("not was possible to send the email %v", err)
	}

	return nil
}
