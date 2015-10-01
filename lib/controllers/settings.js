/*
 * PufferPanel - Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var SettingModels = (process.env.NODE_ENV === 'test') ? Rfr('tests/models/settings.js') : Rfr('lib/models/settings.js');

/** @namespace */
var SettingsController = {};


/**
 * Returns the value of a setting in the database.
 *
 * @param {setting} setting - The requested setting value.
 * @param {SettingsController~getCb} next - Callback to handle response
 */
SettingsController.get = function (setting, next) {

    SettingModels.select( { key: setting }, function (err, results) {

        if (err) {
            return next(err);
        }

        return next(null, results.value);
    });
};
/**
 * @callback SettingsController~getCb
 * @param {Error} err - Error if one exists, otherwise null
 * @param {Object} results - Setting requested
 */

SettingsController.getAllSettings = function (next) {

    SettingModels.selectAll(function (err, results) {

        if (err) {
            return next(err);
        }

        // Process the Results
        var formattedResults = {};
        for (var p in results) {
            if (results.hasOwnProperty(p)) {
                formattedResults[results[p].key] = results[p].value;
            }
        }
        return next(null, formattedResults);

    });
};


module.exports = SettingsController;
