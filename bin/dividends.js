const log4js = require("log4js");

const fetchDividends = require("../lib/dividends");
const db = require("../lib/postgres");
const logger = require("../lib/logger");

async function main() {
  await db.init();

  fetchDividends()
    .then((items) => {
      logger.info(`received ${items.length} records`);
      console.log(items);
      return db.Dividend.bulkCreate(items, {
        updateOnDuplicate: ["ticker", "exDate"],
      });
    })
    .finally(() => {
      log4js.shutdown();
      db.close();
    });
}

main();
