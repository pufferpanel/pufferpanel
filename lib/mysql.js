/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Path = require('path')
var Configuration = require(Path.join(__dirname, '../configuration.json'))
var Mysql = require('mysql')
var Connection = Mysql.createConnection({
  host: Configuration.mysql.host,
  user: Configuration.mysql.username,
  password: Configuration.mysql.password,
  database: Configuration.mysql.database
})

module.exports = Connection
