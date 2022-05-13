const superagent = require("superagent");
const moment = require("moment");

function toDate(value) {
  const d = moment(value);
  return d.isValid() ? d.toDate() : null;
}

function ticker(title) {
  const m = title.match(/\(([ A-Z.]+)\)/);
  return m.length ? m[1].trim() : null;
}

function fetchTranscripts(size = 40, page = 1, from, to) {
  const today = moment().format("YYYY-MM-DD");

  // Should default to 0 in unix time.
  const f = moment(from || "1970-01-01").unix();
  const t = moment(to || "1970-01-01").unix();

  const url = `https://seekingalpha.com/api/v3/articles?cacheBuster=${today}&filter[category]=earnings::earnings-call-transcripts&filter[since]=${from}&filter[until]=${to}&include=author,primaryTickers,secondaryTickers&isMounting=true&page[size]=${size}&page[number]=${page}`;

  return superagent
    .get(url)
    .set(
      "Accept",
      "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
    )
    .set("Host", "seekingalpha.com")
    .set(
      "User-Agent",
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15"
    )
    .set("Accept-Language", "en-us")
    .set("Accept-Encoding", "gzip, deflate, br")
    .set("Connection", "keep-alive")
    .then((res) => res.body)
    .then((res) => res.data)
    .then((articles) =>
      articles.map((a) => ({
        source: "https://seekingalpha.com/api/v3/articles",
        byline: "Seeking Alpha",
        date: toDate(a.attributes?.publishedOn),
        title: a.attributes?.title,
        tags: ["earnings call"],
        url: a.links?.self
          ? new URL(a.links.self, "https://seekingalpha.com").toString()
          : null,
      }))
    )
    .then((articles) =>
      articles.map((a) => {
        const t = ticker(a.title);
        return t ? { ...a, tags: [t].concat(a.tags) } : a;
      })
    );
}

module.exports = fetchTranscripts;
