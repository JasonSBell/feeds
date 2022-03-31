const log4js = require("log4js");

const logger = require("../lib/logger");
const config = require("../lib/config");
const mongo = require("../lib/mongo");
const app = require("../lib/app");

async function main() {
  await mongo.init();

  const server = app.listen(config.port, () => {
    logger.info(`listening at http://localhost:${config.port}`);
  });

  process.on("SIGINT", async () => {
    await server.close();
    await mongo.close();
    await log4js.shutdown();
  });
}

main();
