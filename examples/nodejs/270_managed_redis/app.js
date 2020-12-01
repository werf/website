const express = require('express');
const bodyParser = require('body-parser');
const redis = require("redis");

// Connection to Redis
const host = process.env.REDIS_HOST;
const port = process.env.REDIS_PORT;

const client = redis.createClient(port, host);
client.on("error", function(error) {
  console.error(error);
  process.exit(1)
});
//
global_errors = [];

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

  let isSent = false
  function echoData(data){
    if (isSent) return;
    res.send(JSON.stringify(data));
    isSent = true;
  }

  let return_dataset = []
  client.keys('*', function (err, keys) {
    let i = 0
    keys.forEach( (l) => {
      client.get(l, function(err, value) {
        i++;
        if (err) {
          echoData({
            "result": "error",
            "comment": err.toString()
          })
        } else {
          return_dataset.push({
            'id': l,
            'label': value
          })
        }
        if (i === keys.length) {
          echoData(return_dataset)
        }
      })
    })
    if (keys.length===0) {
      echoData([])
    }
  })
});
//// Add label ////
app.post('/api/labels', function (req, res) {
  halt_if_errors(global_errors, res);

  client.keys('*', (err, keys) => {
    let parsed_keys = keys.map(val => (isNaN(parseInt(val)))? 0 : parseInt(val));
    let latest_id = Math.max(...parsed_keys);
    latest_id++;

    client.set(latest_id, req.body.label, function (err) {
      if (err) {
        res.send(JSON.stringify({
          "result": "error",
          "comment": err.toString()
        }));
      } else {
        res.send(JSON.stringify({"id": latest_id, "label": req.body.label}));
      }
    });
  });
});
//// Get label ////
app.get('/api/labels/:id', function (req, res) {
  halt_if_errors(global_errors, res);

  client.get(req.params.id, function(err, row) {
    if (err) {
      res.send(JSON.stringify({
        "result": "error",
        "comment": err.toString()
      }));
    } else {
      res.send(JSON.stringify({"id": req.params.id, "label": row})); // TODO: check what the hell is returned here
    }
  });
});
//// Modify label ////
app.post('/api/labels/:id', function (req, res) {
  halt_if_errors(global_errors, res);

  client.set(req.params.id, req.body.label, function(err) {
    if (err) {
      res.send(JSON.stringify({
        "result": "error",
        "comment": err.toString()
      }));
    } else {
      // TODO: check what the hell is returned here
      res.send(JSON.stringify({"id": req.params.id, "label": req.body.label}));
    }
  });
});
//// Delete label ////
app.delete('/api/labels/:id', function (req, res) {
  halt_if_errors(global_errors, res);
  client.del(req.params.id, function(err) {
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
