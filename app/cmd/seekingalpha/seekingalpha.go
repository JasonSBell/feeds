package seekingalpha

import (
	"github.com/allokate-ai/feeds/app/cmd/seekingalpha/rss"
	"github.com/allokate-ai/feeds/app/cmd/seekingalpha/transcripts"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "seekingalpha",
	Short: "Tooling for scraping data from seekingalpha.com",
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.feeds.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	Cmd.AddCommand(transcripts.Cmd)
	Cmd.AddCommand(rss.Cmd)
}
