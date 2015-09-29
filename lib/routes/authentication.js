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
var Logger = Rfr('lib/logger.js');
var Boom = require('boom');
var Recaptcha = require('nodejs-nocaptcha-recaptcha');
var SettingsController = Rfr('lib/controllers/settings.js');
var AuthenticationController = Rfr('lib/controllers/authentication.js');
var UserController = Rfr('lib/controllers/user.js');
var UserVisibleError = Rfr('lib/errors/UserVisibleError.js');

var Routes = {};

Routes.postPassword = function (request, reply) {

    SettingsController.get('captcha', function (err, data) {

        if (err) {
            return reply(Boom.badImplementation());
        }

        Recaptcha(request.payload['g-recaptcha-response'], data.secret, function (success) {

            if (success) {

                AuthenticationController.generatePasswordReset(request.payload.email, request.info.remoteAddress, function (err) {

                    if (err && !(err instanceof UserVisibleError)) {
                        Logger.error('An error occured while attempting to perform a password reset.', err);
                        return reply(Boom.badImplementation());
                    }

                    if (!err) {
                        request.session.flash('passwordSuccess', 'An email has been sent to this email address requesting a password reset.');
                    } else {
                        if (err.messageIsString === false) {
                            _.each(err.message, function (value) {
                                request.session.flash('passwordError', value);
                            });
                        } else {
                            request.session.flash('passwordError', err.message);
                        }
                    }

                    return reply.redirect('/auth/password');

                });

            } else {
                request.session.flash('passwordError', 'The captcha was not filled out correctly.');
                return reply.redirect('/auth/password');
            }

        });

    });

};

Routes.postLogin = function (request, reply) {

    AuthenticationController.loginUser(request.payload.email, request.payload.password, request.payload.totpToken, request.info.remoteAddress, function (err, data) {

        if (err || typeof data === 'string') {

            if (err) {
                Logger.error(err);
            }

            request.session.flash('loginError', data);
            reply.redirect('/auth/login');

        } else {

            if (typeof data.language === undefined) {
                data.language = 'en';
            }
            reply.state('pp_language', data.language, {
                ttl: 60 * 60 * 24 * 1000,
                path: '/',
                isSecure: false,
                isHttpOnly: true
            });

            request.auth.session.set(data);
            reply.redirect('/');

        }
    });

};

Routes.postRegister = function (request, reply) {

    reply(Boom.notFound());

};

Routes.postTotp = function (request, reply) {

    UserController.isTOTPEnabled(request.payload.check, function (err, data) {

        if (err) {
            Logger.error(err);
        }

        reply(data.toString());

    });

};

Routes.getLogin = function (request, reply) {

    reply.view('auth/login', {
        flash: request.session.flash('loginError')
    });

};

Routes.getLogout = function (request, reply) {

    request.auth.session.clear();
    reply.redirect('/auth/login');

};

Routes.getRegister = function (request, reply) {

    reply.view('auth/register', {
        flash: request.session.flash('registerError'),
        token: request.params.token
    });

};

Routes.getPassword = function (request, reply) {

    reply.view('auth/password', {
        flashFailure: request.session.flash('passwordError'),
        flashSuccess: request.session.flash('passwordSuccess')
    });

};

module.exports = Routes;
