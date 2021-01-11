// Get configuration from envs
const REDIS_URI = "redis://"+ "root"+":" + process.env.REDIS_PASSWORD +"@"+process.env.REDIS_HOST+":"+ process.env.REDIS_PORT || "redis://127.0.0.1:6379";
const SESSION_TTL = process.env.SESSION_TTL || 3600;

// Init ExpressJS
var express = require('express');
var app = express();

// Redis
const redis = require("redis");
let redisClient = redis.createClient(REDIS_URI);
// Session
const expSession = require("express-session");

// Connect sessions to Redis & configure
let redisStore = require("connect-redis")(expSession);
var session = expSession({
  secret: "keyboard cat", // TODO: CHANGE ME! 'keyboard cat' is insecure!
  resave: false,
  saveUninitialized: false,
  store: new redisStore({ client: redisClient, ttl: SESSION_TTL }),
});
app.use(session);
console.log(`Connect to REDIS ${REDIS_URI}`);

app.get('/', function (req, res) {
  if (req.session.views) {
    req.session.views++
    res.setHeader('Content-Type', 'text/html')
    res.write('<p>views: ' + req.session.views + '</p>')
    res.write('<p>expires in: ' + (req.session.cookie.maxAge / 1000) + 's</p>')
    res.end()
  } else {
    req.session.views = 1
    res.end('welcome to the session demo. refresh!')
  }
});
app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});
