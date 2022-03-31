const mongoose = require("mongoose");

const config = require("./config");
const logger = require("./logger");

mongoose.connection.on("connected", () => {
  logger.info("connected to mongo");
});
mongoose.connection.on("error", (err) => {
  logger.error("mongo err " + err);
});
mongoose.connection.on("disconnected", () => {
  logger.info("disconnected from mongo");
});
mongoose.connection.on("reconnected", () => {
  logger.info("reconnected to mongo");
});
mongoose.connection.on("close", () => {
  logger.info("connection to mongo closed");
});

const ArticleSchema = new mongoose.Schema({
  url: { type: String, required: true },
  title: { type: String, required: true },
  date: { type: Date, required: true },
  feed: { type: String },
  byline: { type: String },
  publisher: { type: String },
  content: { type: String },
  markdown: { type: String },
  tickers: {
    type: [
      {
        type: String,
      },
    ],
  },
});
ArticleSchema.index({ title: 1, date: 1 }, { unique: true });

const Article = mongoose.model("Article", ArticleSchema);

module.exports = {
  init: () => {
    logger.info(
      `connecting to mongo host=${config.mongo.host}  port=${config.mongo.port} database=${config.mongo.database}`
    );
    return mongoose.connect(
      `mongodb://${config.mongo.user}:${config.mongo.password}@${config.mongo.host}:${config.mongo.port}/${config.mongo.database}?authSource=admin&retryWrites=true&writeConcern=majority`
    );
  },
  close: () => {
    logger.info(
      `closing connection to mongo host=${config.mongo.host}  port=${config.mongo.port} database=${config.mongo.database}`
    );
    return mongoose.connection.close();
  },
  Article,
};
