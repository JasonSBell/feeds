package main

import (
	"os"

	"github.com/allokate-ai/feeds/app/cmd/cnbc"
	"github.com/allokate-ai/feeds/app/cmd/congress"
	"github.com/allokate-ai/feeds/app/cmd/geekwire"
	"github.com/allokate-ai/feeds/app/cmd/investing"
	"github.com/allokate-ai/feeds/app/cmd/marketwatch"
	"github.com/allokate-ai/feeds/app/cmd/nasdaq"
	"github.com/allokate-ai/feeds/app/cmd/seekingalpha"
	"github.com/allokate-ai/feeds/app/cmd/wsj"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "feeds",
	Short: "A tool various data feeds",
	Long:  `This is a tool used to scrape various RSS news feeds, API, and other sources for data used by Allokate.`,
}

func main() {
	// Add top level flags here.
	// Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	Cmd.AddCommand(cnbc.Cmd)
	Cmd.AddCommand(congress.Cmd)
	Cmd.AddCommand(investing.Cmd)
	Cmd.AddCommand(marketwatch.Cmd)
	Cmd.AddCommand(nasdaq.Cmd)
	Cmd.AddCommand(seekingalpha.Cmd)
	Cmd.AddCommand(wsj.Cmd)
	Cmd.AddCommand(geekwire.Cmd)

	// Execute adds all child commands to the root command and sets flags appropriately.
	// This is called by main.main(). It only needs to happen once to the rootCmd.
	err := Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
