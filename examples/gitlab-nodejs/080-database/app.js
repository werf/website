var express = require('express');
var app = express();
const db = require('./queries')

const bodyParser = require('body-parser')
app.use(bodyParser.json())
app.use(
    bodyParser.urlencoded({
      extended: true,
    })
)

app.get('/', function (req, res) {
  res.send('Hello World! <a href="/labels/">Labels list</a><br><form method="post" action="/labels"><input type="text" name="label" /><input type="submit" /></form>');
});
app.get('/labels', db.getLabels);
app.post('/labels', db.createLabel);

app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});
