package main

// Download Individual Disclosures
// For downloading individual pieces of information, there are two primary folders where information is store. These two folders are data/ and disclosure_docs/.

// All data will be accessed via the main data URL at https://house-stock-watcher-data.s3-us-west-2.amazonaws.com.

// You can get the list of available disclosures by fetching the endpoint of https://house-stock-watcher-data.s3-us-west-2.amazonaws.com/data/filemap.xml
// Why is there no traditional API? That is because this service is statically hosted and cannot support POST requests or do any processing besides serving a file. This is the most cost effective option to run this service - as I run it at my own cost.
// Note:You can just skip these two steps and grab https://house-stock-watcher-data.s3-us-west-2.amazonaws.com/data/all_transactions.json. This file will only include the transactions that have been transcribed. The below method just indiciates that a disclosure occured that day, but not that the transaction data is available, as it may not have been transcribed yet.

// Here is an example in Javascript of fetching the filemap.xml file.

//   fetch('https://house-stock-watcher-data.s3-us-west-2.amazonaws.com/data/filemap.xml')
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

//   fetch('https://house-stock-watcher-data.s3-us-west-2.amazonaws.com/<your_transaction_report_for_x_x_xxxx>.json')
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
	"regexp"
	"strings"
	"time"
)

type HouseTrade struct {
	TransactionDate *time.Time
	DisclosureDate  *time.Time
	Url             string
	Name            string
	Owner           string
	Ticker          string
	Type            TransactionType
	Amount          string
}

func AllHouseTrades() ([]HouseTrade, error) {

	uri := "https://house-stock-watcher-data.s3-us-west-2.amazonaws.com/data/all_transactions.json"

	// Craft the request for the page.
	req, _ := http.NewRequest("GET", uri, nil)

	// Make the request.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return []HouseTrade{}, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return []HouseTrade{}, err
	}

	// Define the structure of the response from the API endpoint.
	var body []struct {
		Representative  string `json:"representative"`
		TransactionDate string `json:"transaction_date"`
		DisclosureDate  string `json:"disclosure_date"`
		PtrLink         string `json:"ptr_link"`
		Ticker          string `json:"ticker"`
		Owner           string `json:"owner"`
		Amount          string `json:"amount"`
		Type            string `json:"type"`
	}

	if err := json.Unmarshal(content, &body); err != nil {
		return []HouseTrade{}, err
	}

	r := regexp.MustCompile(`^ *Hon. +`)

	trades := []HouseTrade{}
	for _, trade := range body {
		trades = append(trades, HouseTrade{
			Name:            strings.TrimSpace(r.ReplaceAllString(trade.Representative, "")),
			TransactionDate: ParseDate(trade.TransactionDate),
			DisclosureDate:  ParseDate(trade.DisclosureDate),
			Url:             trade.PtrLink,
			Ticker:          strings.ToUpper(strings.TrimSpace(trade.Ticker)),
			Owner:           strings.Trim(trade.Owner, "- "),
			Amount:          trade.Amount,
			Type:            StringToTransactionType(trade.Type),
		})
	}

	return trades, err
}

func AllHouseTradesOnDate(date time.Time) ([]HouseTrade, error) {
	startOfDay := date.UTC().Truncate(24 * time.Hour).Add(-time.Microsecond)
	endOfDay := startOfDay.Add(24*time.Hour + 2*time.Microsecond)

	allTrades, err := AllHouseTrades()
	if err != nil {
		return allTrades, err
	}

	trades := []HouseTrade{}
	for _, trade := range allTrades {
		if trade.TransactionDate != nil && trade.TransactionDate.After(startOfDay) && trade.TransactionDate.Before(endOfDay) {
			trades = append(trades, trade)
		}
	}

	return trades, err
}

func AllHouseTradesSince(date time.Time) ([]HouseTrade, error) {
	allTrades, err := AllHouseTrades()
	if err != nil {
		return allTrades, err
	}

	trades := []HouseTrade{}
	for _, trade := range allTrades {
		if trade.TransactionDate != nil && trade.TransactionDate.After(date) {
			trades = append(trades, trade)
		}
	}

	return trades, err
}
