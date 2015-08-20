/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Logger = require('./logger.js');
var Configuration = require('./config.js');

var config = Configuration.rethinkdb || {
  host: 'localhost',
  port: 28015,
  database: 'pufferpanel'
};

try {

  var Rethink = require('rethinkdbdash')({
    host: config.host,
    port: config.port,
    db: config.database
  });

  // Set silent: true to attempt using this function.
  // Rethink.getPoolMaster().on('log', function (message) {
  //   if (message.indexOf('Fail to create a new connection for the connection') > -1) {
  //     Logger.error('There was an error attempting to connect to the RethinkDB server. Automatically attempting reconnection...', message);
  //   }
  // });

} catch (Exception) {
  Logger.error('Unable to connect to the RethinkDB provided in the configuration.', Exception.stack);
}

module.exports = Rethink;
