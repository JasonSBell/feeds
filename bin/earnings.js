const log4js = require("log4js");

const fetchEarnings = require("../lib/earnings");
const db = require("../lib/postgres");
const logger = require("../lib/logger");

async function main() {
  await db.init();

  fetchEarnings()
    .then((items) => {
      logger.info(`received ${items.length} records`);
      console.log(items);
      return db.Earnings.bulkCreate(items, {
        updateOnDuplicate: ["ticker"],
      });
    })
    .finally(() => {
      log4js.shutdown();
      db.close();
    });
}

main();
