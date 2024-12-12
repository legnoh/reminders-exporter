package cmd

import (
	"github.com/legnoh/reminders-exporter/pkg/reminder"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Cmd struct {
	Path string
	Args []string
}

type Config struct {
	Port    string
	Filters []reminder.Filter
}

var (
	cfgFile string
	conf    Config
	confDir string
	version string
	debug   bool
	log     *logrus.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "reminders-exporter",
	Version: version,
	Short:   "Reminder app data exporter for prometheus",
	Long: `This tool provides daemon service for Apple Reminder app data exporter for Prometheus.

If you have any questions, please visit github site.
https://github.com/legnoh/reminders-exporter`,
}

func Execute() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if rootCmd.Execute() != nil {
		log.Fatal("Root execute is failed... exit")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debug log")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	debug, err := rootCmd.PersistentFlags().GetBool("debug")
	if err != nil {
		log.Fatal(err)
	}
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}
}
