package email

type Client interface {
	Send(template Template, to []string, metadata map[string]string) error
}
