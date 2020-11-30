const express = require('express');
const bodyParser = require('body-parser');
const sqlite3 = require('sqlite3');
const mailgun = require("mailgun-js");

// Connection to SQLite
global_errors = [];
let sqlite_file = process.env.SQLITE_FILE;
if (! sqlite_file) {
  console.log('Environment variable SQLITE_FILE is not set! I will use in-memory database.');
  sqlite_file = ':memory:';
}
let db = new sqlite3.Database(sqlite_file, (err) => {
  if (err) {
    global_errors.append('could not connect to SQLite database')
    console.error(err.message);
  }
  console.log('Connected to SQLite database.');
});
// Connection to Mailgun
let mg;
try {
  mg = mailgun({
    apiKey: process.env.MAILGUN_APIKEY,
    domain: process.env.MAILGUN_DOMAIN
  });
} catch (error) {
  console.error(error)
}
// Errors viewer
function halt_if_errors(global_errors, http_res) {
  if (global_errors.length) {
    http_res.send(JSON.stringify({
      "result": "error",
      "comment": global_errors.join("\n")
    }))
    process.exit(1);
  }
}
// TODO: check if db has `labels` table

// App main logic //
let app = express();
app.use(bodyParser.urlencoded({ extended: true }));
app.get('/', function (req, res) {
  halt_if_errors(global_errors, res);
  res.send('Hello World! See /api/labels');
});
//// Get labels ////
app.get('/api/labels', function (req, res) {
  halt_if_errors(global_errors, res);
  let sql = `SELECT * FROM labels`;
  db.all(sql, [], (err, rows) => {
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
  halt_if_errors(global_errors, res);

  let sql = `INSERT INTO \`labels\` (\`label\`) VALUES ('${req.body.label}')`;
  db.run(sql, [], function (err) {
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
  halt_if_errors(global_errors, res);

  let sql = `SELECT * FROM \`labels\` WHERE \`id\` = ${req.params.id}`;
  db.get(sql, [], (err, row) => {
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
  halt_if_errors(global_errors, res);

  let sql = `UPDATE \`labels\` set \`label\`='${req.body.label}' WHERE \`id\` = ${req.params.id}`;
  db.run(sql, [], (err) => {
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
  halt_if_errors(global_errors, res);

  let sql = `DELETE from \`labels\` WHERE \`id\` = ${req.params.id}`;
  db.run(sql, [], (err) => {
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
//// Generate and send report to email ////
app.get('/api/send_report', function (req, res) {
  halt_if_errors(global_errors, res);

  try {
    email = {
      from: "Mailgun Sandbox <postmaster@" + process.env.MAILGUN_DOMAIN + ">",
      to: process.env.REPORT_RECIEVER,
      subject: "Report @ " + new Date().toISOString(),
      text: "Here is report: " + new Date().toISOString()
    }
    mg.messages().send(email, function(err, body){
      if (err) {
        console.error("error while sending e-mail")
        console.error(err)
        res.send(JSON.stringify({
          "result": "error",
          "comment": err.toString()
        }));
      } else {
        res.send(JSON.stringify({"result": true}));
      }
    });
  } catch (err) {
    res.send(JSON.stringify({
      "result": "error",
      "comment": err.toString()
    }));
  }
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
