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
var Rethink = Rfr('lib/rethink.js');

/** @namespace */
var AdminSettingsModel = {};

AdminSettingsModel.update = function (keyName, newKeyValue) {

    Rethink.table('settings').filter({ key: keyName }).update({ value: newKeyValue }).run().then(function () {

        return next();
    }).error(function (err) {

        return next(err);
    });

};
