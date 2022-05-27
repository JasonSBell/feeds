package main

import (
	"log"
	"time"

	"github.com/allokate-ai/feeds/app/internal/event"
	"github.com/mmcdole/gofeed"
)

func main() {

	// Define the source url(s) for the feed.
	urls := []string{
		"https://feeds.a.dj.com/rss/WSJcomUSBusiness.xml",
		"https://feeds.a.dj.com/rss/RSSMarketsMain.xml",
	}

	for _, url := range urls {

		// Parse the feed.
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(url)

		// Iterate over each item in the feed and publish the article information to the event system.
		for _, item := range feed.Items {

			timestamp, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.Published)
			if err != nil {
				log.Fatal(err)
			}

			author := ""
			if len(item.Authors) > 0 {
				author = item.Author.Email
			}

			// Create the event
			article := event.ArticlePublished{
				Source: url,
				Byline: author,
				Title:  item.Title,
				Url:    item.Link,
				Date:   timestamp,
				Tags:   []string{},
			}

			// Send it!!
			if _, err := event.EmitArticlePublishedEvent(article); err != nil {
				log.Fatal(err)
			} else {
				log.Printf("Article '%s' published on %s (%s)", article.Title, article.Date.Local(), article.Date.Local())
			}
		}
	}

}
