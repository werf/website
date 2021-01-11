const pg = require("pg");
const pgconnectionString = "postgresql://" +process.env.POSTGRESQL_LOGIN + ":" +
                                            process.env.POSTGRESQL_PASSWORD + "@" +
                                            process.env.POSTGRESQL_HOST + ":" +
                                            process.env.POSTGRESQL_PORT + "/" +
                                            process.env.POSTGRESQL_DATABASE || "postgresql://127.0.0.1/postgres";
// Postgres connect
const pool = new pg.Pool({
  connectionString: pgconnectionString,
});
pool.on("error", (err, _client) => {
  console.error("Unexpected error on idle client", err);
  process.exit(-1);
});

const getLabels = (request, response) => {
  pool.query('SELECT * FROM labels', (error, results) => {
    if (error) {
      throw error
    }
    response.status(200).json(results.rows)
  })
}
const createLabel = (request, response) => {
  const { label } = request.body
  pool.query('INSERT INTO labels (label) VALUES ($1)', [label], (error, results) => {
    if (error) {
      throw error
    }
    response.status(201).send(`Label added`)
  })
}

module.exports = {
  createLabel,
  getLabels
}