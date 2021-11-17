// https://sequelize.org/master/manual/migrations.html#dynamic-configuration
module.exports = {
  development: {
    username: 'root',
    password: 'null',
    database: 'werf-guide-app',
    host: 'mysql',
    dialect: 'mysql',
  },
  test: {
    username: 'root',
    password: null,
    database: 'werf-guide-app',
    host: 'mysql',
    dialect: 'mysql',
  },
  production: {
    username: process.env.DB_USERNAME,
    password: process.env.DB_PASSWORD,
    database: 'werf-guide-app',
    host: 'mysql',
    dialect: 'mysql',
  },
};
