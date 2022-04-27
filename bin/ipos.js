const log4js = require("log4js");

const fetchIPOs = require("../lib/ipos");
const event = require("../lib/event");
const logger = require("../lib/logger");

async function main() {
  fetchIPOs()
    .then((items) => {
      logger.info(`received ${items.length} records`);

      return Promise.allSettled(
        items.map((item) => {
          console.log(item);
          event.publish({
            timestamp: new Date(),
            name: "ipo",
            source: "feeds",
            body: item,
          });
        })
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
