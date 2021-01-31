package config

import "github.com/sirupsen/logrus"

var configuration Config

type Config struct {
	LogLevel       logrus.Level
	logLevel       int
	NetworkConfigs NetworkConfigs
	EnableIMAP     bool
	EnableSMTP     bool
	EnableHTTP     bool
}

type NetworkConfigs struct {
	SMTPHost string
	IMAPHost string
	HTTPHost string

	SMTPPort int
	IMAPPort int
	HTTPPort int
}

func GetConfig() Config {
	configuration.LogLevel = logrus.Level(configuration.logLevel)
	return configuration
}
