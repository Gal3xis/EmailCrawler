package conf

import (
	"EmailCrawler/src/paths"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	Url      string
	Port     int
	Username string
	Password string
}

type MailboxConfig struct {
	Mailbox         string
	SaveFolder      string
	SavingStructure string
	MailOffset      int
	DeleteMails     bool
	MinAgeInDays    int
	MinEmailsToKeep int
}

func (c *Config) Read() error {
	content, err := os.ReadFile(filepath.Join(paths.ConfigPath, "emailCrawlerConfig.conf"))
	if err != nil {
		return err
	}
	tokens, err := tokenize(string(content))
	if err != nil {
		return err
	}
	err = c.parseTokensToConfig(tokens)
	if err != nil {
		return err
	}
	return nil
}

func tokenize(fileContent string) (map[string][]string, error) {
	tokens := map[string][]string{}
	lines := strings.Split(fileContent, "\n")
	var tag string = ""
	var configLines []string
	for _, line := range lines {
		if len(line) == 0 || line[0] == ';' {
			continue
		}
		line = strings.TrimSpace(line)
		line = strings.ReplaceAll(line, " ", "")
		if len(line) == 0 {
			continue
		}
		if line[0] == '[' {
			if tag != "" {
				tokens[tag] = configLines
				tag = ""
				configLines = []string{}
			}
			regex, err := regexp.Compile(`\[(.*?)]`)
			if err != nil {
				return nil, err
			}
			tags := regex.FindStringSubmatch(line)[1:]
			if len(tags) > 1 {
				return nil, TooManyTagsError{
					Config: "Mailboxes",
					Tags:   tags,
				}
			}
			if len(tags) == 0 {
				return nil, EmptyTagError{
					Config: "Mailboxes",
				}
			}
			tag = tags[0]
			continue
		}
		configLines = append(configLines, line)
	}
	if tag != "" {
		tokens[tag] = configLines
	}
	return tokens, nil
}

func (c *Config) parseTokensToConfig(tokens map[string][]string) error {
	for key, value := range tokens {
		switch key {
		case "Connection":
			err := c.parseConnection(value)
			if err != nil {
				return err
			}
		default:
			err := c.parseMailbox(key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) parseConnection(lines []string) error {
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		value := parts[1]
		switch key {
		case "Url":
			c.Url = value
		case "Port":
			var err error
			if c.Port, err = strconv.Atoi(value); err != nil {
				return err
			}
		case "Username":
			c.Username = value
		case "Password":
			c.Password = value
		default:
			return InvalidConfigKeyError{
				Config: "Connection",
				Key:    value,
			}
		}
	}
	return nil
}

func (c *Config) parseMailbox(name string, lines []string) error {
	mailboxConfig := MailboxConfig{
		Mailbox: name,
	}
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		key := parts[0]
		value := parts[1]
		switch key {
		case "MailOffset":
			if valueInt, err := strconv.Atoi(value); err != nil {
				return err
			} else {
				mailboxConfig.MailOffset = valueInt + 1
			}
		case "SaveFolder":
			mailboxConfig.SaveFolder = value
		case "SavingStructure":
			mailboxConfig.SavingStructure = value
		case "DeleteMails":
			var err error
			if mailboxConfig.DeleteMails, err = strconv.ParseBool(value); err != nil {
				return err
			}
		case "MinAgeInDaysToDelete":
			var err error
			if mailboxConfig.MinAgeInDays, err = strconv.Atoi(value); err != nil {
				return err
			}
		case "MinEmailsToKeep":
			var err error
			if mailboxConfig.MinEmailsToKeep, err = strconv.Atoi(value); err != nil {
				return err
			}
		default:
			return InvalidConfigKeyError{
				Config: "Connection",
				Key:    value,
			}
		}
	}
	if c.Mailboxes == nil {
		c.Mailboxes = map[string]MailboxConfig{}
	}
	c.Mailboxes[name] = mailboxConfig
	return nil
}

func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

type EmptyTagError struct {
	Config string
}

func (e EmptyTagError) Error() string {
	return fmt.Sprintf("Tag cannot be empty.")
}

type TooManyTagsError struct {
	Config string
	Tags   []string
}

func (e TooManyTagsError) Error() string {
	return fmt.Sprintf("Only one Tag is allowed, but found %d. %s", len(e.Tags), e.Tags)
}

type InvalidConfigKeyError struct {
	Config string
	Key    string
}

func (e InvalidConfigKeyError) Error() string {
	return fmt.Sprintf("Key %s not known in Config %s", e.Key, e.Config)
}
