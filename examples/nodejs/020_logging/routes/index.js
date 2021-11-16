const express = require('express');
const router = express.Router();

/* GET home page. */
router.get('/', function (req, res, next) {
  res.render('index', { title: 'Express' });
});

router.get('/err', function (req, res, next) {
  throw new Error("Hello from an unhandler error")
});


module.exports = router;
