const log4js = require("log4js");

const fetchHouseTrades = require("../lib/house");
const event = require("../lib/event");
const logger = require("../lib/logger");

async function main() {
  fetchHouseTrades()
    .then((items) => {
      logger.info(`received ${items.length} records`);

      return Promise.allSettled(
        items.map((item) =>
          event.publish({
            timestamp: new Date(),
            name: "congressional_trade",
            source: "feeds",
            body: { ...item, body: "House of Representatives" },
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
