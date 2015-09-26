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
var AdminSettingsModel = Rfr('lib/models/admin/settings.js');
var Joi = require('joi');

/** @namespace */
var AdminSettingsController = {};

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
            return next(null, errorResponse);
        }

        AdminSettingsModel.update({
            'companyName': name
        }, function (err) {

            if (err) {
                return next(err);
            }

            return next();
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
            return next(null, errorResponse);
        }

        AdminSettingsModel.update({
            'useSecureConnection': secureConnection,
            'enableSubusers': subUsers
        }, function (err) {

            if (err) {
                return next(err);
            }

            return next();
        });

    });

};

module.exports = AdminSettingsController;
