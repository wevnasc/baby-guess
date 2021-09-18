package email

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type message struct {
	body    string
	subject string
}

func newMessage(template Template) (*message, error) {
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

	return &message{
		body:    string(body),
		subject: fmt.Sprintf("Subject:%s\n", subject),
	}, nil

}

func (e *message) get(meta map[string]string) []byte {
	subject := e.subject
	body := e.body

	for key, value := range meta {
		pattern := fmt.Sprintf("{%s}", key)
		subject = strings.ReplaceAll(subject, pattern, value)
		body = strings.ReplaceAll(body, pattern, value)
	}

	return []byte(subject + body)
}
