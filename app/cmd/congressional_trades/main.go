package main

import (
	"fmt"
	"log"

	"github.com/allokate-ai/feeds/app/internal/event"
)

func main() {

	senateTrades, err := AllSenateTrades()
	if err != nil {
		log.Fatal(err)
	}

	for _, trade := range senateTrades {

		// Skip trades that do not have a transaction date.
		if trade.TransactionDate == nil {
			continue
		}

		// Skip trades that do not have a proper name.
		if trade.Name == "" {
			continue
		}

		e := event.CongressionalTrade{
			Body:            "Senate",
			TransactionDate: *trade.TransactionDate,
			DisclosureDate:  trade.DisclosureDate,
			Url:             trade.Url,
			Name:            trade.Name,
			Owner:           trade.Owner,
			Ticker:          trade.Ticker,
			AssetType:       trade.AssetType,
			Type:            string(trade.Type),
			Comment:         trade.Comment,
			Amount:          trade.Amount,
		}

		// Send it!!
		if _, err := event.EmitCongressionalTradeEvent(e); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Congressional trade event published for %s member %s trading %s on %s", e.Body, trade.Name, trade.Ticker, trade.TransactionDate.Local())
		}

	}

	fmt.Println("Fetched", len(senateTrades), "senate transactions")

	houseTrades, err := AllHouseTrades()
	if err != nil {
		log.Fatal(err)
	}

	for _, trade := range houseTrades {

		// Skip trades that do not have a transaction date.
		if trade.TransactionDate == nil {
			continue
		}

		// Skip trades that do not have a proper name.
		if trade.Name == "" {
			continue
		}

		e := event.CongressionalTrade{
			Body:            "House of Representatives",
			TransactionDate: *trade.TransactionDate,
			DisclosureDate:  trade.DisclosureDate,
			Url:             trade.Url,
			Name:            trade.Name,
			Owner:           trade.Owner,
			Ticker:          trade.Ticker,
			Type:            string(trade.Type),
			Amount:          trade.Amount,
		}

		// Send it!!
		if _, err := event.EmitCongressionalTradeEvent(e); err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Congressional trade event published for %s member %s trading %s on %s", e.Body, trade.Name, trade.Ticker, trade.TransactionDate.Local())
		}

	}

	fmt.Println("Fetched", len(houseTrades), "house transactions")
}
