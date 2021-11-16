const express = require('express');
const router = express.Router();

/* Just pong. */
router.get('/', function (req, res, next) {
  res.log.debug('we are being pinged');
  res.send('pong\n');
});

module.exports = router;
