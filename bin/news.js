const log4js = require("log4js");

const parseRSSNewsFeed = require("../lib/news");
const db = require("../lib/mongo");
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
  await db.init();

  Promise.allSettled(
    feeds.map((feed) => {
      return parseRSSNewsFeed(feed)
        .then((items) => {
          logger.info(`received ${items.length} records from ${feed}`);
          return db.Article.insertMany(items);
        })
        .catch((err) => {
          if (err) {
            if (err.code === 11000) {
              // Duplicate
            } else {
              logger.error("article insert error " + err);
            }
          }
        });
    })
  )
    .then((results) => {
      results.forEach((result, i) => {
        if (result.status === "rejected") {
          logger.error(`failed to process feed ${feeds[i]}`);
        }
      });
    })
    .finally(() => {
      logger.info("done");
      log4js.shutdown();
      db.close();
    });
}

main();
