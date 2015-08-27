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
var CoreSettings = Rfr('lib/models/coresettings.js');

/** @namespace */
var CoreSettingsController = {};


/**
 * Returns the value of a setting in the database.
 *
 * @param {setting} setting - The requested setting value.
 */
CoreSettingsController.returnSetting = function (setting, next) {

    CoreSettings.select( { key: setting }, function (err, results) {

        if (err) {

            Logger.error('An error occured in CoreSettingsController.returnSetting at selection.', err);
            return next(err);
        }

        return next(null, results.value);

    });

};

module.exports = CoreSettingsController;
