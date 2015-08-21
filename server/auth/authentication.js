/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Randomstring = require('randomstring');
var Bcrypt = require('bcrypt');
var Notp = require('notp');
var Base32 = require('thirty-two');

var Authentication = {};

/**
 * @callback validateLoginCallback
 * @param {Error} err - Error that occurred during execution, otherwise undefined
 * @param {Boolean} success - If the login credentials were valid
 * @param {Object|String} data - If err is undefined, then the user is returned, otherwise a message with the failure
 *   reason
 */
/**
 * Determines whether given credentials are valid.
 *
 * @param {String} email - Email address of user
 * @param {String} password - Password for user (may be hashed)
 * @param {String} [totp_token] - TOTP token if the user has totp enabled, otherwise may be omitted
 * @param {validateLoginCallback} callback - Function to call with results
 */
Authentication.validateLogin = function (email, password, totp_token, callback) {

  //check if totp_token was provided. If it was not, then callback is not defined and we need to adjust
  if (callback === undefined) {
    callback = totp_token;
  }

  var Rethink = requireFromRoot('lib/rethink.js');
  Rethink.table('users').filter(Rethink.row('email').eq(email)).run().then(function (users) {

    if (users.length !== 1) {
      return callback(undefined, false, 'No account with that information could be found in the system.');
    }

    var user = users[0];

    if (user.use_totp === 1) {
      if (!Notp.totp.verify(totp_token, Base32.decode(user.totp_secret), { time: 30 })) {
        return callback(undefined, false, 'TOTP token was invalid.');
      }
    }

    if (!Bcrypt.compareSync(password, Authentication.updatePasswordHash(user.password))) {
      return callback(undefined, false, 'Email or password was incorrect.');
    }

    return callback(undefined, true, user);

  }).error(function (err) {
    return callback(err, false, 'There was an error processing this request.');
  });

};

/**
 * @callback createSessionCallback
 * @param {Error} err - Error that occurred, otherwise undefined
 * @param {String} session - Session token
 */
/**
 * Creates a new session for the given user
 *
 * @param {String} userId - ID of user
 * @param {String} ipAddr - IP Address
 * @param {createSessionCallback} callback - Function to handle response
 */
Authentication.createSession = function (userId, ipAddr, callback) {

  var sessionId = Randomstring.generate(12);

  var Rethink = requireFromRoot('lib/rethink');
  Rethink.table('users').get(userId).update({
    session_id: sessionId,
    session_ip: ipAddr
  }).run().then(function () {
    return callback(undefined, sessionId);
  }).error(function (err) {
    return callback(err, undefined);
  });

};

/**
 * @callback isTOTPEnabledCallback
 * @param {Error} err - Error that occurred, otherwise undefined
 * @param {Boolean} isEnabled - True if TOTP is enabled, otherwise false
 */
/**
 * Gets if a given user's TOTP option is enabled.
 *
 * @param {String} email - Email of user
 * @param {isTOTPEnabledCallback} callback - Function to handle response
 */
Authentication.isTOTPEnabled = function (email, callback) {

  var Rethink = requireFromRoot('lib/rethink');
  Rethink.table('users').filter(Rethink.row('email').eq(email)).run().then(function (user) {
    if (user.length !== 1 || user[0].use_totp === 0) {
      return callback(undefined, false);
    }
    return callback(undefined, true);
  }).error(function (err) {
    return callback(err, false);
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
