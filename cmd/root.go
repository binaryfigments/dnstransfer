package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var log = logrus.New()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dnstransfer",
	Short: "This application can check for zone transfers on domains.",
	Long:  `This application can check for zone transfers on domains.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dnstransfer.yaml)")
	rootCmd.PersistentFlags().BoolP("json", "j", false, "JSON output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "all output for debugging")
	rootCmd.PersistentFlags().String("logfile", "", "Output to logfile.")
	rootCmd.PersistentFlags().String("nameserver", "8.8.8.8", "What nameserver to use.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dnstransfer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dnstransfer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
