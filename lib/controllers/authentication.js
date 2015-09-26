/*
 * PufferPanel - Reinventing the way game servers are managed.
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
var UserModels = (process.env.NODE_ENV === 'test') ? Rfr('tests/models/users.js') : Rfr('lib/models/users.js');
var RandomString = require('randomstring');

/** @namespace */
var AuthenticationController = {};


//TODO: Very likely this should be in UserController instead of here, as it's a user function
//Or at least make a proxy function
/**
 * Determines whether given credentials are valid.
 *
 * @param {String} email - User's email
 * @param {String} password - User's password
 * @param {String} totpToken - Request's TOTP token (may be undefined)
 * @param {Authentication~loginUserCb} next - Callback to handle response
 */
AuthenticationController.loginUser = function (email, password, totpToken, ipaddress, next) {

    UserModels.select({ email: email }, function (err, user) {

        if (err) {
            return next(err);
        }

        if (!user) {
            return next(null, 'No account with that email exists');
        }

        if (user.useTotp) {
            if (!Notp.totp.verify(totpToken, Base32.decode(user.totpSecret), { time: 30 })) {
                return next(null, 'TOTP token was invalid');
            }
        }

        if (!Bcrypt.compareSync(password, user.password)) {
            return next(null, 'Email or password was incorrect');
        }

        var sessionTokenValue = RandomString.generate(15);
        user.sessionToken = sessionTokenValue;
        user.sessionIp = ipaddress;

        UserModels.update(user.id, { sessionToken: sessionTokenValue, sessionIp: ipaddress }, function (err) {

            if (err) {

                return next(err);
            }

            return next(null, user);
        });
    });
};
/**
 * @callback AuthenticationController~loginUserCb
 * @param {Error} err - Error if one occurred, otherwise null
 * @param {Object|String} data - The User who is now logged on, otherwise a String with the failure reason
 */

AuthenticationController.validateAccountPassword = function (id, password, next) {

    UserModels.select({ id: id }, function (err, user) {

        if (err) {
            return next(err);
        }

        if (!Bcrypt.compareSync(password, user.password)) {
            return next(false);
        }

        return next(null);

    });

};

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
