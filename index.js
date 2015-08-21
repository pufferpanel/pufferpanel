/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
global.requireFromRoot = function (name) {
  return require(__dirname + '/' + name);
};
var Yargs = require('yargs').argv;
var Logger = requireFromRoot('lib/logger');

Logger.prepare(Yargs);

Logger.info('+ =========================================== +');
Logger.info('| PufferPanel logs all information and errors |');
Logger.info('| into the logs/ directory. Please check      |');
Logger.info('| there before asking for help with bugs.     |');
Logger.info('|                                             |');
Logger.info('| \x1b[41mSubmit bug reports at the following link:\x1b[0m   |');
Logger.info('| https://github.com/PufferPanel/PufferPanel  |');
Logger.info('+ =========================================== +');

// Include HapiJS Routing Mechanisms
requireFromRoot('lib/routes');

process.on('SIGINT', function () {

  Logger.warn('Recieved SIGINT. Preparing for shutdown...');
  Logger.info('All shutdown parameters complete. Stopping...\n');
  process.exit();

});
