package rss

import (
	"log"
	"time"

	"github.com/allokate-ai/feeds/app/internal/event"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "rss",
	Short: "Scrape news article from investing.com's RSS feeds",
	Run: func(cmd *cobra.Command, args []string) {

		// Define the source url for the feed.
		url := "https://www.investing.com/rss/news_25.rss"

		// Parse the feed.
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(url)

		// Iterate over each item in the feed and publish the article information to the event system.
		for _, item := range feed.Items {
			loc, err := time.LoadLocation("UTC")
			timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", item.Published, loc)
			if err != nil {
				log.Fatal(err)
			}

			author := ""
			if len(item.Authors) > 0 {
				author = item.Author.Email
			}

			// Create the event
			article := event.ArticlePublished{
				Source:   url,
				SiteName: "Investing.com",
				Byline:   author,
				Title:    item.Title,
				Url:      item.Link,
				Date:     timestamp,
				Tags:     []string{},
			}

			// Send it!!
			if _, err := event.EmitArticlePublishedEvent(article); err != nil {
				log.Fatal(err)
			} else {
				log.Printf("Article '%s' published on %s (%s)", article.Title, article.Date.Local(), article.Date.Local())
			}
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.feeds.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
