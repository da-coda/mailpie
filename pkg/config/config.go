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
	}
	EnableIMAP bool `yaml:"enable_imap" flag:"enableImap"`
	EnableSMTP bool `yaml:"enable_smtp" flag:"enableSmtp"`
	EnableHTTP bool `yaml:"enable_http" flag:"enableHttp"`
}

func GetConfig() Config {
	return configuration
}
