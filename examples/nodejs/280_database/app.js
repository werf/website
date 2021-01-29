const express = require('express');
const bodyParser = require('body-parser');

// Connection to PostgreSQL
const Pool = require('pg').Pool
let client;
try {
  client = new Pool({
    user: process.env.POSTGRESQL_USER,
    host: process.env.POSTGRESQL_HOST,
    database: process.env.POSTGRESQL_DB,
    password: process.env.POSTGRESQL_PASSWORD,
    port: process.env.POSTGRESQL_PORT,
  })
} catch (err) {
  console.error('connection to PSQL failed')
  console.error(err);
  process.exit(1);
}

// App main logic //
let app = express();
app.use(bodyParser.urlencoded({ extended: true }));
app.get('/', function (req, res) {
  halt_if_errors(global_errors, res);
  res.send('Hello World! See /api/labels');
});
//// Get labels ////
app.get('/api/labels', function (req, res) {
  let sql = `SELECT * FROM labels`;
  client.query(sql, (err, rows) => {
    if (err) {
      res.send(JSON.stringify({
        "result": "error",
        "comment": err.toString()
      }));
      console.error(err);
    } else {
      res.send(JSON.stringify(rows));
    }
  })
});
//// Add label ////
app.post('/api/labels', function (req, res) {
  let sql = `INSERT INTO \`labels\` (\`label\`) VALUES ('${req.body.label}')`;
  client.query(sql, function (err, rows) {
    if (err) {
      res.send(JSON.stringify({
        "result": "error",
        "comment": err.toString()
      }));
    } else {
      res.send(JSON.stringify({"id": this.lastID, "label": req.body.label}));
    }
  });
});
//// Get label ////
app.get('/api/labels/:id', function (req, res) {
  let sql = `SELECT * FROM \`labels\` WHERE \`id\` = ${req.params.id}`;
  client.query(sql, (err, row) => {
    if (err) {
      res.send(JSON.stringify({
        "result": "error",
        "comment": err.toString()
      }));
    } else {
      res.send(JSON.stringify(row));
    }
  });
});
//// Modify label ////
app.post('/api/labels/:id', function (req, res) {
  let sql = `UPDATE \`labels\` set \`label\`='${req.body.label}' WHERE \`id\` = ${req.params.id}`;
  client.query(sql, (err, rows) => {
    if (err) {
      res.send(JSON.stringify({
        "result": "error",
        "comment": err.toString()
      }));
    } else {
      res.send(JSON.stringify({"id": req.params.id, "label": req.body.label}));
    }
  });
});
//// Delete label ////
app.delete('/api/labels/:id', function (req, res) {
  let sql = `DELETE from \`labels\` WHERE \`id\` = ${req.params.id}`;
  client.query(sql, (err, rows) => {
    if (err) {
      res.send(JSON.stringify({
        "result": "error",
        "comment": err.toString()
      }));
    } else {
      res.send(JSON.stringify({"result": true}));
    }
  });
});
////
app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});

// Closing DB connection on sigterm
process.on('SIGTERM', () => {
  db.close((err) => {
    if (err) return console.error(err.message);
    console.log('Close the database connection.');
  });
})
