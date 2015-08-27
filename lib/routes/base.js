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
var Base32 = require('thirty-two');
var Logger = Rfr('lib/logger.js');
var Servers = Rfr('lib/controllers/server.js');
var Users = Rfr('lib/controllers/user.js');

var Routes = {
    get: {
        servers: function (request, reply) {

            Servers.getServersFor(request.auth.credentials.id, function (err, servers) {

                if (err !== undefined) {

                    Logger.error(err);
                }

                reply.view('base/index', {
                    servers: servers || [],
                    user: request.auth.credentials
                });

            });
        },
        language: function (request, reply) {

            // Handle setting language here
        },
        totp: function (request, reply) {

            reply.view('base/totp', {
                flash: request.session.flash('totpError'),
                user: request.auth.credentials
            });
        }
    },
    post: {
        totp: {
            generateToken: function (request, reply) {

                Users.generateTOTP(request, function (err, resp) {

                    if (err) {
                        Logger.error(err);
                        request.session.flash('totpError', err);
                    }

                    reply.view('base/totp-popup', {
                        flash: request.session.flash('totpError'),
                        totp: resp
                    });
                });

            },
            verifyToken: function (request, reply) {

                Users.enableTotp(request, function (err) {

                    if (err) {

                        request.session.flash('totpError', err.toString());
                    }

                    // Probably a better way to do this.
                    reply('hodor');
                });

            }
        }
    }
};

module.exports = Routes;
