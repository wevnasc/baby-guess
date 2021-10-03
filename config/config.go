package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port     string
	Secret   string
	DBHost   string
	DBUser   string
	DBPass   string
	DBPort   string
	DBName   string
	DBSSL    string
	SMTPHost string
	SMTPUser string
	SMTPPass string
	SMTPPort string
	SMTPFrom string
}

func readFromLocal() map[string]string {
	f, err := os.Open("local.env")
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v\n", err))
	}

	scanner := bufio.NewScanner(f)
	conf := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		keyVal := strings.Split(line, "=")

		key := strings.TrimSpace(keyVal[0])
		val := strings.TrimSpace(keyVal[1])

		conf[key] = val
	}
	return conf
}

func New(local string) *Config {
	if local == "true" {
		conf := readFromLocal()

		return &Config{
			Port:     conf["PORT"],
			Secret:   conf["AUTH_SECRET"],
			DBHost:   conf["DB_HOST"],
			DBUser:   conf["DB_USER"],
			DBPass:   conf["DB_PASS"],
			DBPort:   conf["DB_PORT"],
			DBName:   conf["DB_NAME"],
			DBSSL:    conf["DB_SSL"],
			SMTPHost: conf["SMTP_HOST"],
			SMTPUser: conf["SMTP_USER"],
			SMTPPass: conf["SMTP_PASS"],
			SMTPPort: conf["SMTP_PORT"],
			SMTPFrom: conf["SMTP_FROM"],
		}
	}

	return &Config{
		Port:     os.Getenv("PORT"),
		Secret:   os.Getenv("AUTH_SECRET"),
		DBHost:   os.Getenv("DB_HOST"),
		DBUser:   os.Getenv("DB_USER"),
		DBPass:   os.Getenv("DB_PASS"),
		DBPort:   os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		DBSSL:    os.Getenv("DB_SSL"),
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		SMTPFrom: os.Getenv("SMTP_FROM"),
	}
}
