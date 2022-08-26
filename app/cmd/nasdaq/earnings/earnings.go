package earnings

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/allokate-ai/feeds/app/internal/event"
	"github.com/spf13/cobra"
)

type Earnings struct {
	Date   time.Time
	Ticker string
}

func EarningsOnDate(date time.Time) ([]Earnings, error) {

	date = date.UTC().Truncate(24 * time.Hour)

	uri := fmt.Sprintf("https://api.nasdaq.com/api/calendar/earnings?date=%s", date.Format("2006-01-02"))

	// Craft the request for the page.
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set(
		"Accept",
		"application/json",
	)
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15",
	)
	req.Header.Set("Accept-Language", "en-us")

	// Make the request.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return []Earnings{}, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Earnings{}, err
	}

	// Define the structure of the response from the SeekingAlpha API endpoint.
	var body struct {
		Data struct {
			Rows []struct {
				Symbol string
			}
		}
	}

	if err := json.Unmarshal(content, &body); err != nil {
		return []Earnings{}, err
	}

	earnings := []Earnings{}
	for _, item := range body.Data.Rows {

		earnings = append(earnings, Earnings{
			Date:   date,
			Ticker: item.Symbol,
		})
	}

	return earnings, err
}

var Cmd = &cobra.Command{
	Use:   "earnings",
	Short: "Scrape company earnings reporting data",
	Run: func(cmd *cobra.Command, args []string) {

		today := time.Now().UTC().Truncate(24 * time.Hour)

		earnings, err := EarningsOnDate(today)
		if err != nil {
			log.Fatal(err)
		}

		for _, item := range earnings {
			// Create the event
			earnings := event.Earnings{
				Date:   item.Date,
				Ticker: item.Ticker,
			}

			// Send it!!
			if _, err := event.EmitEarningsEvent(earnings); err != nil {
				log.Fatal(err)
			} else {
				log.Printf("%s reporting earnings on %s", earnings.Ticker, earnings.Date)
			}
		}

		fmt.Println(len(earnings), "reporting earnings for", today.Format("2006-01-02"))
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
