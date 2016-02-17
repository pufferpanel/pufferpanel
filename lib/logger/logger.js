/*
 Copyright 2016 Joshua Taylor <lordralex@gmail.com>

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
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

var parseOpt = function (option, def) {
    if (option === undefined) {
        return def;
    } else if (option === 'true') {
        return true;
    } else if (option === 'false') {
        return false;
    }
    return option;
};

var loggerOptions = {
    logFolder: cliArgs.logFolder || './logs/',
    console: {
        level: parseOpt(cliArgs.consoleLevel, 'info'),
        showMeta: parseOpt(cliArgs.showMeta, false)
    },
    file: {
        level: parseOpt(cliArgs.fileLevel, 'verbose')
    }
};

var redirectConsole = parseOpt(cliArgs.redirectConsole, true);

var formatFile = function (options) {

    var data = Object.keys(options.meta).length || options.meta instanceof Error ? options.meta : null;

    return Util.format(
        '[%s] [%s] %s%s',
        new Date().toFormat('YYYY-MM-DD HH24:MI:SS'),
        options.level.toUpperCase(),
        options.message,
        data ? '\n' + formatData(data) : ''
    );
};

var formatConsole = function (options) {

    var data = options.meta;
    if (!(data instanceof Error)) {
        if (Object.keys(options.meta).length === 0 || !loggerOptions.console.showMeta) {
            data = null;
        }
    }

    return Util.format(
        '[%s] [%s] %s%s',
        new Date().toFormat('YYYY-MM-DD HH24:MI:SS'),
        colorMapper[options.level.toLowerCase()] || options.level.toUpperCase(),
        options.message,
        data ? '\n' + formatData(data) : ''
    );
};

var formatData = function (data) {

    if (data instanceof Error) {
        return data.stack;
    }

    if (!data) {
        return JSON.stringify(data, null, 4);
    }

    return '';
};


var formatArgs = function (args) {

    return Util.format.apply(Util.format, Array.prototype.slice.call(args));
};

/**
 * Initializes the logging system
 * @private
 */
Logger._initialize = function () {
    Logger._logger = new (Winston.Logger)({
        transports: [
            new (Winston.transports.File)({
                filename: loggerOptions.logFolder + 'proc.log',
                maxsize: 1000000,
                maxFiles: 10,
                level: loggerOptions.file.level,
                json: false,
                formatter: formatFile,
                zippedArchive: true,
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
    // Also makes sure the installer logs folder exists
    Fs.ensureDir('./logs/installer', function (err) {

        if (err !== null) {
            console.error(err);
        }
    });

    if (redirectConsole) {

        Logger.debug('Overriding console logging functions');

        //Store old references
        Logger._oldConsole = {
            log: console.log,
            info: console.info,
            warn: console.warn,
            error: console.error,
            debug: console.debug
        };

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
};


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

    if (data === undefined) {
        data = null;
    }

    Logger._logger.log(level, message, data);
};

/**
 * Returns true if the application is running in verbose mode.
 * @return {boolean} True if running verbose, false otherwise.
 */
Logger.runningVerbose = function () {

    return (loggerOptions.console.level === 'verbose');
};

Logger.shutdown = function () {
    var oldConsole = Logger._oldConsole;

    if (oldConsole) {
        console.log = oldConsole.log;
        console.info = oldConsole.info;
        console.warn = oldConsole.warn;
        console.error = oldConsole.error;
        console.debug = oldConsole.debug;
    }

    Logger._logger.close();
};

process.on('uncaughtException', function (ex) {

    Logger.error('Unhandled exception occurred', ex);
});

module.exports = Logger;