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
var Logger = Rfr('lib/logger.js');
var Fs = require('fs-extra');
var Boom = require('boom');
var Util = require('util');
var Routes = {
    get: {
        index: function (request, reply) {

            Fs.readFile('./.git/HEAD', function (err, data) {

                if (err) {
                    Logger.error('An error occured while attempting to process .git information.', err);
                    return Boom.badImplementation();
                }

                if (data.indexOf('ref: ') > -1) {

                    var ref = data.toString().split(' ');
                    Fs.readFile(Util.format('./.git/%s', ref[1].trim()), function (err, moreData) {

                        if (err) {
                            Logger.error('An error occured while attempting to process more detailed .git information.', err);
                            return Boom.badImplementation();
                        }

                        return reply.view('admin/index', {
                            version: ref[1].trim(),
                            sha: moreData.toString().trim()
                        });

                    });

                } else {

                    return reply.view('admin/index', {
                        version: 'master',
                        sha: data.toString().trim()
                    });

                }

            });

        },
        settings: {
            router: function (request, reply) {

                if (request.params.action === 'global') {
                    Routes.get.settings.general(request, reply);
                }

                if (request.params.action === 'urls') {
                    Routes.get.settings.urls(request, reply);
                }

                if (request.params.action === 'email') {
                    Routes.get.settings.email(request, reply);
                }

                if (request.params.action === 'captcha') {
                    Routes.get.settings.captcha(request, reply);
                }

                return Boom.notFound();

            },
            general: function (request, reply) {

                return reply.view('admin/settings/general', {
                    flashSuccess: request.session.flash('generalSuccess'),
                    flashFailure: request.session.flash('generalFailure')
                });

            }
        }
    },
    post: {
        settings: {
            general: {
                updateCompanyName: function (request, reply) {

                    AdminSettingsController.updateCompanyName(request.payload.companyName, function (err, response) {

                        if (err) {
                            Logger.error('An error occured attempting to update the company name.', err);
                            return Boom.badImplementation();
                        }

                        if (!response) {
                            request.session.flash('generalSuccess', 'The comapny name for this instance has been updated.');
                        } else {

                            if (typeof response !== 'string') {
                                _.each(response, function (value) {
                                    request.session.flash('generalFailure', value);
                                });
                            } else {
                                request.session.flash('generalFailure', response);
                            }

                        }

                        return reply.redirect('/admin/settings/general');
                    });

                },
                generalSettings: function (request, reply) {

                }
            }
        }
    }
};

module.exports = Routes;
