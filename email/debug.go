package email

import (
	"fmt"
	"strings"
)

type DebugClient struct{}

func (c *DebugClient) Send(template Template, to []string, metadata map[string]string) error {
	email, err := newMessage("debug@example.com.br", template)

	if err != nil {
		fmt.Printf("%v \n", err)
		return nil
	}

	fmt.Printf("\nTo: %v\n", to)

	lines := strings.Split(string(email.get(metadata)), "\n")

	for _, line := range lines {
		fmt.Println(line)
	}

	fmt.Println("----------------------------------------")

	return nil
}
