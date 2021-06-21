package config

import (
	"github.com/sirupsen/logrus"
)

var configuration Config

type Config struct {
	LogrusLevel    logrus.Level `yaml:"-"`
	LogLevel       int          `yaml:"log_level" flag:"logLevel"`
	NetworkConfigs struct {
		SMTP struct {
			Host string `flag:"smtpHost"`
			Port int    `flag:"smtpPort"`
		}
		IMAP struct {
			Host string `flag:"imapHost"`
			Port int    `flag:"imapPort"`
		}
		HTTP struct {
			Host string `flag:"httpHost"`
			Port int    `flag:"httpPort"`
		}
		API struct {
			Host string `flag:"apiHost"`
			Port int    `flag:"apiPort"`
		}
	}
	DisableIMAP bool `yaml:"disable_imap" flag:"disableImap"`
	DisableSMTP bool `yaml:"disable_smtp" flag:"disableSmtp"`
	DisableHTTP bool `yaml:"disable_http" flag:"disableHttp"`
}

func GetConfig() Config {
	return configuration
}
