/*
 * PufferPanel - Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var _ = require('underscore');
var Rfr = require('rfr');
var Base32 = require('thirty-two');
var Crypto = require('crypto');
var Notp = require('notp');
var Joi = require('joi');
var AuthenticationController = Rfr('lib/controllers/authentication.js');
var UserModels = (process.env.NODE_ENV === 'test') ? Rfr('tests/models/users.js') : Rfr('lib/models/users.js');

/** @namespace */
var UserController = {};

UserController.getData = function (id, next) {

    UserModels.select({ id: id }, function (err, user) {

        return next(err || null, user || false);
    });

};

UserController.updatePassword = function (id, oldPassword, newPassword, next) {

    var schema = {
        np: Joi.string().regex(/((?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,})/).required().label('New Password'),
        op: Joi.string().required().label('Current Password')
    };
    var errorResponse = [];

    // Validate Data
    Joi.validate({ np: newPassword, op: oldPassword }, schema, { abortEarly: false }, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(null, errorResponse);
        }

        AuthenticationController.validateAccountPassword(id, oldPassword, function (err) {

            if (err && err !== false) {
                return next(err);
            }

            if (err === false) {
                return next(null, 'The password provided was invalid for this account.');
            }

            UserModels.update(id, { password: AuthenticationController.generatePasswordHash(newPassword) }, function (err) {
                return next(err);
            });

        });

    });

};

UserController.updateNotifications = function (id, password, notificationSuccess, notificationFailure, next) {

    var schema = {
        pass: Joi.string().required().label('Password'),
        ns: Joi.boolean().label('Notification on Success'),
        nf: Joi.boolean().label('Notification on Failure')
    };
    var errorResponse = [];

    // Validate Data
    Joi.validate({ pass: password, ns: notificationSuccess, nf: notificationFailure }, schema, { abortEarly: false }, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(null, errorResponse);
        }

        AuthenticationController.validateAccountPassword(id, password, function (err) {

            if (err && err !== false) {
                return next(err);
            }

            if (err === false) {
                return next(null, 'The password provided was invalid for this account.');
            }

            // Make String into a Bool
            notificationFailure = (notificationFailure === 'true');
            notificationSuccess = (notificationSuccess === 'true');

            UserModels.update(id, { notifications: { loginFailure: notificationFailure, loginSuccess: notificationSuccess } }, function (err) {
                return next(err);
            });

        });

    });

};

UserController.updateEmail = function (id, password, newEmail, next) {

    var schema = {
        pass: Joi.string().required().label('Password'),
        email: Joi.string().email().label('New Email')
    };
    var errorResponse = [];

    Joi.validate({ pass: password, email: newEmail }, schema, { abortEarly: false }, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(null, errorResponse);
        }

        AuthenticationController.validateAccountPassword(id, password, function (err) {

            if (err && err !== false) {
                return next(err);
            }

            if (err === false) {
                return next(null, 'The password provided was invalid for this account.');
            }

            UserModels.update(id, { email: newEmail }, function (err) {
                return next(err);
            });

        });

    });

};

/**
 * Generates TOTP details for a user and saves the secret to the database.
 * Does not enable it, only generates the secret token and the QR code image.
 * @param {String} id - User Id
 * @param {UserController~generateTotpCb} next - Callback to handle response
 */
UserController.generateTOTP = function (id, next) {

    var secret = Base32.encode(Crypto.randomBytes(16)).toString().replace(/=/g, '');
    var totp = {
        secret: secret,
        image: 'https://chart.googleapis.com/chart?chs=166x166&chld=L|0&cht=qr&chl=otpauth://totp/PufferPanel?secret=' + secret
    };

    UserModels.update(id, { totpSecret: totp.secret }, function (err) {

        if (err) {
            return next(err);
        }

        return next(null, totp);
    });
};
/**
 * @callback UserController~generateTotpCb
 * @param {Error} err - Error if one occurred, otherwise null
 * @param {Object} totp - TOTP response with 2 properties: secret and image
 */

/**
 * Gets if a given user's TOTP option is enabled.
 * If the user does not exist, then the callback will receive false.
 *
 * @param {String} email - Email of user
 * @param {UserController~isTotpEnabledCb} next - Callback to handle response
 */
UserController.isTOTPEnabled = function (email, next) {

    UserModels.select({ email: email }, function (err, user) {

        return next(err || null, user !== null && user !== undefined && user.useTotp === true);
    });
};
/**
 * @callback UserController~isTotpEnabledCb
 * @param {Error} err - Error if one occurred, otherwise null
 * @param {Boolean} data - Whether or not the user's TOTP is enabled, or false if the user does not exist
 */


/**
 * Validates request to enable TOTP.
 * @param {String} id - User id
 * @param {String} token - TOTP token
 * @param {String} secret - TOTP secret
 * @param {UserController~enableTotpCb} next - Callback to handle response
 */
//TODO: Possibly move TOTP checks to another lib for unit testing
UserController.toggleTOTP = function (id, token, secret, next) {

    if (!Notp.totp.verify(token, Base32.decode(secret), { time: 30 })) {
        return next(new Error('Unable to validate the TOTP token.'));
    }

    UserModels.select({ id: id }, function (err, user) {

        if (err) {
            return next(err);
        }

        user.useTotp = !user.useTotp;
        UserModels.update(user.id, { useTotp: user.useTotp }, function (err) {

            if (err) {
                return next(err);
            }

            return next(null, user);
        });

    });

};
/**
 * @callback UserController~enableTotpCb
 * @param {Error} err - Error if one occurred, otherwise null
 * @param {Object} user - User involved with TOTP updated
 */

module.exports = UserController;
