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
var LocationsModel = {};

LocationsModel.selectAll = function (next) {

    Rethink.table('locations').run().then(function (locations) {

        return next(null, locations);
    }).error(function (err) {

        return next(err);
    });
};

module.exports = LocationsModel;
