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
require('colors');
var Fs = require('fs-extra');
var Winston = require('winston');
var Util = require('util');
var Vargs = require('yargs');

var cliArgs = Vargs.argv;

/** @namespace */
var Logger = {};

var colorMapper = {
    error: 'ERROR'.red,
    warn: 'WARN'.yellow,
    unknown: 'UNKNOWN'.yellow
};

var loggerOptions = {
    logFolder: cliArgs.logFolder || './logs/',
    console: {
        level: cliArgs.consoleLevel || 'info',
        showMeta: cliArgs.showMeta || false
    },
    file: {
        level: cliArgs.fileLevel || 'verbose'
    }
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


var _logger = new (Winston.Logger)({
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


// Make sure that the logs directory exists
Fs.ensureDir('./logs', function (err) {

    if (err !== null) {
        console.error(err);
    }
});


/**
 * Log a DEBUG level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Logger.debug = function (message, data) {

    this.log('debug', message, data);
};


/**
 * Log a VERBOSE level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Logger.verbose = function (message, data) {

    this.log('verbose', message, data);
};


/**
 * Log an INFO level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Logger.info = function (message, data) {

    this.log('info', message, data);
};


/**
 * Log a WARN level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Logger.warn = function (message, data) {

    this.log('warn', message, data);
};


/**
 * Log an ERROR level message
 *
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Logger.error = function (message, data) {

    this.log('error', message, data);
};


/**
 * Log a message of the specified level
 *
 * @param {string} level Level of log
 * @param {string} message - Message to log
 * @param {Object} [data] - Data to log
 */
Logger.log = function (level, message, data) {

    _logger.log(level, message, data ? { meta: data } : null);
};

/**
 * Returns true if the application is running in verbose mode.
 * @return {bool} True if running verbose, false otherwise.
 */
Logger.runningVerbose = function () {
    return (loggerOptions.console.level === 'verbose');
};


//if we are in tests, we should not override the console.log
if (process.env.NODE_ENV !== 'test') {

    console.log = function () {

        Logger.info(formatArgs(arguments));
    };

    console.info = function () {

        Logger.info(formatArgs(arguments));
    };

    console.warn = function () {

        Logger.warn(formatArgs(arguments));
    };

    console.error = function () {

        Logger.error(formatArgs(arguments));
    };

    console.debug = function () {

        Logger.debug(formatArgs(arguments));
    };
}

module.exports = Logger;
