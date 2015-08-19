/**
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
var Winston = require('winston');
var Util = require('util');
require('colors');

var Log = {};

var colorMapper = {
  error: 'ERROR'.red,
  warn: 'WARN'.yellow
};

Log.prepare = function (options) {

  // Make sure that the logs directory exists
  Fs.ensureDir('./logs', function (err) {
    if (err !== null) {
      console.log(err);
    }
  });

  logger = new (Winston.Logger)({
    transports: [
      new (Winston.transports.DailyRotateFile)({
        dirname: './logs/',
        level: options.fileLevel || 'verbose',
        json: false,
        formatter: formatFile,
        handleExceptions: true
      }),
      new (Winston.transports.Console)({
        level: options.consoleLevel || 'info',
        formatter: formatConsole,
        handleExceptions: true
      })
    ]
  });

  //Override console functions so all logging is through Winston
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

};

Log.debug = function (message, data) {

  this.log('debug', message, data);

};

Log.verbose = function (message, data) {

  this.log('verbose', message, data);

};

Log.info = function (message, data) {

  this.log('info', message, data);

};

Log.warn = function (message, data) {

  this.log('warn', message, data);

};

Log.error = function (message, data) {

  this.log('error', message, data);

};

Log.log = function (level, message, data) {

  logger.log(level, message, data ? { meta: data } : null);

};

var formatFile = function (options) {
  return Util.format(
    '[%s] [%s] %s%s',
    new Date().toFormat('YYYY-MM-DD HH24:MI:SS'),
    options.level.toUpperCase(),
    options.message,
    options.meta && Object.keys(options.meta).length ? '\n\t' + JSON.stringify(options.meta) : ''
  );
};

var formatConsole = function (options) {
  return Util.format(
    '[%s] [%s] %s',
    new Date().toFormat('YYYY-MM-DD HH24:MI:SS'),
    colorMapper[options.level.toLowerCase()] || options.level.toUpperCase(),
    options.message
  );
};

var formatArgs = function (args) {
  return Util.format.apply(Util.format, Array.prototype.slice.call(args));
};

module.exports = Log;
