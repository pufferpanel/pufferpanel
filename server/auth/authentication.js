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
var Database = requireFromRoot('lib/database');
var Authentication = {};

Authentication.validateCredentials = function (email, password, totptoken, callback) {

  Database.get('users', 'email', email, function (user, err) {

      if (err !== undefined) {
        return callback('There was an error processing this request.', false, err);
      }

      if (user.length !== 1) {
        return callback('No account with that information could be found in the system.', false, null);
      }

      if (user[0].use_totp === 1) {

        if (!Notp.totp.verify(totptoken, Base32.decode(user[0].totp_secret), { time: 30 })) {
          return callback('TOTP Token was invalid.', false, null);
        }

      }

      if (!Bcrypt.compareSync(password, Authentication.prototype.updatePasswordHash(user[0].password))) {
        return callback('Email or password was incorrect.', false, null);
      }

      return callback(user[0], true, null);

    }
  );

};

Authentication.createSession = function (userId, ipAddr, callback) {

  var session = {
    session_id: Randomstring.generate(12),
    session_ip: ipAddr
  };

  Database.set('users', userId, session, function (err) {

    if (err !== undefined) {
      return callback(err, null);
    }

    return callback(null, session.session_id);

  });

};

Authentication.isTOTPEnabled = function (email, callback) {

  Database.get('users', 'email', email, function (user, err) {

    if (err !== undefined) {
      return callback(err, false);
    }

    return callback(null, user.length === 1 ? user[0].use_totp : false);

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
