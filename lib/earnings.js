const superagent = require("superagent");
const moment = require("moment");

function toDate(value) {
  const d = moment(value, "MM/DD/YYYY", true);
  return d.isValid() ? d.toDate() : null;
}

function fetchEarnings(date = new Date()) {
  return superagent
    .get("https://api.nasdaq.com/api/calendar/earnings")
    .query(`date=${moment(date).format("YYYY-MM-DD")}`)
    .set(
      "User-Agent",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
    )
    .set("accept-language", "en-US,en")
    .then((res) => res.body)
    .then((data) => data.data.rows)
    .then((earnings) => {
      // Earnings may be null if no one reported that day.
      if (!earnings) {
        return [];
      }

      const list = earnings.map((e) => ({
        date: toDate(date),
        ticker: e.symbol,
      }));

      // Once before (8-18-2022 data) Nasdaq has included the same ticker more than once. So dedupe the array.
      let unique = [];
      const l = [];
      list.forEach((e) => {
        if (!unique.includes(e.ticker)) {
          unique.push(e.ticker);
          l.push(e);
        }
      });

      return l;
    });
}

module.exports = fetchEarnings;
