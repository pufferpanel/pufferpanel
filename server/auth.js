/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Logger = Rfr('lib/logger.js');
var Bcrypt = require('bcrypt');
var Notp = require('notp');
var Base32 = require('thirty-two');
var R = Rfr('lib/rethink.js');
var Authentication = {};

/**
 * Determines whether given credentials are valid.
 *
 * @param {request} object - The Hapi request object.
 * @param {callback} callback - Callback for function.
 */
Authentication.validateLogin = function (request, callback) {

  if (typeof callback === undefined) {
    return callback({ error: 'No callback function was assigned to this request.' });
  }

  R.table('users').filter(R.row('email').eq(request.payload.email)).run().then(function (user) {

    if (user.length !== 1) {
      return callback({ error: 'No account with that information could be found in the system.' });
    }

    user = user[0];
    if (user.use_totp === 1) {
      if (!Notp.totp.verify(request.payload.totp_token, Base32.decode(user.totp_secret), { time: 30 })) {
        return callback({ error: 'TOTP token was invalid.' });
      }
    }

    if (!Bcrypt.compareSync(request.payload.password, Authentication.updatePasswordHash(user.password))) {
      return callback({ error: 'Email or password was incorrect.' });
    }

    return callback({
      success: true,
      session: user
    });

  }).error(function (err) {
    Logger.error(err);
    return callback({ error: 'There was an error processing this request.' });
  });

};

/**
 * Gets if a given user's TOTP option is enabled.
 *
 * @param {String} email - Email of user
 * @param {callback} callback - Callback that returns true or false if TOTP is enabled.
 */
Authentication.isTOTPEnabled = function (email, callback) {

  R.table('users').filter(R.row('email').eq(email)).run().then(function (user) {
    if (user.length !== 1 || user[0].use_totp === 0) {
      return callback(false);
    }
    return callback(true);
  }).error(function (err) {
    Logger.error(err);
    return callback(false);
  });

};

/**
 * Updates a password stored in PHP's BCrypt format to NodeJS's BCrypt format
 *
 * @param {String} password - Password hash to convert
 * @returns {String} Updated password
 */
Authentication.updatePasswordHash = function (password) {
  return password.replace(/^\$2y(.+)$/i, '\$2a$1');
};

/**
 * Generates a {@link bcrypt}-hashed password
 *
 * @param {String} rawpassword - Password to hash
 * @returns {String} Hashed form of the password
 */
Authentication.generatePasswordHash = function (rawpassword) {
  return Bcrypt.hashSync(rawpassword, 10);
};

module.exports = Authentication;
