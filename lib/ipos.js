const superagent = require("superagent");
const moment = require("moment");

function toDate(value) {
  const d = moment(value, "MM/DD/YYYY", true);
  return d.isValid() ? d.toDate() : null;
}

function fetchIPOs(date = new Date()) {
  date = new Date(moment(date).format("YYYY-MM-DD"));

  return superagent
    .get("https://api.nasdaq.com/api/ipo/calendar")
    .query(`date=${moment(date).format("YYYY-MM")}`)
    .set(
      "User-Agent",
      "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
    )
    .set("accept-language", "en-US,en")
    .then((res) => res.body)
    .then((data) => data.data)
    .then(
      ({
        priced: { rows: priced },
        upcoming: {
          upcomingTable: { rows: upcoming },
        },
        filed: { rows: filed },
        withdrawn: { rows: withdrawn },
      }) => {
        return upcoming.map((d) => ({
          ...d,
          expectedPriceDate: toDate(d.expectedPriceDate),
        }));
      }
    );
}

module.exports = fetchIPOs;
