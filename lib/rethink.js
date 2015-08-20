/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Configuration = require('./config.js');

var config = Configuration.rethinkdb || {
  host: 'localhost',
  port: 28015,
  database: 'pufferpanel'
};

var Rethink = {};

Rethink.openConnection = function () {
  var connection = require('rethinkdbdash')({
    host: config.host,
    port: config.port,
    db: config.database
  });
  return connection;
};

module.exports = Rethink;
