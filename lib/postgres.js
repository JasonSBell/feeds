const { Sequelize, DataTypes } = require("sequelize");

const config = require("./config");
const logger = require("./logger");

const sequelize = new Sequelize(
  `postgres://${config.postgres.user}:${config.postgres.password}@${config.postgres.host}:${config.postgres.port}/${config.postgres.database}`,
  {
    logging: false,
  }
);

const Company = sequelize.define(
  "Company",
  {
    // Model attributes are defined here
    cik: {
      type: DataTypes.INTEGER,
      allowNull: false,
      primaryKey: true,
    },
    ticker: {
      type: DataTypes.STRING,
      allowNull: false,
      primaryKey: true,
    },
    name: {
      type: DataTypes.STRING,
      allowNull: false,
    },
  },
  {
    tableName: "companies",
    underscored: true,
  }
);

const Earnings = sequelize.define(
  "Earnings",
  {
    // Model attributes are defined here
    date: {
      type: DataTypes.DATEONLY,
      allowNull: false,
      primaryKey: true,
    },
    ticker: {
      type: DataTypes.STRING,
      allowNull: false,
      primaryKey: true,
    },
  },
  {
    tableName: "earnings",
    underscored: true,
  }
);

const Dividend = sequelize.define(
  "Dividend",
  {
    // Model attributes are defined here
    exDate: {
      type: DataTypes.DATEONLY,
      allowNull: false,
      primaryKey: true,
    },
    ticker: {
      type: DataTypes.STRING,
      allowNull: false,
      primaryKey: true,
    },
    dividendRate: {
      type: DataTypes.DOUBLE,
    },
    recordDate: {
      type: DataTypes.DATEONLY,
    },
    paymentDate: {
      type: DataTypes.DATEONLY,
    },
    announcementDate: {
      type: DataTypes.DATEONLY,
    },
  },
  {
    tableName: "dividends",
    underscored: true,
  }
);

const Split = sequelize.define(
  "Split",
  {
    // Model attributes are defined here
    date: {
      type: DataTypes.DATEONLY,
      allowNull: false,
      primaryKey: true,
    },
    ticker: {
      type: DataTypes.STRING,
      allowNull: false,
      primaryKey: true,
    },
    ratio: {
      type: DataTypes.STRING,
    },
    executionDate: {
      type: DataTypes.DATEONLY,
    },
    announcementDate: {
      type: DataTypes.DATEONLY,
    },
  },
  {
    tableName: "splits",
    underscored: true,
  }
);

const Article = sequelize.define(
  "Article",
  {
    // Model attributes are defined here
    date: {
      type: DataTypes.DATEONLY,
      allowNull: false,
      primaryKey: true,
    },
    title: {
      type: DataTypes.STRING,
      allowNull: false,
      primaryKey: true,
    },
    creator: {
      type: DataTypes.STRING,
    },
    link: {
      type: DataTypes.STRING,
      allowNull: false,
    },
  },
  {
    tableName: "articles",
    underscored: true,
  }
);

Split.belongsTo(Company, { foreignKey: "ticker", targetKey: "ticker" });
Earnings.belongsTo(Company, { foreignKey: "ticker", targetKey: "ticker" });
Dividend.belongsTo(Company, { foreignKey: "ticker", targetKey: "ticker" });
Company.hasMany(Earnings, { foreignKey: "ticker" });
Company.hasMany(Dividend, { foreignKey: "ticker" });

module.exports = {
  sequelize,
  init: () => {
    logger.info(
      `connecting to postgres host=${config.postgres.host}  port =${config.postgres.port} database=${config.postgres.database}`
    );
    return sequelize.sync();
  },
  close: () => {
    logger.info(
      `closing connection to postgres host=${config.postgres.host}  port =${config.postgres.port} database=${config.postgres.database}`
    );
    return sequelize.close();
  },
  Split,
  Earnings,
  Company,
  Dividend,
  Article,
};
