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

/** @namespace */
var AuthenticationController = {};


/**
 * Determines whether given credentials are valid.
 *
 * @param {String} email - User's email
 * @param {String} password - User's password
 * @param {String} totpToken - Request's TOTP token (may be undefined)
 * @param {Authentication~loginUserCb} next - Callback to handle response
 */
AuthenticationController.loginUser = function (email, password, totpToken, next) {

    Users.select({ email: email }, function (err, user) {

        if (err !== undefined) {
            return next(err);
        }

        if (user === undefined) {
            return next(null, 'No account with that email exists');
        }

        if (user.use_totp === 1) {
            if (!Notp.totp.verify(totpToken, Base32.decode(user.totp_secret), { time: 30 })) {
                return next(null, 'TOTP token was invalid');
            }
        }

        if (!Bcrypt.compareSync(password, user.password)) {
            return next(null, 'Email or password was incorrect');
        }

        return next(null, user);
    });
};
/**
 * @callback AuthenticationController~loginUserCb
 * @param {Error} err - Error if one occurred, otherwise null
 * @param {Object|String} data - The User who is now logged on, otherwise a String with the failure reason
 */


/**
 * Gets if a given user's TOTP option is enabled.
 * If the user does not exist, then the callback will receive false.
 *
 * @param {String} email - Email of user
 * @param {Authentication~isTotpEnabledCb} next - Callback to handle response
 */
AuthenticationController.isTOTPEnabled = function (email, next) {

    Users.select({ email: email }, function (err, user) {

        return next(err || null, user !== undefined && user.use_totp === 1);
    });
};
/**
 * @callback AuthenticationController~isTotpEnabledCb
 * @param {Error} err - Error if one occurred, otherwise null
 * @param {Boolean} data - Whether or not the user's TOTP is enabled, or false if the user does not exist
 */


/**
 * Generates a {@link bcrypt}-hashed password
 *
 * @param {String} rawpassword - Password to hash
 * @returns {String} Hashed form of the password
 */
AuthenticationController.generatePasswordHash = function (rawpassword) {

    return Bcrypt.hashSync(rawpassword, 10);
};

module.exports = AuthenticationController;
