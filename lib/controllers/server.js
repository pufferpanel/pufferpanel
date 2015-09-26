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
var ServerModels = (process.env.NODE_ENV === 'test') ? Rfr('tests/models/servers.js') : Rfr('lib/models/servers.js');

/** @namespace */
var ServerController = {};

/**
 * Gets all servers this user has access to.
 * This would include servers they are owner of, and are a sub-user of.
 *
 * @param {String} userId - User's id
 * @param {Servers~getServersForCb} callback - Callback to handle response
 */
ServerController.getServersFor = function (userId, next) {

    var userServers = [];

    ServerModels.getByOwner(userId, function (err, servers) {

        if (err) {
            return next(err);
        }

        userServers = servers;

        //TODO: get servers the user is a sub-user of and append to userServers

        return next(null, userServers);
    });
};

/**
 * @callback ServerController~getServersForCb
 * @params {Error} err - Error if one occurred, otherwise null
 * @params {Object[]} servers - Array of Servers which the user has access to
 */

module.exports = ServerController;
