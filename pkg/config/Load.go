package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilename = "mailpie"
	envPrefix             = "MAILPIE"
)

func Load(execute func()) *cobra.Command {
	configuration = Config{}

	rootCmd := &cobra.Command{
		Use:   "mailpie",
		Short: "Mailpie the ",
		Long:  `Demonstrate how to get cobra flags to bind to viper properly`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			execute()
		},
	}
	rootCmd.Flags().IntVarP(&configuration.logLevel, "loglevel", "l", int(logrus.WarnLevel), "Possible log leves are:\n0 - Panic\n1 - Fatal\n2 - Error\n3 - Warn\n4 - Info\n5 - Debug\n6 - Trace")
	rootCmd.Flags().StringVarP(&configuration.NetworkConfigs.IMAPHost, "imaphost", "i", "0.0.0.0", "IMAP-host which Mailpie is listening to - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	rootCmd.Flags().StringVarP(&configuration.NetworkConfigs.SMTPHost, "smtphost", "s", "0.0.0.0", "SMTP-host which Mailpie is listening to - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	rootCmd.Flags().StringVarP(&configuration.NetworkConfigs.HTTPHost, "httphost", "h", "0.0.0.0", "HTTP-host where Mailpie serves ths SPA - Use 127.0.0.1 for local access & 0.0.0.0 for network access")
	rootCmd.Flags().IntVarP(&configuration.NetworkConfigs.IMAPPort, "imapport", "I", 1143, "IMAP-port where Mailpie is listening")
	rootCmd.Flags().IntVarP(&configuration.NetworkConfigs.SMTPPort, "smtpport", "S", 1025, "SMTP-port where Mailpie is listening")
	rootCmd.Flags().IntVarP(&configuration.NetworkConfigs.HTTPPort, "httpport", "H", 8000, "HTTP-port where Mailpie serves ths SPA")
	rootCmd.Flags().BoolVarP(&configuration.EnableIMAP, "enableimap", "", true, "Enables the IMAP handler")
	rootCmd.Flags().BoolVarP(&configuration.EnableSMTP, "enablesmtp", "", true, "Enables the SMTP handler")
	rootCmd.Flags().BoolVarP(&configuration.EnableHTTP, "enablehttp", "", true, "Enables the SPA")
	rootCmd.PersistentFlags().BoolP("help", "", false, "helping myself with the command")
	return rootCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(defaultConfigFilename)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	v.AddConfigPath(fmt.Sprintf("%s%s", homeDir, "/.config"))
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	return bindFlags(cmd, v)
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) error {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			err := v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
			if err != nil {
				logrus.WithError(err).Error("Unable to bind env")
			}
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				logrus.WithError(err).Error("Unable to set flag value")
			}
		}
		if f.Name != "help" {
			v.Set(f.Name, f.Value)
		}
	})
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configFile := fmt.Sprintf("%s%s", homeDir, "/.config/mailpie.yml")
	v.SetConfigType("yaml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return v.SafeWriteConfigAs(configFile)
	}
	return nil
}
