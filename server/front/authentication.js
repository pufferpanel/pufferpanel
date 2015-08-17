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
var Bcrypt = require('bcrypt')
var Logger = require(Path.join(__dirname, '../../lib/logger.js'))
var Mysql = require(Path.join(__dirname, '../../lib/mysql.js'))

Mysql.connect()

function Authentication () {}

Authentication.prototype.validateCredentials = function (email, password) {

  Mysql.query('SELECT password FROM users WHERE email = ?', [email], function (error, result) {

    if (error) {
      Logger.error('An error occured with a MySQL operation in Authentication.validateLogin() at step 1.', error)
      return false
    }

    if (result.affectedRows !== 1) return false

    var userPassword = Authentication.prototype.updatePasswordHash(result[0].password)

    return Bcrypt.compareSync(password, userPassword)

  })

}

// Helper function to convert from PHP password_hash
Authentication.prototype.updatePasswordHash = function (password) {
  return password.replace(/^\$2y(.+)$/i, '\$2a$1')
}

Authentication.prototype.generatePasswordHash = function (rawpassword) {
  return Bcrypt.hashSync(rawpassword, 10)
}

module.exports = Authentication
