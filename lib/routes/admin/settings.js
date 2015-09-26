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
var Util = require('util');
var AdminSettingsController = Rfr('lib/controllers/admin/settings.js');
var Routes = {};

Routes.router = function (request, reply) {

    if (request.params.action === 'general') {
        Routes.getGeneralSettings(request, reply);
    }

    if (request.params.action === 'urls') {
        //Routes.settings.urls(request, reply);
    }

    if (request.params.action === 'email') {
        //Routes.get.settings.email(request, reply);
    }

    if (request.params.action === 'captcha') {
        //Routes.get.settings.captcha(request, reply);
    }

    if (reply._replied === false) {
        return reply(Boom.notFound());
    }

};

Routes.getGeneralSettings = function (request, reply) {

    return reply.view('admin/settings/general', {
        flashSuccess: request.session.flash('generalSuccess'),
        flashFailure: request.session.flash('generalFailure')
    });

};

Routes.postUpdateCompanyName = function (request, reply) {

    AdminSettingsController.updateCompanyName(request.payload.companyName, function (err, response) {

        if (err) {
            Logger.error('An error occured attempting to update the company name.', err);
            return reply(Boom.badImplementation());
        }

        if (!response) {
            request.session.flash('generalSuccess', 'The company name for this instance has been updated.');
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

};

Routes.postGeneralSettings = function (request, reply) {

    request.payload.enableSubusers = !(typeof request.payload.enableSubusers === 'undefined');
    request.payload.requireSecure = !(typeof request.payload.requireSecure === 'undefined');

    AdminSettingsController.updateGeneralSettings(request.payload.requireSecure, request.payload.enableSubusers, function (err, response) {

        if (err) {
            Logger.error('An error occured attempting to update the company name.', err);
            return reply(Boom.badImplementation());
        }

        if (!response) {
            request.session.flash('generalSuccess', 'The general settings for this instance have been updated.');
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

};

module.exports = Routes;
