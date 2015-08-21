/**
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
 * @param {Error} err Error that occurred during execution, otherwise undefined
 * @param {boolean} success If the login credentials were valid
 * @param {Object|string} data If err is undefined, then the user is returned, otherwise a message with the failure
 *   reason
 */
/**
 * Determines whether given credentials are valid.
 *
 * @param {string} email Email address of user
 * @param {string} password Password for user (may be hashed)
 * @param {string} totp_token TOTP token if the user has totp enabled, otherwise may be omitted
 * @param {validateLoginCallback} callback Function to call with results
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
 * @callback loginUserCallback
 * @param {Error} err Error that occurred, otherwise undefined
 * @param {boolean} success Whether the user was successfully logged in
 * @param {Object|string} data The user data if the login was valid, otherwise a failure message
 */
/**
 * Attempts to log a user. The given credentials are first verified, then a session is created.
 *
 * @param {string} email Email of user
 * @param {string} password Password for user
 * @param {string} totptoken Totp token for user
 * @param {string} ipAddr IP address to create session for
 * @param callback Function to handle response
 */
Authentication.loginUser = function (email, password, totptoken, ipAddr, callback) {
  Authentication.validateLogin(email, password, totptoken, function (err, success, data) {
    if (err !== undefined) {
      return callback(err, false, data);
    }

    if (!success) {
      return callback(undefined, false, data);
    }

    Authentication.createSession(data.id, ipAddr, function (err, sessionId) {
      data.session_id = sessionId;
      return callback(err, err === undefined, data);
    });
  });
};

Authentication.isTOTPEnabled = function (email, callback) {

  var Rethink = requireFromRoot('lib/rethink');
  Rethink.table('users').filter(Rethink.row('email').eq(email)).run().then(function (user) {
    if (user.length !== 1 || user[0].use_totp === 0) {
      return callback(false);
    }
    return callback(true);
  }).error(function (err) {
    Logger.error(err);
  });

};

// Helper function to convert from PHP password_hash
Authentication.updatePasswordHash = function (password) {
  return password.replace(/^\$2y(.+)$/i, '\$2a$1');
};

Authentication.generatePasswordHash = function (rawpassword) {
  return Bcrypt.hashSync(rawpassword, 10);
};

module.exports = Authentication;
