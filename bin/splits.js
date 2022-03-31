const log4js = require("log4js");

const fetchSplits = require("../lib/splits");
const db = require("../lib/postgres");
const logger = require("../lib/logger");

async function main() {
  await db.init();

  fetchSplits()
    .then((items) => {
      logger.info(`received ${items.length} records`);
      console.log(items);
      return db.Split.bulkCreate(items, {
        updateOnDuplicate: ["ticker"],
      });
    })
    .finally(() => {
      log4js.shutdown();
      db.close();
    });
}

main();
