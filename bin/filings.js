let Parser = require("rss-parser");

const feeds = [
  "https://data.sec.gov/rss?cik=1535527&type=3,4,5&exclude=true&count=100",
  // "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type=10-k&company=&dateb=&owner=include&start=0&count=100&output=atom",
  // "https://www.sec.gov/cgi-bin/browse-edgar?action=getcurrent&CIK=&type=10-q&company=&dateb=&owner=include&start=0&count=100&output=atom",
];

(async () => {
  feeds.forEach(async (url) => {
    let parser = new Parser({
      headers: {
        "User-Agent":
          "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
        "accept-language": "en-US,en",
      },
    });

    let feed = await parser.parseURL(url);

    feed.items.forEach((item) => {
      console.log(
        item
        //   {
        //   creator: item.creator,
        //   title: item.title,
        //   link: item.link,
        //   isoDate: item.isoDate,
        // }
      );
    });
    console.log(feed.items.length);
  });
})();
