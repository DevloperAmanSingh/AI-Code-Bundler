package cmd

import (
	"log"
	"os"

	"github.com/AkhilSharma90/AI-Code-Bundler/internal/review"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Let an AI review your codefuse-project.txt",
	Long:  `Let an AI review the codefuse-project.txt you generated with the bundle command.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey := viper.GetString("codefuse")
		if apiKey == "" {
			log.Fatal(`Api key is required for review.`)
		}
		dat, err := os.ReadFile("codefuse-project.txt")
		if err != nil {
			log.Fatal("Could not find codefuse-project.txt. Did you forget to run the \"codefuse bundle\" command?")
		}
		review.Review(string(dat), apiKey)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.Flags().String("codefuse_api_key", "", "Your Code AI Review API key ")
	err := viper.BindPFlag("codefuse_api_key", reviewCmd.Flags().Lookup("codefuse_api_key"))
	if err != nil {
		log.Fatal(err)
	}
}
