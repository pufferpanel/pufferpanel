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
var Rethink = Rfr('lib/rethink.js');
var Async = require('async');

/** @namespace */
var SettingsModel = {};

SettingsModel.select = function (criteria, next) {

    Rethink.table('settings').filter(criteria).run().then(function (setting) {

        return next(null, _.first(setting));
    }).error(function (err) {

        return next(err);
    });
};

SettingsModel.selectAll = function (next) {

    Rethink.table('settings').run().then(function (settings) {

        return next(null, settings);
    }).error(function (err) {

        return next(err);
    });

};

SettingsModel.update = function (updates, next) {

    Async.forEachOf(updates, function (newValue, setKey, callback) {

        Rethink.table('settings').filter({ key: setKey }).update({ value: newValue }).run().then(function () {

            return callback();
        }).error(function (err) {

            return callback(err);
        });

    }, function (err) {

        if (err) {
            return next(err);
        }

        return next();
    });

};

module.exports = SettingsModel;
