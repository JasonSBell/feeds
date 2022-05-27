package event

import (
	"encoding/json"
	"time"

	"github.com/allokate-ai/events/pkg/client"
	"github.com/allokate-ai/events/pkg/events"
	"github.com/google/uuid"
)

var c *client.Client

// Define the JSON body structure for publishing an article (article.created).
type Article struct {
	Source string    `json:"source"`
	Byline string    `json:"byline"`
	Title  string    `json:"title"`
	Url    string    `json:"url"`
	Date   time.Time `json:"date"`
	Tags   []string  `json:"tags"`
}

// Define the JSON body structure for publishing an article (congressional_trade).
type CongressionalTrade struct {
	Body            string     `json:"body"`
	TransactionDate time.Time  `json:"transactionDate"`
	DisclosureDate  *time.Time `json:"disclosureDate"`
	Url             string     `json:"url"`
	Name            string     `json:"name"`
	Owner           string     `json:"owner"`
	Ticker          string     `json:"ticker"`
	AssetType       string     `json:"assetType"`
	Type            string     `json:"type"`
	Comment         string     `json:"comment"`
	Amount          string     `json:"amount"`
}

func Client() *client.Client {
	if c == nil {
		// Declare a client that will be used to publish new articles.
		cli, err := client.NewClient("http://192.168.4.64:8092", nil)
		if err != nil {
			panic(err)
		}
		c = cli
	}

	return c

}

func EmitArticlePublishedEvent(article Article) (events.GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(article)
	if err != nil {
		return events.GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := events.GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "news",
		Source:    "feeds",
		Body:      data,
	}

	return Client().Publish(e)
}

func EmitCongressionalTradeEvent(trade CongressionalTrade) (events.GenericEvent, error) {
	// Serialize the body to a JSON string.
	data, err := json.Marshal(trade)
	if err != nil {
		return events.GenericEvent{}, err
	}

	// Define the event for publishing articles.
	e := events.GenericEvent{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Name:      "congressional_trade",
		Source:    "feeds",
		Body:      data,
	}

	return Client().Publish(e)
}
