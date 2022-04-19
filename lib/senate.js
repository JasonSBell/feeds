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

const superagent = require("superagent");
const moment = require("moment");

function toDate(value) {
  const d = moment(value, "MM/DD/YYYY", true);
  return d.isValid() ? d.toDate() : null;
}

function fetchSenateTrades(date = new Date()) {
  return superagent
    .get(
      "https://senate-stock-watcher-data.s3-us-west-2.amazonaws.com/aggregate/all_transactions.json"
    )
    .then((res) => res.body)
    .then((trades) => {
      return trades.map((e) => ({
        transactionDate: toDate(e.transaction_date),
        disclosureDate: toDate(e.disclosure_date),
        url: e.ptr_link,
        name: e.senator,
        owner: e.owner,
        ticker: e.ticker,
        assetType: e.asset_type,
        type: e.type,
        comment: e.comment,
        amount: e.amount,
      }));
    });
}

module.exports = fetchSenateTrades;
