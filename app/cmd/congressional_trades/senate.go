package main

// Download Individual Disclosures
// For downloading individual pieces of information, there are two primary folders where information is store. These two folders are data/ and aggregate/.

// All files can be accessed via the main data URL at https://senate-stock-watcher-data.s3-us-west-2.amazonaws.com.

// You can get the list of available disclosures by fetching the endpoint of https://senate-stock-watcher-data.s3-us-west-2.amazonaws.com/aggregate/filemap.xml
// Why is there no traditional API? That is because this service is statically hosted and cannot support POST requests or do any processing besides serving a file. This is the most cost effective option to run this service - as I run it at my own cost.
// Note:You can just skip these two steps and grab https://senate-stock-watcher-data.s3-us-west-2.amazonaws.com/aggregate/all_transactions.json. This file will only include the transactions that have been discovered. The below method just indiciates that a disclosure occured that day, but not that the transaction data is available, as it may not have been transcribed yet.

// Here is an example in Javascript of fetching the filemap.xml file.

//   fetch('https://senate-stock-watcher-data.s3-us-west-2.amazonaws.com/aggregate/filemap.xml')
//     .then((response) => response.text())
//     .then((response) => {
//       const parser = new DOMParser()
//       const xml = parser.parseFromString(response, 'text/xml')
//       const results = [].slice.call( xml.getElementsByTagName('Key') ).filter((key) => key.textContent.includes('.json'))
//       const files = results.map(file => file.textContent.split('/')[1])
//     })
//     .catch((response) => {
//       console.log(response)
//     })

//     // results ['data/transaction_report_for_01_01_2021.json',
//     'data/transaction_report_for_01_02_2021.json',
//     ...,
//     'data/transaction_report_for_12_31_2021.json']

// You can then get the JSON for a single days disclosure by easily fetching the file found in the Key tag.

//   fetch('https://senate-stock-watcher-data.s3-us-west-2.amazonaws.com/data/<your_transaction_report_for_x_x_xxxx>.json')
//     .then((response) => response.json())
//     .then((response) => {
//       // here is your json response for that day. Each object is a distinct document filed that day.
//       // Each object has a 'transactions' key that will contain all the trade information on that day.
//     })
//     .catch((response) => {
//       console.log(response)
//     })

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type SenateTrade struct {
	TransactionDate *time.Time
	DisclosureDate  *time.Time
	Url             string
	Name            string
	Owner           string
	Ticker          string
	AssetType       string
	Type            TransactionType
	Comment         string
	Amount          string
}

func AllSenateTrades() ([]SenateTrade, error) {

	uri := "https://senate-stock-watcher-data.s3-us-west-2.amazonaws.com/aggregate/all_transactions.json"

	// Craft the request for the page.
	req, _ := http.NewRequest("GET", uri, nil)

	// Make the request.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return []SenateTrade{}, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return []SenateTrade{}, err
	}

	// Define the structure of the response from the API endpoint.
	var body []struct {
		Senator         string `json:"senator"`
		TransactionDate string `json:"transaction_date"`
		DisclosureDate  string `json:"disclosure_date"`
		PtrLink         string `json:"ptr_link"`
		Ticker          string `json:"ticker"`
		Owner           string `json:"owner"`
		Amount          string `json:"amount"`
		Type            string `json:"type"`
		AssetType       string `json:"asset_type"`
		Comment         string `json:"comment"`
	}

	if err := json.Unmarshal(content, &body); err != nil {
		return []SenateTrade{}, err
	}

	trades := []SenateTrade{}
	for _, trade := range body {
		trades = append(trades, SenateTrade{
			Name:            strings.TrimSpace(trade.Senator),
			TransactionDate: ParseDate(trade.TransactionDate),
			DisclosureDate:  ParseDate(trade.DisclosureDate),
			Url:             trade.PtrLink,
			Ticker:          strings.ToUpper(strings.TrimSpace(trade.Ticker)),
			Owner:           strings.Trim(trade.Owner, "- "),
			Amount:          trade.Amount,
			AssetType:       trade.AssetType,
			Comment:         trade.Comment,
			Type:            StringToTransactionType(trade.Type),
		})
	}

	return trades, err
}

func AllSenateTradesOnDate(date time.Time) ([]SenateTrade, error) {
	startOfDay := date.UTC().Truncate(24 * time.Hour).Add(-time.Microsecond)
	endOfDay := startOfDay.Add(24*time.Hour + 2*time.Microsecond)

	allTrades, err := AllSenateTrades()
	if err != nil {
		return allTrades, err
	}

	trades := []SenateTrade{}
	for _, trade := range allTrades {
		if trade.TransactionDate != nil && trade.TransactionDate.After(startOfDay) && trade.TransactionDate.Before(endOfDay) {
			trades = append(trades, trade)
		}
	}

	return trades, err
}

func AllSenateTradesSince(date time.Time) ([]SenateTrade, error) {
	allTrades, err := AllSenateTrades()
	if err != nil {
		return allTrades, err
	}

	trades := []SenateTrade{}
	for _, trade := range allTrades {
		if trade.TransactionDate != nil && trade.TransactionDate.After(date) {
			trades = append(trades, trade)
		}
	}

	return trades, err
}
