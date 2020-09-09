/* eslint no-unused-vars: "warn" */
"use strict";

var sqlFixtures = require("sql-fixtures");

var dbConfig = {
  client: "pg",
  connection: process.env.DATABASE_URL,
  searchPath: ["public"],
  pool: { min: 0, max: 1, idleTimeoutMillis: 500 },
};

var options = {
  unique: true,
  showWarning: false,
};

var dataSpec = {
  users: [
    {
      name: "admin",
      email: "admin@krovatka",
      password: "supersecret",
    },
  ],
};

sqlFixtures.create(dbConfig, dataSpec, function (err, result) {
  if (err) {
    return;
  } else {
    console.log(result);
  }
});   