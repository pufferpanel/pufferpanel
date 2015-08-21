/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Logger = requireFromRoot('lib/logger');
var Authentication = requireFromRoot('server/auth/authentication');

var AuthService = {};

/**
 * @callback loginUserCallback
 * @param {Boolean} success - Whether the user was successfully logged in
 * @param {Object|String} data - The user data if the login was valid, otherwise a failure message
 */
/**
 * Attempts to log a user. The given credentials are first verified, then a session is created.
 *
 * @param {String} email - Email of user
 * @param {String} password - Password for user
 * @param {String} totptoken - Totp token for user
 * @param {String} ipAddr - IP address to create session for
 * @param {loginUserCallback} callback - Function to handle response
 */
AuthService.loginUser = function (email, password, totptoken, ipAddr, callback) {
  Authentication.validateLogin(email, password, totptoken, function (err, success, data) {

    if (err !== undefined) {
      Logger.error(err);
      return callback(false, data);
    }

    if (!success) {
      return callback(false, data);
    }

    Authentication.createSession(data.id, ipAddr, function (err, sessionId) {
      if (err !== undefined) {
        Logger.error(err);
        return callback(false, 'Error creating session');
      }
      data.session_id = sessionId;
      return callback(true, data);
    });

  });
};

/**
 * Gets if a given user's TOTP option is enabled.
 *
 * @param {String} email - Email of user
 * @returns {Boolean} isEnabled - True if TOTP is enabled, otherwise false
 */
AuthService.isTOTPEnabled = function (email) {
  Authentication.isTOTPEnabled(email, function (err, isEnabled) {
    if (err !== undefined) {
      Logger.error(err);
      return false;
    }
    return isEnabled;
  });
};

module.exports = AuthService;
