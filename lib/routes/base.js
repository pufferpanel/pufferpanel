/*
 * PufferPanel — Reinventing the way game servers are managed.
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
        language: function () {
            // Handle setting language here
        },
        totp: function (request, reply) {

            reply.view('base/totp', {
                flash: request.session.flash('totpError'),
                user: request.auth.credentials
            });
        },
        account: function (request, reply) {

            reply.view('base/account', {
                flashError: request.session.flash('accountError'),
                flashSuccess: request.session.flash('accountSuccess'),
                user: request.auth.credentials
            });
        }
    },
    post: {
        accountUpdate: function (request, reply) {

            if (request.params.action === 'password') {

                if (request.payload.newPassword !== request.payload.newPasswordAgain) {
                    request.session.flash('accountError', 'An error occured while attempting to update the account password. The new passwords did not match.');
                }

                Users.updatePassword(request.auth.credentials.id, request.payload.currentPassword, request.payload.newPassword, function (err, response) {

                    if (err) {
                        return reply.view('code/500').code(500);
                    }

                    if (!response) {
                        request.session.flash('accountSuccess', 'Your password has been successfully updated.');
                    } else {

                        if (typeof response !== 'string') {
                            _.each(response, function (value) {
                                request.session.flash('accountError', value);
                            });
                        } else {
                            request.session.flash('accountError', response);
                        }

                    }

                    return reply.redirect('/account');

                });
            }

            if (request.params.action === 'email') {

                // Response is only set if there is an error thrown by Joi
                Users.updateEmail(request.auth.credentials.id, request.payload.currentPassword, request.payload.newEmail, function (err, response) {

                    if (err) {
                        Logger.error(err);
                        return reply.view('code/500').code(500);
                    }

                    if (!response) {
                        request.session.flash('accountSuccess', 'Your email has been successfully updated.');
                    } else {

                        if (typeof response !== 'string') {
                            _.each(response, function (value) {
                                request.session.flash('accountError', value);
                            });
                        } else {
                            request.session.flash('accountError', response);
                        }

                    }

                    return reply.redirect('/account');

                });

            }

            if (request.params.action === 'notifications') {

                Users.updateNotifications(request.auth.credentials.id, request.payload.currentPassword, request.payload.loginSuccess, request.payload.loginFailure, function (err, response) {

                    if (err) {
                        Logger.error(err);
                        return reply.view('code/500').code(500);
                    }

                    if (!response) {
                        request.session.flash('accountSuccess', 'Your notification preferences have been successfully updated.');
                    } else {

                        if (typeof response !== 'string') {
                            _.each(response, function (value) {
                                request.session.flash('accountError', value);
                            });
                        } else {
                            request.session.flash('accountError', response);
                        }

                    }

                    return reply.redirect('/account');

                });

            }

            // return reply.view('code/500').code(500);

        },
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

            }
        }
    }
};

module.exports = Routes;
