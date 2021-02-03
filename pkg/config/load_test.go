package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LoadConfigUnitSuite struct {
	suite.Suite
}

func (suite LoadConfigUnitSuite) TestLoad_Flags() {
	flags := []string{"logLevel", "imapHost", "smtpHost", "httpHost", "imapPort", "smtpPort", "httpPort", "enableImap", "enableSmtp", "enableHttp"}
	initFlags()
	flag.Parse()
	for _, expectedFlag := range flags {
		assert.NotNil(suite.T(), flag.Lookup(expectedFlag), "expected flag '%s', not found", expectedFlag)
	}
}

func TestLoad(t *testing.T) {
	suite.Run(t, new(LoadConfigUnitSuite))
}
