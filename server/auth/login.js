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
var Logger = require(Path.join(__dirname, '../../lib/logger.js'))
var Mysql = require(Path.join(__dirname, '../../lib/mysql.js'))

Mysql.connect()

function Authentication () {}

Authentication.prototype.testlogin = function () {

  Logger.info('Called testlogin() function successfully.')
  Mysql.query('SELECT * FROM users', function (error, rows, fields) {
    if (error) throw error
  })

}

module.exports = Authentication
