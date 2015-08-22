/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
require('date-utils');
var Fs = require('fs-extra');
var Path = require('path');
var Winston = require('winston');
var Util = require('util');
var Yargs = require('yargs').argv;
require('colors');

var Log = {};

var colorMapper = {
  error: 'ERROR'.red,
  warn: 'WARN'.yellow
};

var loggerOptions = {
  logFolder: Yargs.logFolder || './logs/',
  console: {
    level: Yargs.consoleLevel || 'info',
    showMeta: Yargs.showMeta || false
  },
  file: {
    level: Yargs.fileLevel || 'verbose'
  }
};

// Make sure that the logs directory exists
Fs.ensureDir('./logs', function (err) {
  if (err !== null) {
    console.log(err);
  }
});

/**
 * Log a DEBUG level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Log.debug = function (message, data) {
  this.log('debug', message, data);
};

/**
 * Log a VERBOSE level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Log.verbose = function (message, data) {
  this.log('verbose', message, data);
};

/**
 * Log an INFO level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Log.info = function (message, data) {
  this.log('info', message, data);
};

/**
 * Log a WARN level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Log.warn = function (message, data) {
  this.log('warn', message, data);
};

/**
 * Log an ERROR level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Log.error = function (message, data) {
  this.log('error', message, data);
};

/**
 * Log a message of the specified level
 *
 * @param {string} level Level of log
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Log.log = function (level, message, data) {
  logger.log(level, message, data ? { meta: data } : null);
};

var formatFile = function (options) {

  return Util.format(
    '[%s] [%s] %s%s',
    new Date().toFormat('YYYY-MM-DD HH24:MI:SS'),
    options.level.toUpperCase(),
    options.message,
    options.meta && Object.keys(options.meta).length ? '\n' + JSON.stringify(options.meta, null, 4) : ''
  );
};

var formatConsole = function (options) {
  return Util.format(
    '[%s] [%s] %s%s',
    new Date().toFormat('YYYY-MM-DD HH24:MI:SS'),
    colorMapper[options.level.toLowerCase()] || options.level.toUpperCase(),
    options.message,
    loggerOptions.console.showMeta && options.meta && Object.keys(options.meta).length ? '\n' + JSON.stringify(options.meta, null, 4) : ''
  );
};

var formatArgs = function (args) {
  return Util.format.apply(Util.format, Array.prototype.slice.call(args));
};

var logger = new (Winston.Logger)({
  transports: [
    new (Winston.transports.DailyRotateFile)({
      dirname: loggerOptions.logFolder,
      level: loggerOptions.file.level,
      json: false,
      formatter: formatFile,
      handleExceptions: true
    }),
    new (Winston.transports.Console)({
      level: loggerOptions.console.level,
      formatter: formatConsole,
      handleExceptions: true
    })
  ],
  exitOnError: false
});

//Override console functions so all logging is through Winston
if (process.env.NODE_ENV !== 'test') {
  console.log = function () {
    Log.info(formatArgs(arguments));
  };

  console.info = function () {
    Log.info(formatArgs(arguments));
  };

  console.warn = function () {
    Log.warn(formatArgs(arguments));
  };

  console.error = function () {
    Log.error(formatArgs(arguments));
  };

  console.debug = function () {
    Log.debug(formatArgs(arguments));
  };
}

module.exports = Log;
