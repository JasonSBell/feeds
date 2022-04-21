const log4js = require("log4js");

const fetchEarnings = require("../lib/earnings");
const event = require("../lib/event");
const logger = require("../lib/logger");

async function main() {
  fetchEarnings(new Date("2022-04-18"))
    .then((items) => {
      logger.info(`received ${items.length} records`);

      return Promise.allSettled(
        items.map((item) =>
          event.publish({
            timestamp: new Date(),
            name: "earnings",
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
