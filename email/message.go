package email

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type message struct {
	body   string
	header string
}

func newMessage(from string, template Template) (*message, error) {
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
		body:   string(body),
		header: fmt.Sprintf("From:%s\nSubject:%s\n", from, subject),
	}, nil

}

func (e *message) get(meta map[string]string) []byte {
	header := e.header
	body := e.body

	for key, value := range meta {
		pattern := fmt.Sprintf("{%s}", key)
		header = strings.ReplaceAll(header, pattern, value)
		body = strings.ReplaceAll(body, pattern, value)
	}

	return []byte(header + body)
}
