/*
 * PufferPanel � Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Base32 = require('thirty-two');
var Crypto = require('crypto');
var Notp = require('notp');
var UserModels = Rfr('lib/models/users.js');

/** @namespace */
var UserController = {};

/**
 * Generates TOTP details for a user and saves the secret to the database.
 * Does not enable it, only generates the secret token and the QR code image.
 * @param {Object} request - Hapi request object
 * @param {UserController~generateTotpCb} next - Callback to handle response
 */
//TODO: Remove request here, pass just what we need
UserController.generateTOTP = function (request, next) {

    var tempSecret = Base32.encode(Crypto.randomBytes(16)).toString().replace(/=/g, '');
    var totp = {
        secret: tempSecret,
        image: 'https://chart.googleapis.com/chart?chs=166x166&chld=L|0&cht=qr&chl=otpauth://totp/PufferPanel?secret=' + tempSecret
    };

    UserModels.update(request.auth.credentials.id, { totp_secret: totp.secret }, function (err) {

        if (err) {
            return next(err);
        }

        //TODO: This should occur outside this controller
        request.auth.session.set('totp_secret', totp.secret);
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

        return next(err || null, user !== undefined && user.use_totp === 1);
    });
};
/**
 * @callback UserController~isTotpEnabledCb
 * @param {Error} err - Error if one occurred, otherwise null
 * @param {Boolean} data - Whether or not the user's TOTP is enabled, or false if the user does not exist
 */


/**
 * Validates request to enable TOTP.
 * @param  {Object} request - Hapi request object
 * @param  {UserController~enableTotpCb} next - Callback to handle response
 */
//TODO: Also remove request usage here
UserController.enableTotp = function (request, next) {

    if (!Notp.totp.verify(request.payload.token, Base32.decode(request.auth.credentials.totp_secret), { time: 30 })) {
        return next(new Error('Unable to validate the TOTP token.'));
    }

    UserModels.update(request.auth.credentials.id, { use_totp: 1 }, function (err) {

        if (err) {
            return next(err);
        }

        request.auth.session.set('use_totp', 1);
        return next(null);
    });
};
/**
 * @callback UserController~enableTotpCb
 * @param {Error} err - Error if one occurred, otherwise null
 */

module.exports = UserController;
