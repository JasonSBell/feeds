const superagent = require("superagent");
const moment = require("moment");

function toDate(value) {
  const d = moment(value, "MM/DD/YYYY", true);
  return d.isValid() ? d.toDate() : null;
}

async function fetchDividends(date = new Date()) {
  return superagent
    .get("https://api.nasdaq.com/api/calendar/dividends")
    .query(`date=${moment(date).format("YYYY-MM-DD")}`)
    .set(
      "User-Agent",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
    )
    .set("accept-language", "en-US,en")
    .then((res) => res.body)
    .then((data) => data.data.calendar.rows || [])
    .then((dividends) =>
      dividends.map((d) => ({
        name: d.companyName,
        ticker: d.symbol,
        exDate: toDate(d.dividend_Ex_Date),
        dividendRate: d.dividend_Rate,
        recordDate: toDate(d.record_Date),
        paymentDate: toDate(d.payment_Date),
        announcementDate: toDate(d.announcement_Date),
      }))
    );
}

module.exports = fetchDividends;
