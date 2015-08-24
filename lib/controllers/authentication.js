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
var Bcrypt = require('bcrypt');
var Notp = require('notp');
var Base32 = require('thirty-two');
var Users = Rfr('lib/models/users.js');
var Authentication = {};


/**
 * Determines whether given credentials are valid.
 *
 * @param {String} email - User's email
 * @param {String} password - User's password
 * @param {String} totpToken - Request's TOTP token (may be undefined)
 * @param {callback} next - Callback to handle response
 */
Authentication.loginUser = function (email, password, totpToken, next) {

  Users.select({ email: email }, function (err, user) {

    if (err !== undefined) {
      return next(err);
    }

    if (user === undefined) {
      return next(undefined, false, 'No account with that email exists');
    }

    if (user.use_totp === 1) {
      if (!Notp.totp.verify(totpToken, Base32.decode(user.totp_secret), { time: 30 })) {
        return next(undefined, false, 'TOTP token was invalid');
      }
    }

    if (!Bcrypt.compareSync(password, user.password)) {
      return next(undefined, false, 'Email or password was incorrect');
    }

    return next(undefined, true, user);

  });

};


/**
 * Gets if a given user's TOTP option is enabled.
 *
 * @param {String} email - Email of user
 * @param {callback} next - Callback to handle response
 */
Authentication.isTOTPEnabled = function (email, next) {

  Users.select({ email: email }, function (err, user) {

    return next(err, user !== undefined && user.use_totp === 1);

  });

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
