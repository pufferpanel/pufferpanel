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
var Randomstring = require('randomstring')
var Bcrypt = require('bcrypt')
var Logger = require(Path.join(__dirname, '../../lib/logger.js'))
var Mysql = require(Path.join(__dirname, '../../lib/mysql.js'))
var Notp = require('notp')
var Base32 = require('thirty-two')

Mysql.connect()

function Authentication () {}

Authentication.prototype.validateCredentials = function (requestip, data, callback) {

  var email = data.email
  var password = data.password

  Mysql.query('SELECT password, use_totp, totp_secret FROM users WHERE email = ?', [email], function (error, result) {

    if (error) {
      Logger.error('An error occured with a MySQL operation in Authentication.validateLogin() at step 1.', error)
      return callback('A MySQL error occured.', false)
    }

    if (result.length !== 1) return callback('No account with that information could be found in the system.', false)

    // TOTP Checks
    if (result[0].use_totp === 1) {
      if (!Notp.totp.verify(data.totp_token, Base32.decode(result[0].totp_secret), { time: 30 })) {
        return callback('TOTP Token was invalid.', false)
      }
    }

    var userPassword = Authentication.prototype.updatePasswordHash(result[0].password)

    if (!Bcrypt.compareSync(password, userPassword)) {
      return callback('Email or password was incorrect.', false)
    }

    var session = {
      id: Randomstring.generate(12),
      ip: requestip
    }

    Mysql.query('UPDATE users SET session_id = ?, session_ip = ? WHERE email = ?', [session.id, session.ip, email], function (error, result) {

      if (error) {
        Logger.error('An error occured with a MySQL operation in Authentication.validateLogin() at step 2.', error)
        return callback('A MySQL error occured.', false)
      }

      return callback(session.id, true)

    })

  })

}

Authentication.prototype.TOTPEnabled = function (email, callback) {

  Mysql.query('SELECT use_totp FROM users WHERE email = ?', [email], function (error, result) {

    if (error) {
      Logger.error('An error occured with a MySQL operation in Authentication.TOTPEnabled() at step 1.', error)
      return callback(null, false)
    }

    if (result.length !== 1) return callback(null, false)

    return callback(null, result[0].use_totp)

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
