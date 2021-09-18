package email

import (
	"fmt"
	"net/smtp"
)

type SmtpClient struct {
	Host string
	Port string
	User string
	auth smtp.Auth
}

func NewSmtpClient(host string, port string, user string, password string) *SmtpClient {
	return &SmtpClient{
		Host: host,
		Port: port,
		User: user,
		auth: smtp.PlainAuth("", user, password, host),
	}
}

func (c *SmtpClient) address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *SmtpClient) Send(template Template, to []string, meta map[string]string) error {
	message, err := newMessage(template)

	if err != nil {
		return err
	}

	if err := smtp.SendMail(
		c.address(),
		c.auth,
		c.User,
		to,
		message.get(meta),
	); err != nil {
		return fmt.Errorf("not was possible to send the email %v", err)
	}

	return nil
}
