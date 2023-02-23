package transcripts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/allokate-ai/events/app/pkg/client"
	"github.com/spf13/cobra"
)

type EarningsCallTranscript struct {
	Date  time.Time
	Title string
	Url   string
}

type DateRange struct {
	From time.Time
	To   time.Time
}

type PagedEarningsCallTranscripts struct {
	Count       int
	TotalPages  int
	Transcripts []EarningsCallTranscript
	Size        int
	Page        int
	DateRange   DateRange
}

func EarningsCallTranscripts(from time.Time, to time.Time, size int, page int, cookies []*http.Cookie) (PagedEarningsCallTranscripts, error) {

	uri := fmt.Sprintf("https://seekingalpha.com/api/v3/articles?filter[category]=earnings::earnings-call-transcripts&filter[since]=%d&filter[until]=%d&include=author,primaryTickers,secondaryTickers&isMounting=true&page[size]=%d&page[number]=%d", from.UTC().Unix(), to.UTC().Unix(), size, page)

	// Craft the request for the page.
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "seekingalpha.com")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15",
	)
	req.Header.Set("Accept-Language", "en-us,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Make the request.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return PagedEarningsCallTranscripts{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PagedEarningsCallTranscripts{}, fmt.Errorf("received bad status code %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return PagedEarningsCallTranscripts{}, err
	}

	// Define the structure of the response from the SeekingAlpha API endpoint.
	var body struct {
		Data []struct {
			Attributes struct {
				PublishOn string
				Title     string
			}
			Links struct {
				Self string
			}
		}
		Meta struct {
			Page struct {
				Total           int
				TotalPages      int
				MinMaxPublishOn struct {
					Min int
					Max int
				}
			}
		}
	}

	if err := json.Unmarshal(content, &body); err != nil {
		return PagedEarningsCallTranscripts{}, err
	}

	var transcripts []EarningsCallTranscript
	for _, transcript := range body.Data {
		timestamp, err := time.Parse(time.RFC3339, transcript.Attributes.PublishOn)
		if err != nil {
			return PagedEarningsCallTranscripts{}, err
		}

		u := url.URL{
			Scheme: "https",
			Host:   "seekingalpha.com",
			Path:   transcript.Links.Self,
		}

		transcripts = append(transcripts, EarningsCallTranscript{
			Date:  timestamp,
			Title: transcript.Attributes.Title,
			Url:   u.String(),
		})
	}

	return PagedEarningsCallTranscripts{
		Count:       body.Meta.Page.Total,
		TotalPages:  body.Meta.Page.TotalPages,
		Transcripts: transcripts,
		Size:        size,
		Page:        page,
		DateRange: DateRange{
			From: time.Unix(int64(body.Meta.Page.MinMaxPublishOn.Min), 0),
			To:   time.Unix(int64(body.Meta.Page.MinMaxPublishOn.Max), 0),
		},
	}, err
}

func EarningsCallTranscriptsFromDate(date time.Time, cookies []*http.Cookie) ([]EarningsCallTranscript, error) {

	// Calculate the date range from the start of yesterday (UTC) to midnight today.
	from := date.UTC().Truncate(24 * time.Hour)
	to := from.Add(24 * time.Hour)

	// Page starts at 1, not 0 according to their API.
	page := 1

	var transcripts []EarningsCallTranscript
	for {
		data, err := EarningsCallTranscripts(from, to, 250, page, cookies)
		if err != nil {
			return transcripts, err
		}

		transcripts = append(transcripts, data.Transcripts...)

		// Break if we fetched the last page.
		if data.Page >= data.TotalPages {
			break
		}

		// Get the next page.
		page = data.Page + 1
	}

	return transcripts, nil
}

// rootCmd represents the base command when called without any subcommands
var Cmd = &cobra.Command{
	Use: "transcripts",
	// Short: "A tool various data feeds",
	// Long:  `This is a tool used to scrape various RSS news feeds, API, and other sources for data used by Allokate.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		yesterday := time.Now().Add(-24 * time.Hour)

		cookies, err := GetSeekingAlphaCookies()
		if err != nil {
			log.Fatal(fmt.Errorf("failed to get cookies: %w", err))
		}

		transcripts, err := EarningsCallTranscriptsFromDate(yesterday, cookies)
		if err != nil {
			log.Fatal(err)
		}

		for _, transcript := range transcripts {

			// Create the event
			article := client.ArticlePublished{
				Source:   "https://seekingalpha.com/api/v3/articles",
				SiteName: "Seeking Alpha",
				Byline:   "Seeking Alpha",
				Title:    transcript.Title,
				Url:      transcript.Url,
				Date:     transcript.Date,
			}

			// Send it!!
			if _, err := client.Default().EmitArticlePublishedEvent("feeds.seekingalpha.transcripts", article); err != nil {
				log.Fatal(err)
			} else {
				log.Printf("Transcript '%s' published on %s (%s)", article.Title, article.Date.Local(), article.Date.Local())
			}

		}

		fmt.Println("Fetched", len(transcripts), "transcripts for", yesterday.Format("2006-01-02"))
	},
}

func GetSeekingAlphaCookies() ([]*http.Cookie, error) {

	// Craft the request for the page.
	req, _ := http.NewRequest("GET", "https://seekingalpha.com/earnings/earnings-call-transcripts", nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "seekingalpha.com")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15",
	)
	req.Header.Set("Accept-Language", "en-us,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://seekingalpha.com/earnings/earnings-call-transcripts")

	// Make the request.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Cookies(), nil

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
