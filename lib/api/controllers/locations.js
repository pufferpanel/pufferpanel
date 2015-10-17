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
var LocationModels = (process.env.NODE_ENV === 'test') ? Rfr('tests/models/locations.js') : Rfr('lib/api/models/locations.js');

/** @namespace */
var LocationController = {};

/**
 * Return all locations in database
 *
 * @param {LocationController~getAllLocationsCb} callback - Callback to handle response
 */
LocationController.getAllLocations = function (next) {

    var locationsArray = [];
    LocationModels.selectAll(function (err, locations) {

        if (err) {
            return next(err);
        }

        _.each(locations, function (value) {
            locationsArray[value.id] = {
                'long': value.long,
                'short': value.short
            };
        });

        return next(null, locationsArray);
    });
};

/**
 * @callback LocationController~getAllLocationsCb
 * @params {Error} err - Error if one occurred, otherwise null
 * @params {Object[]} locations - Array of all locations
 */

module.exports = LocationController;
