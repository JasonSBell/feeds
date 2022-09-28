package marketwatch

import (
	"github.com/allokate-ai/feeds/app/cmd/marketwatch/rss"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "marketwatch",
	Short: "Tooling for scraping data from marketwatch.com",
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.feeds.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	Cmd.AddCommand(rss.Cmd)
}
