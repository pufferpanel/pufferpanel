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
var Authentication = Rfr('lib/controllers/authentication.js');
var User = Rfr('lib/controllers/user.js');
var Logger = Rfr('lib/logger.js');

var Routes = {
    post: {
        login: function (request, reply) {

            Authentication.loginUser(request.payload.email, request.payload.password, request.payload.totpToken, function (err, data) {

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
        },
        totp: function (request, reply) {

            User.isTOTPEnabled(request.payload.check, function (err, data) {

                if (err) {
                    Logger.error(err);
                }

                reply(data.toString());

            });

        },
        register: function () {
        },
        password: function () {
        }
    },
    get: {
        login: function (request, reply) {

            reply.view('auth/login', {
                flash: request.session.flash('loginError')
            });

        },
        logout: function (request, reply) {

            request.auth.session.clear();
            reply.redirect('/auth/login');

        },
        register: function (request, reply) {

            reply.view('auth/register', {
                flash: request.session.flash('registerError'),
                token: request.params.token
            });

        },
        password: function (request, reply) {

            reply.view('auth/password', {
                flash: request.session.flash('passwordError'),
                noshow: false
            });

        }
    }
};

module.exports = Routes;
