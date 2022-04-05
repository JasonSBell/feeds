const superagent = require("superagent");

const config = require("./config");

function publish({ timestamp, name, source, body }) {
  if (!timestamp) {
    timestamp = new Date();
  }

  return superagent
    .put(new URL("/api/events", config.eventServiceAPI).toString())
    .send({ timestamp, name, source, body })
    .then((res) => res.body);
}

module.exports = {
  publish,
};
