package email

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
)

type Template string

const (
	AccountCreated Template = "account_created"
	ItemSelected            = "item_selected"
	ItemApproved            = "item_approved"
	Winner                  = "winner"
	Losener                 = "looser"
)

type Connection struct {
	Host string
	Port string
	User string
	auth smtp.Auth
}

func NewConnection(host string, port string, user string, password string) *Connection {
	return &Connection{
		Host: host,
		Port: port,
		User: user,
		auth: smtp.PlainAuth("", user, password, host),
	}
}

func (c *Connection) address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type Email struct {
	From       string
	Subject    string
	Body       string
	connection *Connection
}

func New(connection *Connection, subject string, body string) *Email {
	return &Email{
		From:       connection.User,
		Subject:    subject,
		Body:       body,
		connection: connection,
	}
}

func NewFromTemplate(connection *Connection, template Template) (*Email, error) {
	bfile := fmt.Sprintf("email/templates/%s_body.txt", template)
	body, err := ioutil.ReadFile(bfile)

	if err != nil {
		return nil, fmt.Errorf("error to create the new email %v", err)
	}

	sfile := fmt.Sprintf("email/templates/%s_subject.txt", template)
	subject, err := ioutil.ReadFile(sfile)

	if err != nil {
		return nil, fmt.Errorf("error to create the new email %v", err)
	}

	return &Email{
		From:       connection.User,
		Body:       string(body),
		Subject:    string(subject),
		connection: connection,
	}, nil

}

func (e *Email) buildMessage(meta map[string]string) []byte {
	subject := fmt.Sprintf("Subject:%s\n", e.Subject)
	body := e.Body

	for key, value := range meta {
		pattern := fmt.Sprintf("{%s}", key)
		subject = strings.ReplaceAll(subject, pattern, value)
		body = strings.ReplaceAll(body, pattern, value)
	}

	return []byte(subject + body)
}

func (e *Email) Send(to []string, meta map[string]string) error {
	if err := smtp.SendMail(
		e.connection.address(),
		e.connection.auth,
		e.From,
		to,
		e.buildMessage(meta),
	); err != nil {
		return fmt.Errorf("not was possible to send the email %v", err)
	}

	return nil
}
