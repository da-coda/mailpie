package config

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	configuration = Config{LogLevel: int(logrus.DebugLevel)}
	config := GetConfig()
	assert.Equal(t, logrus.DebugLevel, logrus.Level(config.LogLevel))
}
