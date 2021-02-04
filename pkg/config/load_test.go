package config

import (
	"bytes"
	"errors"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"strconv"
	"testing"
)

type LoadConfigUnitSuite struct {
	suite.Suite
}

func (suite *LoadConfigUnitSuite) TestLoad_Flags() {
	flags := []string{"logLevel", "imapHost", "smtpHost", "httpHost", "imapPort", "smtpPort", "httpPort", "enableImap", "enableSmtp", "enableHttp"}
	initFlags()
	flag.Parse()
	for _, expectedFlag := range flags {
		assert.NotNil(suite.T(), flag.Lookup(expectedFlag), "expected flag '%s', not found", expectedFlag)
	}
}

func (suite *LoadConfigUnitSuite) TestParseConfig_Valid() {
	configFile := []byte(
		`
log_level: 1
networkconfigs:
    smtp:
        host: 1.2.3.4
        port: 1337
    imap:
        host: 4.3.2.1
        port: 6666
    http:
        host: 1.3.3.7
        port: 80085
enable_imap: true
enable_smtp: true
enable_http: true
`)
	config := Config{}
	expectedConfig := Config{
		LogLevel: 1,
		NetworkConfigs: struct {
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
		}{
			SMTP: struct {
				Host string `flag:"smtpHost"`
				Port int    `flag:"smtpPort"`
			}{
				Host: "1.2.3.4",
				Port: 1337,
			},
			IMAP: struct {
				Host string `flag:"imapHost"`
				Port int    `flag:"imapPort"`
			}{
				Host: "4.3.2.1",
				Port: 6666,
			},
			HTTP: struct {
				Host string `flag:"httpHost"`
				Port int    `flag:"httpPort"`
			}{
				Host: "1.3.3.7",
				Port: 80085,
			},
		},
		EnableIMAP: true,
		EnableSMTP: true,
		EnableHTTP: true,
	}
	configReader := bytes.NewReader(configFile)
	c, err := parseConfig(config, configReader)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), expectedConfig, c)
}

type CorruptionReader struct {
}

func (c CorruptionReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("ups")
}

func (suite *LoadConfigUnitSuite) TestParseConfig_ReadCorruption() {
	corruptedFile := CorruptionReader{}
	c, err := parseConfig(Config{}, corruptedFile)
	assert.Zero(suite.T(), c)
	assert.Error(suite.T(), err)
}

func (suite *LoadConfigUnitSuite) TestParseConfig_InvalidYaml() {
	corruptedYaml := []byte(
		`
log_level: 1
networkconfigs:
    smtp:
        host: 1.2.3.4
        port: 1337
    imap:
        host: 4.3.2.1
        port: 6666
    http:
        host: 1.3.3.7
        port: 80085
enable_imap: true
enable_sm
`)
	configReader := bytes.NewReader(corruptedYaml)
	c, err := parseConfig(Config{}, configReader)
	assert.Zero(suite.T(), c)
	assert.Error(suite.T(), err)
}

type MockFlagStringValue struct {
	Value string
}

func (m MockFlagStringValue) String() string {
	return m.Value
}

func (m *MockFlagStringValue) Set(s string) error {
	m.Value = s
	return nil
}

type MockFlagIntValue struct {
	Value int
}

func (m MockFlagIntValue) String() string {
	return strconv.Itoa(m.Value)
}

func (m *MockFlagIntValue) Set(s string) error {
	var err error
	m.Value, err = strconv.Atoi(s)
	return err
}

type MockFlagBoolValue struct {
	Value bool
}

func (m MockFlagBoolValue) String() string {
	return strconv.FormatBool(m.Value)
}

func (m MockFlagBoolValue) Set(s string) error {
	var err error
	m.Value, err = strconv.ParseBool(s)
	return err
}

func (suite *LoadConfigUnitSuite) TestOverrideValue_String() {
	toBeChanged := struct {
		Test string
	}{Test: "Hello World"}
	reflectField := reflect.ValueOf(&toBeChanged)
	field := reflect.Indirect(reflectField).Field(0)
	f := &flag.Flag{
		Name:     "Test",
		Usage:    "",
		Value:    &MockFlagStringValue{Value: "Test"},
		DefValue: "",
	}
	err := overrideValue(&field, f)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Test", toBeChanged.Test)
}

