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
var SettingsModel = Rfr('lib/models/settings.js');
var Joi = require('joi');
var UserVisibleError = Rfr('lib/errors/USerVisibleError.js');
var SettingsController = Rfr('lib/controllers/settings.js');
var URL = require('url-parse');

/** @namespace */
var AdminSettingsController = {};

AdminSettingsController.updateUrls = function (main, assets, next) {

    var schema = {
        main: Joi.string().required().label('Main URL'),
        assets: Joi.string().required().label('Assets URL')
    };
    var errorResponse = [];

    Joi.validate({ main: main, assets: assets }, schema, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(new UserVisibleError(errorResponse));
        }

        SettingsController.get('useSecureConnection', function (err, value) {

            if (err) {
                return next(err);
            }

            var useSecureConnection = value;

            var parseMain = URL(main, false);
            var parseAssets = URL(assets, false);

            parseMain.set('protocol', (useSecureConnection === true) ? 'https:' : 'http:');
            parseAssets.set('protocol', (useSecureConnection === true) ? 'https:' : 'http:');

            var setMain = parseMain.href;
            var setAssets = parseAssets.href;

            if (setMain.substr(-1) === '/') {
                setMain = setMain.substr(0, setMain.length - 1);
            }

            if (setAssets.substr(-1) === '/') {
                setAssets = setAssets.substr(0, setAssets.length - 1);
            }

            SettingsModel.update({
                'urls': {
                    'assets': setAssets,
                    'main': setMain
                }
            }, function (err) {
                return next(err);
            });

        });

    });

};

AdminSettingsController.updateCompanyName = function (name, next) {

    var schema = {
        companyName: Joi.string().min(3).max(100).required().label('Company Name')
    };
    var errorResponse = [];

    Joi.validate({ companyName: name }, schema, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(new UserVisibleError(errorResponse));
        }

        SettingsModel.update({ 'companyName': name }, function (err) {
            return next(err);
        });

    });

};

AdminSettingsController.updateGeneralSettings = function (secureConnection, subUsers, next) {

    var schema = {
        sc: Joi.boolean().required(),
        su: Joi.boolean().required()
    };
    var errorResponse = [];

    Joi.validate({ sc: secureConnection, su: subUsers }, schema, function (err) {

        if (err) {
            _.each(err.details, function (v) {
                errorResponse.push(v.message);
            });
            return next(new UserVisibleError(errorResponse));
        }

        SettingsModel.update({
            'useSecureConnection': secureConnection,
            'enableSubusers': subUsers
        }, function (err) {
            return next(err);
        });

    });

};

module.exports = AdminSettingsController;
