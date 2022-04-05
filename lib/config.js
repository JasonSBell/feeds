require("dotenv-expand").expand(require("dotenv").config());

module.exports = {
  port: Number.parseInt(process.env.PORT) || 8088,
  postgres: {
    port: process.env.POSTGRES_PORT
      ? Number.parseInt(process.env.POSTGRES_PORT)
      : 5432,
    host: process.env.POSTGRES_HOST ? process.env.POSTGRES_HOST : "localhost",
    user: process.env.POSTGRES_USER ? process.env.POSTGRES_USER : "root",
    password: process.env.POSTGRES_PASSWORD
      ? process.env.POSTGRES_PASSWORD
      : "password",
    database: process.env.POSTGRES_DATABASE
      ? process.env.POSTGRES_DATABASE
      : "allokate",
  },
  fluentd: {
    port: process.env.FLUENTD_PORT
      ? Number.parseInt(process.env.FLUENTD_PORT)
      : 24224,
    host: process.env.FLUENTD_HOST ? process.env.FLUENTD_HOST : "localhost",
    prefix: process.env.FLUENTD_PREFIX ? process.env.FLUENTD_PREFIX : "feeds",
    level: process.env.FLUENTD_LEVEL ? process.env.FLUENTD_LEVEL : "info",
  },
  eventServiceAPI: process.env.EVENT_SERVICE_API || "http://localhost:8094",
};

if (require.main === module) {
  console.log(module.exports);
}
