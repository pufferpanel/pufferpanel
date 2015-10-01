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

/** @namespace */
var TestSettingsModel = {};

TestSettingsModel.select = function (criteria, next) {

    return next(null, _.findWhere(_fakeData, criteria));
};

TestSettingsModel.selectAll = function (next) {

    return next(null, _fakeData);
};

TestSettingsModel.update = function (updates, next) {

    Async.forEachOf(updates, function (newValue, setKey, callback) {

        var setting = _.findWhere(_fakeData, { key: setKey });
        if (setting === undefined) {
            return callback(new Error('No such setting.'));
        }
        _.extend(setting, { value: newValue });
        return callback();

    }, function (err) {
        return next(err);
    });

};

module.exports = TestSettingsModel;
