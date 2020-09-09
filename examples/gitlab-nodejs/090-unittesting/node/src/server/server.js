/* eslint no-unused-vars: "warn" */
"use strict";

var express = require("express");
const cookieParser = require("cookie-parser");
const pg = require("pg");
const bodyParser = require("body-parser");
const path = require("path");
const async = require("async");
const app = express();
const http = require("http").Server(app);
var morgan = require("morgan");
const router = express.Router();
const io = require("socket.io")(http);
var amqp = require("amqplib/callback_api");
const Minio = require("minio");
var formidable = require("formidable");
var md5 = require("md5");
var fs = require("fs");
const ioMetrics = require("socket.io-prometheus");
const promClient = require("prom-client");
const promRegister = promClient.register;
const promBundle = require("express-prom-bundle");
const metricsMiddleware = promBundle({ includeMethod: true });

// Get ENV vars
const PORT = Number(process.env.PORT) || 5000;
const REDIS_URI = process.env.SESSION_REDIS || "redis://127.0.0.1:6379";
const SESSION_TTL = process.env.SESSION_TTL || 3600;
const COOKIE_SECRET = process.env.COOKIE_SECRET || "supersecret";
const pgconnectionString =
  process.env.DATABASE_URL || "postgresql://127.0.0.1/postgres";
const AMQP_URI = process.env.AMQP_URI || "amqp://localhost";
const S3_ENDPOINT = process.env.S3_ENDPOINT || "127.0.0.1";
const S3_PORT = Number(process.env.S3_PORT) || 9000;
const TMP_S3_SSL = process.env.S3_SSL || "true";
const S3_SSL = TMP_S3_SSL.toLowerCase() == "true";
const S3_ACCESS_KEY = process.env.S3_ACCESS_KEY || "SECRET";
const S3_SECRET_KEY = process.env.S3_SECRET_KEY || "SECRET";
const S3_BUCKET = process.env.S3_BUCKET || "avatars";
const CDN_PREFIX = process.env.CDN_PREFIX || "http://127.0.0.1:9000";
const RMQ_QUEUE = process.env.RMQ_QUEUE || "onboard_emails";
const RMQ_MAIN_BUS = process.env.RMQ_MAIN_BUS || "main_bus";
// Enable access log
app.use(morgan("combined"));

// Enable metrics for express
app.use(metricsMiddleware);

// S3 client
var s3Client = new Minio.Client({
  endPoint: S3_ENDPOINT,
  port: S3_PORT,
  useSSL: S3_SSL,
  accessKey: S3_ACCESS_KEY,
  secretKey: S3_SECRET_KEY,
});

