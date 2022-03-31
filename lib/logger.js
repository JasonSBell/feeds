const log4js = require("log4js");

const config = require("./config");

log4js.configure({
  appenders: {
    fluent: {
      type: "log4js-fluent-appender",
      tag_prefix: config.fluentd.prefix,
      options: {
        levelTag: true,
        host: config.fluentd.host,
        port: config.fluentd.port,
      },
    },
    console: { type: "console" },
  },
  categories: {
    default: {
      appenders: ["fluent", "console"],
      level: config.fluentd.level,
    },
  },
});

module.exports = log4js.getLogger();
