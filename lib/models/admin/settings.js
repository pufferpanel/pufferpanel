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
var Async = require('async');
var Rfr = require('rfr');
var Rethink = Rfr('lib/rethink.js');

/** @namespace */
var AdminSettingsModel = {};

AdminSettingsModel.update = function (updates, next) {

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

module.exports = AdminSettingsModel;
