const express = require('express');
const path = require('path');
const pino = require('pino');
const pinoHttp = require('pino-http');

const pingRouter = require('./routes/ping');
const talkersRouter = require('./routes/talkers');

const logger = pino({ level: process.env.LOG_LEVEL || 'info' });
const db = require('./db/models')((msg) => logger.debug(msg));

const app = express();

app.use(pinoHttp({ logger, useLevel: 'info' }));

app.use('/', talkersRouter(db));
app.use('/ping', pingRouter);

app.use((err, req, res, next) => {
  req.log.error(err);
  res.status(500).send({ body: err.message });
});

module.exports = app;
