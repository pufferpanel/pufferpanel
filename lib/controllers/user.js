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

var User = {};

/**
 * Generates TOTP details for a user and saves the secret to the database.
 * Does not enable it, only generates the secret token and the QR code image.
 * @param  {Object} request
 * @param  {Function} next Callback function
 * @return {Object}        The image and totp secret.
 */
User.generateTOTP = function (request, next) {

    var tempSecret = Base32.encode(Crypto.randomBytes(16)).toString().replace(/=/g, '');
    var totp = {
        secret: tempSecret,
        image: 'https://chart.googleapis.com/chart?chs=166x166&chld=L|0&cht=qr&chl=otpauth://totp/PufferPanel?secret=' + tempSecret
    };

    UserModels.update(request.auth.credentials.id, { totp_secret: totp.secret }, function (err) {

        if (err) {

            Logger.error('An error occured while attempting to update user information.', err);
            return next(err);
        }

        request.auth.session.set('totp_secret', totp.secret);
        return next(null, totp);
    });

};

/**
 * Validates request to enable TOTP.
 * @param  {Object}   user
 * @param  {Function} next
 * @return {[type]}
 */
User.enableTotp = function (request, next) {

    if (!Notp.totp.verify(request.payload.token, Base32.decode(request.auth.credentials.totp_secret), { time: 30 })) {
        return next(new Error('Unable to validate the TOTP token.'));
    }

    UserModels.update(request.auth.credentials.id, { use_totp: 1 }, function (err) {

        if (err) {

            Logger.error('An error occured while attempting to update user information.', err);
            return next(err);
        }

        request.auth.session.set('use_totp', 1);
        return next(null);

    });

};

module.exports = User;
