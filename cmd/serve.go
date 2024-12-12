package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/creasty/defaults"
	"github.com/legnoh/reminders-exporter/pkg/collector"
	"github.com/legnoh/reminders-exporter/pkg/reminder"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start reminders-exporter daemon",
	Long: `Start reminders-exporter application.

This command requires config file.
Before start this command, prepare config file with "reminders-exporter init" command.`,
	PreRun: preStartServer,
	Run:    startServer,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	confDir = home + "/.reminders-exporter"

	serveCmd.Flags().StringVarP(&cfgFile, "config", "c", confDir+"/config.yml", "config file path")
}

func preStartServer(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}
	if err := defaults.Set(&conf); err != nil {
		log.Fatal(err)
	}
	if !regexp.MustCompile(`^[0-9]{4}$`).MatchString(conf.Port) {
		log.Fatalf("Specified Port(%s) is invalid format. Please fix to 4-digit code(e.g. 8888)", conf.Port)
	}

	if _, err := os.Stat(reminder.ReminderCLIPath); err != nil {
		log.Fatalf("reminders-cli is not installed this device. Please install it with homebrew! https://github.com/keith/reminders-cli/blob/main/README.md#with-homebrew")
	}
}

func startServer(cmd *cobra.Command, args []string) {
	reg := prometheus.NewRegistry()
	collector := collector.ReminderCollector{Filters: conf.Filters}
	reg.MustRegister(collector)

	server := &http.Server{Addr: ":" + conf.Port}
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	go func() {
		log.Info("Starting Reminder Exporter...")
		log.Debugf("Config File: %s", cfgFile)
		log.Debugf("       Port: %s", conf.Port)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("reminders-exporter server error: %v", err)
		}
		log.Info("Stopped serving new connections....")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Info("reminders-exporter's Graceful shutdown complete.")
}
