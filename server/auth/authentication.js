/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Path = require('path');
var Randomstring = require('randomstring');
var Bcrypt = require('bcrypt');
var Notp = require('notp');
var Base32 = require('thirty-two');
var Rethink = require(Path.join(__dirname, '../../lib/rethink.js'));

var Authentication = {};

Authentication.validateCredentials = function (request, callback) {

  var connection = Rethink.openConnection();

  connection.table('users').filter(Rethink.row('email').eq(request.payload.email)).run().then(function (user) {

    if (user.length !== 1) {
      return callback('No account with that information could be found in the system.', false);
    }

    if (user[0].use_totp === 1) {
      if (!Notp.totp.verify(request.payload.totp_token, Base32.decode(user[0].totp_secret), { time: 30 })) {
        return callback('TOTP Token was invalid.', false, null);
      }
    }

    if (!Bcrypt.compareSync(request.payload.password, Authentication.prototype.updatePasswordHash(user[0].password))) {
      return callback('Email or password was incorrect.', false, null);
    }

    var session = {
      id: Randomstring.generate(12),
      ip: request.info.remoteAddress
    };

    Rethink.table('users').get(user[0].id).update({
      session_id: session.id,
      session_ip: session.ip
    }).run().error(function (err) {
      callback('There was an error creating the session.', false, err);
    });

    return callback(user[0], true, null);

  }).error(function (err) {
    callback('There was an error processing this request.', false, err);
  });

};

Authentication.TOTPEnabled = function (email, callback) {

  var connection = Rethink.openConnection();

  connection.table('users').filter(Rethink.row('email').eq(email)).run().then(function (user) {
    if (user.length !== 1) {
      return callback(null, false);
    }
    return callback(null, user[0].use_totp);
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
