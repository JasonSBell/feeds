const log4js = require("log4js");

const parseRSSNewsFeed = require("../lib/news");
const event = require("../lib/event");
const logger = require("../lib/logger");

const feeds = [
  "https://feeds.a.dj.com/rss/WSJcomUSBusiness.xml",
  "https://feeds.a.dj.com/rss/RSSMarketsMain.xml",
  "http://feeds.marketwatch.com/marketwatch/topstories/",
  // "http://www.fool.com/feed/",
  "http://seekingalpha.com/feed.xml",
  "https://www.investing.com/rss/news_25.rss",
  "https://www.cnbc.com/id/20409666/device/rss/rss.html?x=1",
  // "https://www.nasdaq.com/feed/rssoutbound?symbol=crwd",
  // "https://seekingalpha.com/api/sa/combined/CRWD.xml",
  // "https://seekingalpha.com/sector/transcripts.xml",
  // "https://investor.docusign.com/rss/PressRelease.aspx?LanguageId=1&CategoryWorkflowId=1cb807d2-208f-4bc3-9133-6a9ad45ac3b0&tags=",
  // "https://investor.docusign.com/rss/event.aspx",
  // "https://ir.crowdstrike.com/rss/events.xml",
  // "http://www.nasdaqtrader.com/rss.aspx?feed=currentheadlines&categorylist=2,6,7",
];

async function main() {
  Promise.allSettled(
    feeds.map((feed) => {
      return parseRSSNewsFeed(feed)
        .then((items) => {
          logger.info(`received ${items.length} records from ${feed}`);

          return Promise.allSettled(
            items.map((item) =>
              event.publish({
                timestamp: new Date(),
                name: "news",
                source: "feeds",
                body: { ...item, feed },
              })
            )
          );
        })
        .then((results) => {
          results.forEach((result) => {
            if (result.status === "rejected") {
              logger.error(`Failed to publish event ${result.reason}`);
            } else {
              logger.info(`Published event ${result.value.id}`);
            }
          });
        });
    })
  ).finally(() => {
    logger.info("done");
    log4js.shutdown();
  });
}

main();
