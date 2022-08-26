package nasdaq

import (
	"github.com/allokate-ai/feeds/app/cmd/nasdaq/dividends"
	"github.com/allokate-ai/feeds/app/cmd/nasdaq/earnings"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "nasdaq",
	Short: "Tooling for scraping data from nasdaq.com",
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.feeds.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	Cmd.AddCommand(earnings.Cmd)
	Cmd.AddCommand(dividends.Cmd)
}
