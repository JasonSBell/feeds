const express = require("express");
const { Router } = require("express");
const prometheusExporter = require("express-prometheus-middleware");
const { body, query, validationResult } = require("express-validator");
const bodyParser = require("body-parser");
var cors = require("cors");
const log4js = require("log4js");
const moment = require("moment");
const op = require("sequelize").Op;

const logger = require("./logger");
const db = require("./postgres");

function validateRequest() {
  return (req, res, next) => {
    const errors = validationResult(req);
    if (!errors.isEmpty()) {
      return res.status(400).json({ errors: errors.array() });
    } else {
      next();
    }
  };
}

const app = express();

app.use(
  cors({
    origin: true,
  })
);
app.use(bodyParser.json({ limit: "50mb" }));

app.use(
  prometheusExporter({
    metricsPath: "/metrics",
    collectDefaultMetrics: true,
    requestDurationBuckets: [0.1, 0.5, 1, 1.5],
    requestLengthBuckets: [512, 1024, 5120, 10240, 51200, 102400],
    responseLengthBuckets: [512, 1024, 5120, 10240, 51200, 102400],
  })
);

app.use(log4js.connectLogger(logger, { level: "debug" }));

api = Router();

api.get("/ping", (req, res) => {
  res.json({ message: "pong" });
});

api.get(
  "/earnings",
  [query("date").optional().isDate(), validateRequest()],
  (req, res) => {
    db.Earnings.findAll({
      where: {
        date: req.query.date || moment().format("YYYY-MM-DD"),
      },
      attributes: ["date", "ticker", "Company.cik", "Company.name"],
      raw: true,
      include: [
        {
          model: db.Company,
          attributes: [],
          required: false,
        },
      ],
    })
      .then((items) => res.json(items))
      .catch((err) => res.status(500).json({ message: err }));
  }
);

api.get(
  "/dividends",
  [query("date").optional().isDate(), validateRequest()],
  (req, res) => {
    db.Dividend.findAll({
      where: {
        exDate: req.query.date || moment().format("YYYY-MM-DD"),
      },
      attributes: [
        "exDate",
        "ticker",
        "recordDate",
        "paymentDate",
        "announcementDate",
        "Company.cik",
        "Company.name",
      ],
      raw: true,
      include: [
        {
          model: db.Company,
          attributes: [],
          required: false,
        },
      ],
    })
      .then((items) => res.json(items))
      .catch((err) => res.status(500).json({ message: err }));
  }
);

api.get(
  "/splits",
  [query("date").optional().isDate(), validateRequest()],
  (req, res) => {
    db.Split.findAll({
      where: {
        // date: req.query.date || moment().format("YYYY-MM-DD"),
      },
      attributes: [
        "date",
        "ticker",
        "ratio",
        "executionDate",
        "announcementDate",
        "Company.cik",
        "Company.name",
      ],
      raw: true,
      include: [
        {
          model: db.Company,
          attributes: [],
          required: false,
        },
      ],
    })
      .then((items) => res.json(items))
      .catch((err) => res.status(500).json({ message: err }));
  }
);

api.get(
  "/activity",
  [
    query("tickers").exists().withMessage("is required").bail().isString(),
    query("after").optional().isISO8601().bail().toDate(),
    query("before").optional().isISO8601().bail().toDate(),
    validateRequest(),
  ],
  (req, res) => {
    console.log("here");
    const tickers = req.query.tickers.toUpperCase().split(",");
    const before = req.query.before;
    const after = req.query.after;

    const query = {
      where: {
        ticker: {
          [op.in]: tickers,
        },
      },
    };

    if (before) {
      query.where.date = query.where.date || {};
      query.where.date["$lt"] = before;
    }

    if (after) {
      query.where.date = query.where.date || {};
      query.where.date["$gte"] = after;
    }

    console.log(query);

    Promise.all([
      db.Earnings.findAll(query),
      db.Dividend.findAll(query),
      db.Split.findAll(query),
    ])
      .then(([earnings, dividends, splits]) =>
        res.json({ tickers, earnings, dividends, splits })
      )
      .catch((err) => res.status(500).json({ message: err }));
  }
);

app.use("/api/feeds", api);

module.exports = app;
