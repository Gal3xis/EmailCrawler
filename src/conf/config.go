package conf

import (
	"EmailCrawler/src/paths"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ConfigReader interface {
	Read(path string) error
}

type Config struct {
	Connection
	Mailboxes map[string]MailboxConfig
}

type Connection struct {
	url      string
	port     int
	Username string
	Password string
}

type MailboxConfig struct {
	Mailbox         string
	DeleteMails     bool
	MinAgeInDays    int
	MinEmailsToKeep int
}

func (c *Config) Read() error {
	content, err := os.ReadFile(filepath.Join(paths.ConfigPath, "connection.conf"))
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line[0] == '[' {
			continue
		}
		line = strings.ReplaceAll(line, " ", "")
		line = strings.TrimSpace(line)
		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		value := parts[1]
		switch key {
		case "url":
			c.url = value
		case "port":
			c.port, err = strconv.Atoi(value)
			if err != nil {
				return err
			}
		case "username":
			c.Username = value
		case "password":
			c.Password = value
		default:
			return WrongConfigValueError{
				Config: "Connection",
				Value:  value,
			}
		}
	}
	return nil
}

type WrongConfigValueError struct {
	Config string
	Value  string
}

func (e WrongConfigValueError) Error() string {
	return fmt.Sprintf("Value %s not known in Config %s", e.Value, e.Config)
}
