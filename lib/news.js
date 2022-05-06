let Parser = require("rss-parser");

async function parseRSSNewsFeed(url) {
  const parser = new Parser({
    headers: {
      "User-Agent":
        "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
      "accept-language": "en-US,en",
    },
  });

  const source = await parser.parseURL(url);

  return source.items.map((item) => ({
    source: url,
    byline: item.creator,
    title: item.title,
    url: item.link,
    date: new Date(item.isoDate),
  }));
}

module.exports = parseRSSNewsFeed;
