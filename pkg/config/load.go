package config

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"os/user"
	"reflect"
	"strconv"
)

func Load(flags *flag.FlagSet, arguments []string) error {
	initFlags(flags)
	err := flags.Parse(arguments)
	if err != nil {
		return errors.Wrap(err, "Unable to parse flags")
	}
	createConfig := false
	configPath := flags.Lookup("config").Value.String()
	file, err := os.Open(configPath)
	if err != nil {
		createConfig = true
	} else {
		configuration, err = parseConfig(configuration, file)
		if err != nil {
			//We only log this error and don't return so that we can still run MailPie with the commandline flags
			logrus.WithError(err).Error("Unable to parse config file")
		}
		defer func() {
			err := file.Close()
			if err != nil {
				logrus.WithError(err).Error("Unable to close config file after read.")
			}
		}()
	}
	configuration, err = combineConfigAndFlags(configuration, flags)
	if err != nil {
		return errors.Wrap(err, "Error combining config and flags")
	}
	configuration.LogrusLevel = logrus.Level(configuration.LogLevel)
	if createConfig {
		file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return errors.Wrap(err, "couldn't create file to write config to")
		}
		defer func() {
			err := file.Close()
			if err != nil {
				logrus.WithError(err).Error("Unable to close config file after write.")
			}
		}()
		return errors.Wrap(writeConfig(configuration, file), "unable to create config file")
	}
	return nil
}

func initFlags(flags *flag.FlagSet) {
	flags.Int("logLevel", int(logrus.WarnLevel), "Possible log leves are:\n0 - Panic\n1 - Fatal\n2 - Error\n3 - Warn\n4 - Info\n5 - Debug\n6 - Trace")
	flags.String("imapHost", "0.0.0.0", "IMAP-host which Mailpie is listening to - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	flags.String("smtpHost", "0.0.0.0", "SMTP-host which Mailpie is listening to - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	flags.String("httpHost", "0.0.0.0", "HTTP-host where Mailpie serves ths SPA - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	flags.Int("imapPort", 1143, "IMAP-port where Mailpie is listening")
	flags.Int("smtpPort", 1025, "SMTP-port where Mailpie is listening")
	flags.Int("httpPort", 8000, "HTTP-port where Mailpie serves ths SPA")
	flags.Bool("disableImap", false, "Disable the IMAP handler")
	flags.Bool("disableSmtp", false, "Disable the SMTP handler")
	flags.Bool("disableHttp", false, "Disable the SPA")
	usr, _ := user.Current()
	dir := usr.HomeDir
	flags.String("config", dir+"/.config/mailpie.yml", "sets the config file path. If file not exits, MailPie will create one with default values.")
}

func parseConfig(config Config, file io.Reader) (Config, error) {
	readFile, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(readFile, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func combineConfigAndFlags(config Config, flagSet *flag.FlagSet) (Config, error) {
	reflectConfig := reflect.ValueOf(&config)
	for i := 0; i < reflect.Indirect(reflectConfig).NumField(); i++ {
		field := reflect.Indirect(reflectConfig).Field(i)
		if field.Kind() == reflect.Struct {
			err := parseStruct(field, flagSet)
			if err != nil {
				return Config{}, errors.Wrap(err, "Error parsing nested struct")
			}
			continue
		}
		fieldType := reflect.Indirect(reflectConfig).Type().Field(i)
		err := parseField(field, &fieldType, flagSet)
		if err != nil {
			return Config{}, errors.Wrap(err, "Error parsing struct field")
		}
	}
	return config, nil
}

func parseStruct(structToParse reflect.Value, flagSet *flag.FlagSet) error {
	for i := 0; i < structToParse.NumField(); i++ {
		field := structToParse.Field(i)
		if field.Kind() == reflect.Struct {
			err := parseStruct(field, flagSet)
			if err != nil {
				return errors.Wrap(err, "Error parsing nested struct")
			}
			continue
		}
		fieldType := structToParse.Type().Field(i)
		err := parseField(field, &fieldType, flagSet)
		if err != nil {
			return errors.Wrap(err, "Error parsing struct field")
		}
	}
	return nil
}

func parseField(fieldToParse reflect.Value, fieldType *reflect.StructField, flagSet *flag.FlagSet) error {
	tag := fieldType.Tag
	flagName, exists := tag.Lookup("flag")
	if !exists {
		return nil
	}
	if fieldToParse.IsZero() {
		return errors.Wrap(overrideValue(&fieldToParse, flagSet.Lookup(flagName)), "Unable to override unset config value")
	}
	var visitError error
	flagSet.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			visitError = overrideValue(&fieldToParse, f)
		}
	})
	return errors.Wrap(visitError, "Unable to override config value with set flag")
}

func overrideValue(field *reflect.Value, f *flag.Flag) error {
	flagValue := f.Value
	flagString := flagValue.String()
	switch field.Type().Kind() {
	case reflect.String:
		field.SetString(flagString)

	case reflect.Int:
		atoi, err := strconv.Atoi(flagString)
		if err != nil {
			return errors.Wrapf(err, "unable to parse flag '%s' as int", f.Name)
		}
		field.SetInt(int64(atoi))

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(flagString)
		if err != nil {
			return errors.Wrapf(err, "unable to parse flag '%s' as bool", f.Name)
		}
		field.SetBool(boolValue)
	default:
		return errors.New("unsupported config type: " + field.Type().Kind().String())
	}
	return nil
}

func writeConfig(config Config, writer io.Writer) error {
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal config")
	}
	_, err = writer.Write(configBytes)
	return errors.Wrap(err, "couldn't write to config file")
}
