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

Authentication.validateLogin = function (request, callback) {

  Rethink.table('users').filter(Rethink.row('email').eq(request.payload.email)).run().then(function (user) {

    if (user.length !== 1) {
      return callback({ error: 'No account with that information could be found in the system.' });
    }

    if (user[0].use_totp === 1) {
      if (!Notp.totp.verify(request.payload.totp_token, Base32.decode(user[0].totp_secret), { time: 30 })) {
        return callback({ error: 'TOTP token was invalid.' });
      }
    }

    if (!Bcrypt.compareSync(request.payload.password, Authentication.updatePasswordHash(user[0].password))) {
      return callback({ error: 'Email or password was incorrect.' });
    }

    var session = {
      id: Randomstring.generate(12),
      ip: request.info.remoteAddress
    };

    Rethink.table('users').get(user[0].id).update({
      session_id: session.id,
      session_ip: session.ip
    }).run().error(function (err) {
      return callback({ error: 'An error occured creating the session in the database.' });
    });

    return callback({
      success: true,
      session: user[0]
    });

  }).error(function (err) {
    Logger.error(err);
    return callback({ error: 'There was an error processing this request.' });
  });

};

Authentication.TOTPEnabled = function (email, callback) {

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
