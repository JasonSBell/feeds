package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

func main() {

	// Define the source url for the feed.
	// url := "https://www.sec.gov/news/pressreleases.rss"
	// url := "https://www.sec.gov/rss/litigation/litreleases.xml"
	// "https://www.nasdaq.com/feed/rssoutbound?symbol=crwd",
	// "https://seekingalpha.com/api/sa/combined/CRWD.xml",
	// "https://seekingalpha.com/sector/transcripts.xml",
	// "https://investor.docusign.com/rss/PressRelease.aspx?LanguageId=1&CategoryWorkflowId=1cb807d2-208f-4bc3-9133-6a9ad45ac3b0&tags=",
	// "https://investor.docusign.com/rss/event.aspx",
	// "https://ir.crowdstrike.com/rss/events.xml",
	// "http://www.nasdaqtrader.com/rss.aspx?feed=currentheadlines&categorylist=2,6,7",
	url := "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type=&company=&dateb=&owner=include&start=0&count=100&output=atom"

	// Parse the feed.
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)

	// Iterate over each item in the feed and publish the article information to the event system.
	for _, item := range feed.Items {
		fmt.Println(item.Published, item.Title, item.Link)

	}

}
