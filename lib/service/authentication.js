/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var AuthService = {};

var Auth = requireFromRoot('server/auth/authentication');
var Logger = requireFromRoot('lib/logger');

AuthService.loginUser = function (email, password, totptoken, ip, callback) {
  Auth.validateCredentials(email, password, totptoken, function (data, valid, err) {

    if (err !== undefined) {
      Logger.error('Error talking to database', err);
      return callback(false, data);
    }

    if (!valid) {
      return callback(false, data);
    }

    Auth.createSession(data.userId, function (err, result) {

      if (err !== undefined) {
        Logger.error(err);
        return callback(false, 'An error occured while creating session');
      }

      return callback(true, { user: data, sessionId: result });

    });

  });
};

AuthService.isTOTPEnabledFor = function (email) {

  Auth.isTOTPEnabled(email, function (err, result) {

    if (err !== undefined) {
      Logger.error(err);
      return false;
    }

    return result;

  });

};

module.exports = AuthService;
