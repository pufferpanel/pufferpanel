/*
 * PufferPanel � Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var _ = require('underscore');
var Rfr = require('rfr');
var Rethink = Rfr('lib/rethink.js');

/** @namespace */
var CoreSettingsModel = {};

CoreSettingsModel.select = function (criteria, next) {

    Rethink.table('settings').filter(criteria).run().then(function (setting) {

        return next(null, _.first(setting));
    }).error(function (err) {

        return next(err);
    });
};

module.exports = CoreSettingsModel;