func (suite *LoadConfigUnitSuite) TestOverrideValue_Int() {
	toBeChanged := struct {
		Test int
	}{Test: 1}
	reflectField := reflect.ValueOf(&toBeChanged)
	field := reflect.Indirect(reflectField).Field(0)
	f := &flag.Flag{
		Name:     "Test",
		Usage:    "",
		Value:    &MockFlagIntValue{Value: 2},
		DefValue: "",
	}
	err := overrideValue(&field, f)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, toBeChanged.Test)
}

func (suite *LoadConfigUnitSuite) TestOverrideValue_Bool() {
	toBeChanged := struct {
		Test bool
	}{Test: true}
	reflectField := reflect.ValueOf(&toBeChanged)
	field := reflect.Indirect(reflectField).Field(0)
	f := &flag.Flag{
		Name:     "Test",
		Usage:    "",
		Value:    &MockFlagBoolValue{Value: false},
		DefValue: "",
	}
	err := overrideValue(&field, f)
	assert.Nil(suite.T(), err)
	assert.False(suite.T(), toBeChanged.Test)
}

func (suite *LoadConfigUnitSuite) TestParseField_AlwaysParseZeroConfig() {
	type StructWithZeroField struct {
		IwillBeZero string `flag:"TestFlag"`
	}
	instanceWithZero := StructWithZeroField{}

	reflectField := reflect.ValueOf(&instanceWithZero)
	field := reflect.Indirect(reflectField).Field(0)
	fieldType := reflect.Indirect(reflectField).Type().Field(0)

	flagSet := flag.NewFlagSet("TestParseField_AlwaysParseZeroConfig", flag.PanicOnError)
	flagSet.String("TestFlag", "Test", "")

	err := parseField(field, &fieldType, flagSet)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Test", instanceWithZero.IwillBeZero)
}

func (suite *LoadConfigUnitSuite) TestParseField_OverrideIfFlagSet() {
	type StructWithSetField struct {
		IwillBeSet string `flag:"TestFlag"`
	}
	instanceWithZero := StructWithSetField{IwillBeSet: "RemoveMe"}

	reflectField := reflect.ValueOf(&instanceWithZero)
	field := reflect.Indirect(reflectField).Field(0)
	fieldType := reflect.Indirect(reflectField).Type().Field(0)

	flagSet := flag.NewFlagSet("TestParseField_OverrideIfFlagSet", flag.PanicOnError)
	flagSet.String("TestFlag", "", "")
	err := flagSet.Set("TestFlag", "Test")
	assert.Nil(suite.T(), err)

	err = parseField(field, &fieldType, flagSet)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "Test", instanceWithZero.IwillBeSet)
}

func (suite *LoadConfigUnitSuite) TestParseField_KeepIfFlagNotSet() {
	type StructWithSetField struct {
		IwillBeSet string `flag:"TestFlag"`
	}
	instanceWithZero := StructWithSetField{IwillBeSet: "KeepMe"}

	reflectField := reflect.ValueOf(&instanceWithZero)
	fieldType := reflect.Indirect(reflectField).Type().Field(0)
	field := reflect.Indirect(reflectField).Field(0)

	flagSet := flag.NewFlagSet("TestParseField_KeepIfFlagNotSet", flag.PanicOnError)
	flagSet.String("TestFlag", "Test", "")

	err := parseField(field, &fieldType, flagSet)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "KeepMe", instanceWithZero.IwillBeSet)
}

func (suite *LoadConfigUnitSuite) TestParseField_KeepIfNoFlagtag() {
	type StructWithNoFlagField struct {
		IhaveNoTag string
	}
	instanceWithZero := StructWithNoFlagField{IhaveNoTag: "IgnoreMe"}

	reflectField := reflect.ValueOf(&instanceWithZero)
	field := reflect.Indirect(reflectField).Field(0)
	fieldType := reflect.Indirect(reflectField).Type().Field(0)

	flagSet := flag.NewFlagSet("TestParseField_KeepIfNoFlagtag", flag.PanicOnError)
	flagSet.String("TestFlag", "", "")
	err := flagSet.Set("TestFlag", "Test")
	assert.Nil(suite.T(), err)

	err = parseField(field, &fieldType, flagSet)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "IgnoreMe", instanceWithZero.IhaveNoTag)
}

func TestLoad(t *testing.T) {
	suite.Run(t, new(LoadConfigUnitSuite))
}