// Postgres connect
const pool = new pg.Pool({
  connectionString: pgconnectionString,
});
pool.on("error", (err, _client) => {
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

// amqp channel
var global_amqp_channel = null;

console.log(`Connect to REDIS ${REDIS_URI}`);

app.set("views", path.join(__dirname, "/", "views"));
app.engine("html", require("ejs").renderFile);
app.use(cookieParser(COOKIE_SECRET));
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());
app.use(function (req, res, next) {
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
        pool.connect(function (err, connection, _release) {
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

// main bus
amqp.connect(AMQP_URI, function (error0, connection) {
  if (error0) {
    throw error0;
  }

  connection.createChannel(function (err, amqp_channel) {
    global_amqp_channel = amqp_channel;
    amqp_channel.assertExchange(RMQ_MAIN_BUS, "fanout", { durable: true });
    amqp_channel.assertQueue("", { exclusive: true }, function (
      err,
      amqp_queue
    ) {
      amqp_channel.bindQueue(amqp_queue.queue, RMQ_MAIN_BUS, "");
      amqp_channel.consume(amqp_queue.queue, function (msg) {
        try {
          var data = JSON.parse(msg.content.toString());
        } catch (e) {
          console.log("WRONG JSON; " + e.toString());
          return;
        }
        io.emit(data.type, data.message);
      });
    });
  });
});

// socketio function
io.sockets.on("connection", function (socket) {
  socket.on("login", function (_userdata) {
    socket.handshake.session.save();
    socket.username = socket.handshake.session["key"]["name"];
    if (global_amqp_channel) {
      global_amqp_channel.assertExchange(RMQ_MAIN_BUS, "fanout", {
        durable: true,
      });
      global_amqp_channel.publish(
        RMQ_MAIN_BUS,
        "",
        new Buffer.from(
          JSON.stringify({
            type: "is_online",
            message:
              "🔵 <i>" +
              socket.handshake.session["key"]["name"] +
              " join the chat..</i>",
          })
        )
      );
    }
  });

  socket.on("disconnect", function (_username) {
    if (global_amqp_channel) {
      global_amqp_channel.assertExchange(RMQ_MAIN_BUS, "fanout", {
        durable: true,
      });
      global_amqp_channel.publish(
        RMQ_MAIN_BUS,
        "",
        new Buffer.from(
          JSON.stringify({
            type: "is_online",
            message:
              "🔴 <i>" +
              socket.handshake.session["key"]["name"] +
              " left the chat..</i>",
          })
        )
      );
    }
  });

  socket.on("chat_message", function (message) {
    if (socket.handshake.session["key"]["avatarUrl"] === null) {
      if (global_amqp_channel) {
        global_amqp_channel.assertExchange(RMQ_MAIN_BUS, "fanout", {
          durable: true,
        });
        global_amqp_channel.publish(
          RMQ_MAIN_BUS,
          "",
          new Buffer.from(
            JSON.stringify({
              type: "chat_message",
              message:
                '<img height="32" width="32" src="http://placehold.it/32x32"/> <strong>' +
                socket.handshake.session["key"]["name"] +
                "</strong>: " +
                message,
            })
          )
        );
      }
    } else {
      if (global_amqp_channel) {
        global_amqp_channel.assertExchange(RMQ_MAIN_BUS, "fanout", {
          durable: true,
        });
        global_amqp_channel.publish(
          RMQ_MAIN_BUS,
          "",
          new Buffer.from(
            JSON.stringify({
              type: "chat_message",
              message:
                '<img height="32" width="32" src="' +
                CDN_PREFIX +
                "/" +
                socket.handshake.session["key"]["avatarUrl"] +
                '"/> <strong>' +
                socket.handshake.session["key"]["name"] +
                "</strong>: " +
                message,
            })
          )
        );
      }
    }
  });
});

// Metrics for prometheus from socket.io
ioMetrics(io);

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
      avatar: CDN_PREFIX + "/" + req.session.key["avatarUrl"],
    });
  } else {
    res.redirect("/");
  }
});

router.get("/api/get_avatar", function (req, res) {
  if (req.session.key) {
    res.json({
      error: false,
      avatar: CDN_PREFIX + "/" + req.session.key["avatarUrl"],
    });
  } else {
    res.json({ error: true, avatar: "" });
  }
});

router.post("/api/upload_avatar", function (req, res) {
  if (req.session.key) {
    var form = new formidable.IncomingForm();
    form.parse(req, function (err, _fields, files) {
      if (err) {
        console.log(err.stack);
        res.json({ error: true, message: "Upload error" });
      } else {
        var filedata = files["avatar_file"];
        fs.readFile(filedata["path"], function (_err, buf) {
          var md5Hash = md5(buf);
          var s3_path = [
            md5Hash.substr(md5Hash.length - 2),
            md5Hash.substr(2, 2),
            md5Hash,
          ].join("/");
          var metaData = {
            "Content-Type": filedata["type"],
          };
          s3Client.fPutObject(
            S3_BUCKET,
            s3_path,
            filedata["path"],
            metaData,
            function (e) {
              if (e) {
                console.log(e);
                res.json({ error: true, message: "Upload error" });
              } else {
                const query = {
                  text: 'UPDATE users SET "avatarUrl" = $2 WHERE id = $1',
                  values: [req.session.key["id"], s3_path],
                };
                pool.query(query, (err, _result) => {
                  if (err) {
                    console.log(err.stack);
                    res.json({ error: true, message: "Upload error" });
                  } else {
                    req.session.key["avatarUrl"] = s3_path;
                    res.json({ error: false, message: "Avatar uploaded" });
                  }
                });
              }
            }
          );
        });
      }
    });
  } else {
    res.json({ error: true, message: "You are not logged in" });
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
              var msg = `{"to":"${req.body.user_email}","text":"Hello dear ${req.body.user_name}"}`;

              channel.assertQueue(RMQ_QUEUE, {
                durable: true,
              });

              channel.sendToQueue(RMQ_QUEUE, Buffer.from(msg));
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

router.get("/healthz", function (req, res) {
  pool.query("SELECT NOW() as now", (err, result) => {
    if (err) {
      console.log(err.stack);
      res.status(500).send("ERROR");
    } else {
      res.status(200).send("OK");
    }
  });
});

router.get("/metrics", (req, res) => {
  res.set("Content-Type", promRegister.contentType);
  res.end(promRegister.metrics());
});

app.use("/", router);

const server = http.listen(PORT, function (_req, _res) {
  console.log(`Server is listening on port ${PORT}`);
});
