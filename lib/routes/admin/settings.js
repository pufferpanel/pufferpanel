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
var UserVisibleError = Rfr('lib/errors/UserVisibleError.js');
var Routes = {};

Routes.router = function (request, reply) {

    if (request.params.action === 'general') {
        if (request.method === 'post') {
            return Routes.postGeneralRouter(request, reply);
        }
        return Routes.getGeneralSettings(request, reply);
    }

    if (request.params.action === 'urls') {
        if (request.method === 'post') {
            return Routes.postUpdateUrlSettings(request, reply);
        }
        return Routes.getUrlSettings(request, reply);
    }

    if (request.params.action === 'email') {
        if (request.method === 'post') {
            return Routes.postUpdateEmailSettings(request, reply);
        }
        return Routes.getEmailSettings(request, reply);
    }

    if (request.params.action === 'captcha') {
        //Routes.get.settings.captcha(request, reply);
    }

    if (reply._replied === false) {
        return reply(Boom.notFound());
    }

};

Routes.postGeneralRouter = function (request, reply) {

    if (request.payload.doAction === 'general') {
        return Routes.postUpdateGeneralSettings(request, reply);
    }

    if (request.payload.doAction === 'company') {
        return Routes.postUpdateCompanyName(request, reply);
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

Routes.getUrlSettings = function (request, reply) {

    return reply.view('admin/settings/urls', {
        flashSuccess: request.session.flash('urlSuccess'),
        flashFailure: request.session.flash('urlFailure')
    });

};


Routes.getEmailSettings = function (request, reply) {

    return reply.view('admin/settings/email', {
        flashSuccess: request.session.flash('emailSuccess'),
        flashFailure: request.session.flash('emailFailure')
    });

};

Routes.postUpdateUrlSettings = function (request, reply) {

    AdminSettingsController.updateUrls(request.payload.mainUrl, request.payload.assetsUrl, function (err) {

        // General Error; Errors returned by the controller are an instance of UserVisibleError and can be displayed to users.
        if (err && !err instanceof UserVisibleError) {
            Logger.error('An error occured while attempting to update URL settings.', err);
            return reply(Boom.badImplementation());
        }

        if (!err) {
            request.session.flash('urlSuccess', 'URLs have been updated for this instance.');
        } else {
            if (err.messageIsString === false) {
                _.each(err.message, function (value) {
                    request.session.flash('urlFailure', value);
                });
            } else {
                request.session.flash('urlFailure', err.message);
            }
        }

        return reply.redirect('/admin/settings/urls');

    });

};

Routes.postUpdateCompanyName = function (request, reply) {

    AdminSettingsController.updateCompanyName(request.payload.companyName, function (err, response) {

        if (err && !err instanceof UserVisibleError) {
            Logger.error('An error occured attempting to update the company name.', err);
            return reply(Boom.badImplementation());
        }

        if (!err) {
            request.session.flash('generalSuccess', 'The company name for this instance has been updated.');
        } else {
            if (err.messageIsString === false) {
                _.each(err.message, function (value) {
                    request.session.flash('generalFailure', value);
                });
            } else {
                request.session.flash('generalFailure', err.message);
            }
        }

        return reply.redirect('/admin/settings/general');
    });

};

Routes.postUpdateGeneralSettings = function (request, reply) {

    request.payload.enableSubusers = !(typeof request.payload.enableSubusers === 'undefined');
    request.payload.requireSecure = !(typeof request.payload.requireSecure === 'undefined');

    AdminSettingsController.updateGeneralSettings(request.payload.requireSecure, request.payload.enableSubusers, function (err, response) {

        if (err && !err instanceof UserVisibleError) {
            Logger.error('An error occured attempting to update the company name.', err);
            return reply(Boom.badImplementation());
        }

        if (!err) {
            request.session.flash('generalSuccess', 'The general settings for this instance have been updated.');
        } else {
            if (err.messageIsString === false) {
                _.each(err.message, function (value) {
                    request.session.flash('generalFailure', value);
                });
            } else {
                request.session.flash('generalFailure', err.message);
            }
        }

        return reply.redirect('/admin/settings/general');
    });

};

module.exports = Routes;
