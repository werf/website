var express = require('express');
var app = express();
app.get('/api/', function (req, res) {
  res.send('server rendered text');
});
app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});
