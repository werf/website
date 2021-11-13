const express = require('express');
const router = express.Router();

// [<snippet logger-usage>]
router.get('/', function (req, res, next) {
  res.log.debug('we are being pinged');
  res.send('pong\n');
});
// [<endsnippet logger-usage>]

module.exports = router;
