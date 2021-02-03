package config

import (
	"flag"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/user"
	"reflect"
	"strconv"
)

func Load() {
	initFlags()
	flag.Parse()
	createConfig := false
	file, err := os.Open(flag.Lookup("config").Value.String())
	if err != nil {
		createConfig = true
	} else {
		configuration = parseConfig(configuration, file)
		defer file.Close()
	}
	configuration = combineConfigAndFlags(configuration)
	configuration.LogrusLevel = logrus.Level(configuration.LogLevel)
	if createConfig {
		writeConfig(configuration)
	}
}

func initFlags() {
	flag.Int("logLevel", int(logrus.WarnLevel), "Possible log leves are:\n0 - Panic\n1 - Fatal\n2 - Error\n3 - Warn\n4 - Info\n5 - Debug\n6 - Trace")
	flag.String("imapHost", "0.0.0.0", "IMAP-host which Mailpie is listening to - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	flag.String("smtpHost", "0.0.0.0", "SMTP-host which Mailpie is listening to - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	flag.String("httpHost", "0.0.0.0", "HTTP-host where Mailpie serves ths SPA - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	flag.Int("imapPort", 1143, "IMAP-port where Mailpie is listening")
	flag.Int("smtpPort", 1025, "SMTP-port where Mailpie is listening")
	flag.Int("httpPort", 8000, "HTTP-port where Mailpie serves ths SPA")
	flag.Bool("enableImap", true, "Enables the IMAP handler")
	flag.Bool("enableSmtp", true, "Enables the SMTP handler")
	flag.Bool("enableHttp", true, "Enables the SPA")
	usr, _ := user.Current()
	dir := usr.HomeDir
	flag.String("config", dir+"/.config/mailpie.yml", "sets the config file path. If file not exits, MailPie will create one with default values.")
}

func parseConfig(config Config, file io.Reader) Config {
	readFile, err := io.ReadAll(file)
	if err != nil {
		logrus.WithError(err).Info("no config found, will be created")
		return Config{}
	}
	err = yaml.Unmarshal(readFile, &config)
	if err != nil {
		logrus.WithError(err).Fatal("couldn't unmarshal config file")
		return Config{}
	}
	return config
}

func combineConfigAndFlags(config Config) Config {
	reflectConfig := reflect.ValueOf(&config)
	for i := 0; i < reflect.Indirect(reflectConfig).NumField(); i++ {
		field := reflect.Indirect(reflectConfig).Field(i)
		if field.Kind() == reflect.Struct {
			parseStruct(field)
			continue
		}
		fieldType := reflect.Indirect(reflectConfig).Type().Field(i)
		parseField(field, &fieldType)
	}
	return config
}

func parseStruct(structToParse reflect.Value) {
	for i := 0; i < structToParse.NumField(); i++ {
		field := structToParse.Field(i)
		if field.Kind() == reflect.Struct {
			parseStruct(field)
			continue
		}
		fieldType := structToParse.Type().Field(i)
		parseField(field, &fieldType)
	}
}

func parseField(fieldToParse reflect.Value, fieldType *reflect.StructField) {
	tag := fieldType.Tag
	flagName, exists := tag.Lookup("flag")
	if !exists {
		return
	}
	if fieldToParse.IsZero() {
		overrideValue(&fieldToParse, flag.Lookup(flagName))
		return
	}
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			overrideValue(&fieldToParse, f)
		}
	})

}

func overrideValue(field *reflect.Value, f *flag.Flag) {
	flagValue := f.Value
	flagString := flagValue.String()
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(flagString)

	case reflect.Int:
		atoi, err := strconv.Atoi(flagString)
		if err != nil {
			logrus.WithError(err).WithField("flag", f.Name).Error("unable to parse flag to int")
		}
		field.SetInt(int64(atoi))

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(flagString)
		if err != nil {
			logrus.WithError(err).WithField("flag", f.Name).Error("unable to parse flag to bool")
		}
		field.SetBool(boolValue)
	}
}

func writeConfig(config Config) {
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		logrus.WithError(err).Error("unable to marshal config")
		return
	}
	file, err := os.OpenFile(flag.Lookup("config").Value.String(), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.WithError(err).Error("unable to create config file")
		return
	}
	defer file.Close()
	_, err = file.Write(configBytes)
	if err != nil {
		logrus.WithError(err).Error("unable to write to config file")
		return
	}
}
