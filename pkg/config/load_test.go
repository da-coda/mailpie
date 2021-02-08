package config

import (
	"bytes"
	"errors"
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"reflect"
	"strconv"
	"testing"
)

type LoadConfigUnitSuite struct {
	suite.Suite
}

func (suite *LoadConfigUnitSuite) SetupTest() {
	configuration = Config{}
}

func (suite *LoadConfigUnitSuite) TestLoad_ConfigExists() {
	tmp := os.TempDir()
	confFile, err := os.CreateTemp(tmp, "mailpieTestLoadConfigExists")
	if err != nil {
		suite.T().Skip("Unable to create temp config file for test")
	}
	configFileContent := []byte(
		`
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
disable_imap: false
disable_smtp: false
disable_http: false
`)
	_, err = confFile.Write(configFileContent)
	if err != nil {
		suite.T().Skip("Unable to write temp config file for test")
	}
	err = confFile.Close()
	if err != nil {
		suite.T().Skip("Unable to close temp config file after write")
	}
	flags := flag.NewFlagSet("TestLoad", flag.PanicOnError)
	arguments := []string{"-config", confFile.Name(), "-disableImap", "-smtpHost", "8.8.8.8", "-imapPort", "9999"}
	err = Load(flags, arguments)
	suite.Nil(err)
	config := GetConfig()
	//use config
	suite.Equal(80085, config.NetworkConfigs.HTTP.Port)
	//use flag
	suite.Equal("8.8.8.8", config.NetworkConfigs.SMTP.Host)
	//use flag
	suite.Equal(9999, config.NetworkConfigs.IMAP.Port)
	//use default
	suite.Equal(int(logrus.WarnLevel), config.LogLevel)
	//parse bool flags correctly
	suite.Equal(true, config.DisableIMAP)
	_ = os.Remove(confFile.Name())
}

func (suite *LoadConfigUnitSuite) TestLoad_ConfigNotExists() {
	tmp := os.TempDir()
	flags := flag.NewFlagSet("TestLoad", flag.PanicOnError)
	arguments := []string{"-config", tmp + "/mailpie_test.yml", "-disableImap", "-smtpHost", "8.8.8.8", "-imapPort", "9999"}
	err := Load(flags, arguments)
	suite.Nil(err)
	open, err := os.Open("/tmp/mailpie_test.yml")
	suite.Nil(err)
	var conf Config
	content, err := io.ReadAll(open)
	suite.Nil(err)
	err = yaml.Unmarshal(content, &conf)
	suite.Nil(err)
	_ = os.Remove("/tmp/mailpie_test.yml")
}

//
// initFlags
//

func (suite *LoadConfigUnitSuite) TestInitFlags() {
	flags := []string{"logLevel", "imapHost", "smtpHost", "httpHost", "imapPort", "smtpPort", "httpPort", "disableImap", "disableSmtp", "disableHttp"}
	flagSet := flag.NewFlagSet("TestInitFlags", flag.PanicOnError)
	initFlags(flagSet)
	err := flagSet.Parse([]string{})
	suite.Nil(err)
	for _, expectedFlag := range flags {
		suite.NotNil(flagSet.Lookup(expectedFlag), "expected flag '%s', not found", expectedFlag)
	}
}

//
//	parseConfig
//

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
disable_imap: true
disable_smtp: true
disable_http: true
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
		DisableIMAP: true,
		DisableSMTP: true,
		DisableHTTP: true,
	}
	configReader := bytes.NewReader(configFile)
	c, err := parseConfig(config, configReader)
	suite.Nil(err)
	suite.Equal(expectedConfig, c)
}

type CorruptionReader struct {
}

func (c CorruptionReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("ups")
}

