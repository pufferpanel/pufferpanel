/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
require('colors');
var Rfr = require('rfr');
var Logger = Rfr('lib/logger');

Logger.info('+ =========================================== +');
Logger.info('| PufferPanel logs all information and errors |');
Logger.info('| into the logs/ directory. Please check      |');
Logger.info('| there before asking for help with bugs.     |');
Logger.info('|                                             |');
Logger.info('| '.reset + 'Submit bug reports at the following link: '.red + '  |');
Logger.info('| https://github.com/PufferPanel/PufferPanel  |');
Logger.info('+ =========================================== +');

// Include HapiJS Routing Mechanisms
Rfr('lib/routes.js');

process.on('SIGINT', function () {

    Logger.warn('Recieved SIGINT. Preparing for shutdown...');
    Logger.info('All shutdown parameters complete. Stopping...\n');
    process.exit();
});
