// Package cmd provides the root command for the codefuse tool.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "codefuse",
	Version: Version,
	Short:   "Initialize",
	Long: `Allows you to bundle your codebase and let it be reviewed by an AI.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// otherwise the completion command will be available
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Search the current directory for a config file
	viper.SetConfigType("yaml")
	viper.SetConfigName(".codefuse-config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}