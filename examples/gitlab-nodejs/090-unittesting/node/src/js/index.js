const express = require("express");
const cookieParser = require("cookie-parser");
const pg = require("pg");
const bodyParser = require("body-parser");
const path = require("path");
const async = require("async");
const app = express();
const http = require("http").Server(app);
const router = express.Router();
const io = require("socket.io")(http);
var amqp = require("amqplib/callback_api");

// Get ENV vars
const PORT = process.env.PORT || 8080;
const REDIS_URI = process.env.SESSION_REDIS || "redis://127.0.0.1:6379";
const SESSION_TTL = process.env.SESSION_TTL || 3600;
const COOKIE_SECRET = process.env.COOKIE_SECRET || "supersecret";
const pgconnectionString =
  process.env.DATABASE_URL || "postgresql://127.0.0.1/postgres";
const AMQP_URI = process.env.AMQP_URI || "amqp://localhost";

// Postgres connect
const pool = new pg.Pool({
  connectionString: pgconnectionString,
});
pool.on("error", (err, client) => { // eslint-disable-line no-unused-vars
  console.error("Unexpected error on idle client", err);
  process.exit(-1);
});

// Redis connect
const expSession = require("express-session");
const redis = require("redis");
let redisClient = redis.createClient(REDIS_URI);
let redisStore = require("connect-redis")(expSession);

var session = expSession({
  store: new redisStore({ client: redisClient, ttl: SESSION_TTL }),
  secret: "keyboard cat",
  resave: false,
  saveUninitialized: false,
});
var sharedsession = require("express-socket.io-session");
app.use(session);

// Use express-session in socket.io
io.use(
  sharedsession(session, {
    autoSave: true,
  })
);

console.log(`Connect to REDIS ${REDIS_URI}`);

app.set("views", path.join(__dirname, "../../", "views"));
app.engine("html", require("ejs").renderFile);
app.use(cookieParser(COOKIE_SECRET));
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
app.use(function(req, res, next) {
  res.locals.key = req.session.key;
  next();
});

io.use(function (socket, next) {
  session(socket.handshake, {}, next);
});

// This is an important function.
// This function does the database handling task.
// We also use async here for control flow.

function handle_database(req, type, callback) {
  async.waterfall(
    [
      function (callback) {
        pool.connect(function (err, connection, release) { // eslint-disable-line no-unused-vars
          if (err) {
            // if there is error, stop right away.
            // This will stop the async code execution and goes to last function.
            callback(true);
          } else {
            callback(null, connection);
          }
        });
      },
      function (connection, callback) {
        var SQLquery;
        switch (type) {
          case "login":
            SQLquery =
              "SELECT * FROM users WHERE email='" +
              req.body.user_email +
              "' AND password='" +
              req.body.user_password +
              "'";
            break;
          case "checkEmail":
            SQLquery =
              "SELECT * FROM users WHERE email='" + req.body.user_email + "'";
            break;
          case "register":
            SQLquery =
              "INSERT INTO users(email,password,name) VALUES ('" +
              req.body.user_email +
              "','" +
              req.body.user_password +
              "','" +
              req.body.user_name +
              "')";
            break;
          default:
            break;
        }
        callback(null, connection, SQLquery);
      },
      function (connection, SQLquery, callback) {
        connection.query(SQLquery, function (err, result) {
          connection.release();
          if (!err) {
            if (type === "login") {
              callback(result.rows.length === 0 ? false : result.rows[0]);
            } else if (type === "checkEmail") {
              callback(result.rows.length === 0 ? false : true);
            } else {
              callback(false);
            }
          } else {
            // if there is error, stop right away.
            // This will stop the async code execution and goes to last function.
            callback(true);
          }
        });
      },
    ],
    function (result) {
      // This function gets call after every async task finished.
      if (typeof result === "boolean" && result === true) {
        callback(null);
      } else {
        callback(result);
      }
    }
  );
}

// socketio function
io.sockets.on("connection", function (socket) {
  socket.on("login", function (userdata) { // eslint-disable-line no-unused-vars
    socket.handshake.session.save();
    socket.username = socket.handshake.session["key"]["name"];
    io.emit("is_online", "🔵 <i>" + socket.username + " join the chat..</i>");
  });

  socket.on("disconnect", function (username) { // eslint-disable-line no-unused-vars
    io.emit("is_online", "🔴 <i>" + socket.username + " left the chat..</i>");
  });

  socket.on("chat_message", function (message) {
    io.emit(
      "chat_message",
      "<strong>" + socket.username + "</strong>: " + message
    );
  });
});

/**
    --- Router Code begins here.
**/

router.get("/", function (req, res) {
  if (req.session.key) {
    res.render("chat.ejs", {
      email: req.session.key["email"],
      name: req.session.key["name"],
    });
  } else {
    res.redirect("/login");
  }
});

router.get("/login", function (req, res) {
    if (req.session.key) {
        res.redirect("/");
      } else {
        res.render("login.ejs");
      }
});

router.post("/login", function (req, res) {
  handle_database(req, "login", function (response) {
    if (response === null) {
      res.json({ error: "true", message: "Database error occured" });
    } else {
      if (!response) {
        res.json({
          error: "true",
          message: "Login failed ! Please register",
        });
      } else {
        req.session.key = response;
        res.json({ error: false, message: "Login success." });
      }
    }
  });
});

router.get("/profile", function (req, res) {
  if (req.session.key) {
    res.render("profile.ejs", {
      email: req.session.key["email"],
      name: req.session.key["name"],
    });
  } else {
    res.redirect("/");
  }
});

router.post("/register", function (req, res) {
  handle_database(req, "checkEmail", function (response) {
    if (response === null) {
      res.json({ error: true, message: "This email is already present" });
    } else {
      handle_database(req, "register", function (response) {
        if (response === null) {
          res.json({ error: true, message: "Error while adding user." });
        } else {
          amqp.connect(AMQP_URI, function (error0, connection) {
            if (error0) {
              throw error0;
            }
            connection.createChannel(function (error1, channel) {
              if (error1) {
                throw error1;
              }
              var queue = "onboard_emails";
              var msg = `{"to":"${req.body.user_email}", "text":"Hello, dear ${req.body.user_name}", "subject":"Chat registration"`;

              channel.assertQueue(queue, {
                durable: true,
              });

              channel.sendToQueue(queue, Buffer.from(msg));
              console.log(" [x] Sent %s", msg);
            });
          });
          res.json({ error: false, message: "Registered successfully." });
        }
      });
    }
  });
});

router.get("/logout", function (req, res) {
  if (req.session.key) {
    req.session.destroy(function () {
      res.redirect("/");
    });
  } else {
    res.redirect("/");
  }
});

app.use("/", router);

const server = http.listen(PORT, function () { // eslint-disable-line no-unused-vars
  console.log(`Server is listening on port ${PORT}`);
});
