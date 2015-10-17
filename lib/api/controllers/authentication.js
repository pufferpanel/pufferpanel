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
var Bcrypt = require('bcrypt');
var Notp = require('notp');
var Base32 = require('thirty-two');
var RandomString = require('randomstring');
var Crypto = require('crypto');
var Moment = require('moment');
var Joi = require('joi');
var EmailController = Rfr('lib/api/controllers/email.js');
var UserModel = (process.env.NODE_ENV === 'test') ? Rfr('tests/models/users.js') : Rfr('lib/api/models/users.js');
var ActionsModel = (process.env.NODE_ENV === 'test') ? Rfr('tests/models/actions.js') : Rfr('lib/api/models/actions.js');
var UserVisibleError = Rfr('lib/errors/UserVisibleError.js');

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

    UserModel.select({ email: email }, function (err, user) {

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

        UserModel.update(user.id, { sessionToken: sessionTokenValue, sessionIp: ipaddress }, function (err) {

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

    UserModel.select({ id: id }, function (err, user) {

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

AuthenticationController.encrypt = function (string) {

    if (!process.env.CIPHER_KEY) {
        throw new Error('Unable to encrypt string due to a missing cipher key! In order to encrypt items you must start this process with \'CIPHER_KEY\' environment variable.');
    }

    if (!process.env.CIPHER_ALGO) {
        process.env.CIPHER_ALGO = 'AES-256-CTR';
    }

    var cipher = Crypto.createCipher(process.env.CIPHER_ALGO, process.env.CIPHER_KEY);
    var crypted = cipher.update(string, 'utf8', 'hex');

    crypted += cipher.final('hex');
    return crypted;

};

AuthenticationController.decrypt = function (string) {

    if (!process.env.CIPHER_KEY) {
        throw new Error('Unable to encrypt string due to a missing cipher key! In order to encrypt items you must start this process with \'CIPHER_KEY\' environment variable.');
    }

    if (!process.env.CIPHER_ALGO) {
        process.env.CIPHER_ALGO = 'AES-256-CTR';
    }

    var decipher = Crypto.createDecipher(process.env.CIPHER_ALGO, process.env.CIPHER_KEY);
    var deciphered = decipher.update(string, 'hex', 'utf8');

    deciphered += decipher.final('utf8');
    return deciphered;

};

AuthenticationController.generatePasswordReset = function (email, ip, next) {

    var schema = {
        email: Joi.string().email().required().label('Email')
    };
    var errorResponse = [];

    Joi.validate({ email: email }, schema, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(new UserVisibleError(errorResponse));
        }

        UserModel.select({ email: email }, function (err, user) {

            if (err) {
                return next(err);
            }

            if (!user) {
                return next(new UserVisibleError('No account with that email exists.'));
            }

            var token = RandomString.generate(36);
            ActionsModel.create({
                email: email,
                user: user.id,
                token: token,
                ip: ip,
                time: Moment().toString()
            }, function (err, documentToken) {

                if (err) {
                    return next(err);
                }

                // Send the Email
                var templateData = {
                    subject: 'Account Password Reset Request',
                    token: token,
                    ip: ip,
                    time: Moment().format()
                };

                EmailController.sendTemplate('resetPassword', email, templateData, function (err) {

                    if (err) {

                        var emailError = err;
                        ActionsModel.deleteId(documentToken, function (err) {

                            if (err) {
                                return next({
                                    EmailController: emailError,
                                    ActionsModel: err
                                });
                            }

                            return next(emailError);

                        });

                    }

                    return next();

                });

            });

        });

    });

};

AuthenticationController.validatePasswordReset = function (token, next) {

    var schema = {
        token: Joi.string().length(36).token().required().label('Verification Token')
    };

    Joi.validate({ token: token }, schema, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(new UserVisibleError(errorResponse));
        }

        ActionsModel.select({
            token: token
        }, function (err, data) {

            if (err) {
                return next(err);
            }

            if (!data) {
                return next(new UserVisibleError('The token provided was invalid.'));
            }

            var time = {
                current: Moment(),
                former: Moment(Moment(new Date(data.time)).format())
            };

            if (time.current.isAfter(time.former.add(4, 'hours'))) {
                return next(new UserVisibleError('This token is over 4 hours old and has expired.'));
            }

            // Email Template Data
            var templateData = {
                subject: 'Account Password Successfully Reset',
                newPassword: RandomString.generate(20),
                time: Moment().toString()
            };

            UserModel.update(data.user, { password: AuthenticationController.generatePasswordHash(templateData.newPassword) }, function (err) {

                if (err) {
                    return next(err);
                }

                ActionsModel.deleteId(data.id, function (err) {

                    if (err) {
                        return next(err);
                    }

                    EmailController.sendTemplate('passwordReset', data.email, templateData, function (err) {
                        return next(err);
                    });

                });

            });

        });

    });

};

module.exports = AuthenticationController;