func (suite *LoadConfigUnitSuite) TestParseConfig_ReadCorruption() {
	corruptedFile := CorruptionReader{}
	c, err := parseConfig(Config{}, corruptedFile)
	suite.Zero(c)
	suite.Error(err)
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
disable_imap: true
disable_sm
`)
	configReader := bytes.NewReader(corruptedYaml)
	c, err := parseConfig(Config{}, configReader)
	suite.Zero(c)
	suite.Error(err)
}

//
// combineConfigAndFlags
//

func (suite *LoadConfigUnitSuite) TestCombineConfigAndFlags() {
	config := Config{DisableSMTP: true, DisableHTTP: false}
	config.NetworkConfigs.HTTP.Port = 8080
	config.NetworkConfigs.SMTP.Host = "0.0.0.0"

	flags := flag.NewFlagSet("TestCombineConfigAndFlags", flag.PanicOnError)
	initFlags(flags)
	err := flags.Set("httpPort", "9000")
	suite.Nil(err)
	err = flags.Set("httpHost", "127.0.0.1")
	suite.Nil(err)
	err = flags.Set("disableSmtp", "false")

	configAndFlags, err := combineConfigAndFlags(config, flags)
	suite.Nil(err)
	//use default value of flag if not set in config
	suite.Equal(false, configAndFlags.DisableIMAP, "Should use default value of flags if not set in config")
	//handle nested structs
	suite.Equal("127.0.0.1", configAndFlags.NetworkConfigs.HTTP.Host, "Should handle nested structs")
	//handle nested structs; flags have higher prio then config
	suite.Equal(9000, configAndFlags.NetworkConfigs.HTTP.Port, "In nested structs: Flags have higher prio then config")
	//flags have higher prio then config
	suite.Equal(false, configAndFlags.DisableSMTP, "Flags have higher prio then config")
	//keep config if flag is not set
	suite.Equal("0.0.0.0", configAndFlags.NetworkConfigs.SMTP.Host, "use config value if flag not set")
}

//
// overrideValue
//

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
	suite.Nil(err)
	suite.Equal("Test", toBeChanged.Test)
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
	suite.Nil(err)
	suite.Equal(2, toBeChanged.Test)
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
	suite.Nil(err)
	suite.False(toBeChanged.Test)
}

//
// parseField
//

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
	suite.Nil(err)
	suite.Equal("Test", instanceWithZero.IwillBeZero)
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
	suite.Nil(err)

	err = parseField(field, &fieldType, flagSet)
	suite.Nil(err)
	suite.Equal("Test", instanceWithZero.IwillBeSet)
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
	suite.Nil(err)
	suite.Equal("KeepMe", instanceWithZero.IwillBeSet)
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
	suite.Nil(err)

	err = parseField(field, &fieldType, flagSet)
	suite.Nil(err)
	suite.Equal("IgnoreMe", instanceWithZero.IhaveNoTag)
}

//
// parseStruct
//

func (suite *LoadConfigUnitSuite) TestParseStruct_NestedStruct() {
	type MainStruct struct {
		SubStruct struct {
			LeaveStruct struct {
				ParseMeSenpai string `flag:"parseMeSenpai"`
			}
		}
	}
	toBeParsed := MainStruct{
		SubStruct: struct {
			LeaveStruct struct {
				ParseMeSenpai string `flag:"parseMeSenpai"`
			}
		}{LeaveStruct: struct {
			ParseMeSenpai string `flag:"parseMeSenpai"`
		}{ParseMeSenpai: "ParseMe"}},
	}

	reflectField := reflect.ValueOf(&toBeParsed)
	field := reflect.Indirect(reflectField).Field(0)

	flagSet := flag.NewFlagSet("TestParseStruct_NestedStruct", flag.PanicOnError)
	flagSet.String("parseMeSenpai", "", "")
	err := flagSet.Set("parseMeSenpai", "Test")
	suite.Nil(err)

	err = parseStruct(field, flagSet)
	suite.Nil(err)
	suite.Equal("Test", toBeParsed.SubStruct.LeaveStruct.ParseMeSenpai)
}

//
// writeConfig
//

func (suite *LoadConfigUnitSuite) TestWriteConfig_ValidYaml() {
	var buf []byte
	writer := bytes.NewBuffer(buf)
	config := Config{
		LogLevel: 4,
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
				Host: "127.0.0.1",
				Port: 1,
			},
			IMAP: struct {
				Host string `flag:"imapHost"`
				Port int    `flag:"imapPort"`
			}{
				Host: "0.0.0.0",
				Port: 2,
			},
			HTTP: struct {
				Host string `flag:"httpHost"`
				Port int    `flag:"httpPort"`
			}{
				Host: "0.0.0.0",
				Port: 3,
			},
		},
		DisableIMAP: true,
		DisableSMTP: true,
		DisableHTTP: true,
	}
	err := writeConfig(config, writer)
	suite.Nil(err)
	var resultingConfig Config
	resultingConfig, err = parseConfig(resultingConfig, writer)
	suite.Nil(err)
	suite.Equal(config, resultingConfig)
}

func TestLoad(t *testing.T) {
	suite.Run(t, new(LoadConfigUnitSuite))
}
