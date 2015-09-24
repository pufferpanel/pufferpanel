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
var Locations = Rfr('lib/controllers/locations.js');
var Users = Rfr('lib/controllers/user.js');
var Async = require('async');

var Routes = {
    get: {
        servers: function (request, reply) {

            Async.parallel({
                servers: function (next) {
                    Servers.getServersFor(request.auth.credentials.id, function (err, servers) {

                        if (err) {
                            return next(err);
                        }

                        return next(null, servers || []);
                    });
                },
                locations: function (next) {
                    Locations.getAllLocations(function (err, data) {

                        if (err) {
                            return next(err);
                        }

                        return next(null, data);
                    });
                }
            }, function (err, results) {

                if (err) {
                    Logger.error(err);
                    return reply.view('code/500.html');
                }

                return reply.view('base/servers', {
                    servers: results.servers,
                    locations: results.locations,
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

                Users.generateTOTP(request.auth.credentials.id, function (err, resp) {

                    if (err) {
                        Logger.error(err);
                        request.session.flash('totpError', err);
                    }

                    request.auth.session.set('totp_secret', resp.secret);
                    reply.view('base/totp-popup', {
                        flash: request.session.flash('totpError'),
                        totp: resp
                    });
                });

            },
            verifyToken: function (request, reply) {

                Users.toggleTOTP(request.auth.credentials.id, request.payload.token, request.auth.credentials.totpSecret, function (err) {

                    //TODO: relook at this, because rethinkdb errors can flow up and get shown to client here
                    if (err) {
                        Logger.error(err);
                        request.session.flash('totpError', err.toString());
                    }

                    //TODO: Probably a better way to do this.
                    if (request.payload.ajax) {
                        reply('hodor');
                    } else {
                        reply.redirect('/totp');
                    }

                });

            },
            disableToken: function (request, reply) {


            }
        }
    }
};

module.exports = Routes;
