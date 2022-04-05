const log4js = require("log4js");

const fetchDividends = require("../lib/dividends");
const event = require("../lib/event");
const logger = require("../lib/logger");

async function main() {
  fetchDividends()
    .then((items) => {
      logger.info(`received ${items.length} records`);

      return Promise.allSettled(
        items.map((item) =>
          event.publish({
            timestamp: new Date(),
            name: "dividend",
            source: "feeds",
            body: item,
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
    })
    .finally(() => {
      log4js.shutdown();
    });
}

main();
