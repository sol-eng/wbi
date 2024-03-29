package cmd

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type settings struct {
	// logrus log level
	loglevel string
}

type rootCmd struct {
	cmd *cobra.Command
	cfg *settings
}

func Execute(version string, args []string) {
	newRootCmd(version).Execute(args)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)
	if err := cmd.cmd.Execute(); err != nil {
		// if get to this point and don't fatally log in the subcommand,
		// the Usage help will be printed before the error,
		// which may or may not be the desired behavior
		log.Fatalf("failed with error: %s", err)
	}
}

func setGlobalSettings(cfg *settings) {
	cfg.loglevel = viper.GetString("loglevel")
	setLogLevel(cfg.loglevel)
	setUpLogger()
}
func newRootCmd(version string) *rootCmd {
	root := &rootCmd{cfg: &settings{}}
	cmd := &cobra.Command{
		Use:   "wbi",
		Short: "workbench installer",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// need to set the config values here as the viper values
			// will not be processed until Execute, so can't
			// set them in the initializer.
			// If persistentPreRun is used elsewhere, should
			// remember to setGlobalSettings in the initializer
			setGlobalSettings(root.cfg)
		},
	}
	cmd.Version = version
	// without this, the default version is like `cmd version <version>` so this
	// will just print the version for simpler parsing
	cmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
	cmd.PersistentFlags().String("loglevel", "info", "log level")
	viper.BindPFlag("loglevel", cmd.PersistentFlags().Lookup("loglevel"))
	cmd.AddCommand(newSetupCmd().cmd)
	cmd.AddCommand(newVerifyCmd().cmd)
	cmd.AddCommand(newConfigCmd().cmd)
	cmd.AddCommand(newInstallCmd().cmd)
	cmd.AddCommand(newScanCmd().cmd)
	cmd.AddCommand(newActivateCmd().cmd)

	root.cmd = cmd
	return root
}

func setUpLogger() error {
	// Setup the logger output
	logFile := "wbi-log-" + time.Now().Format("20060102T150405") + ".log"

	var f *os.File
	var err error

	if f, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		// Using fmt to print to stdout since logger is not ready
		fmt.Println(err)
		return err
	}

	log.SetOutput(f)

	// Setup the logger format
	log.SetFormatter(&log.JSONFormatter{})

	return nil
}
