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
var Joi = require('joi');
var ServerModel = Rfr('lib/models/servers.js');
var UserVisibleError = Rfr('lib/errors/UserVisibleError.js');

/** @namespace */
var AdminServerController = {};

AdminServerController.listAllServers = function (next) {

    ServerModel.select({}, function (err, servers) {

        if (err) {
            return next(err);
        }

        return next(null, servers);

    });

};

module.exports = AdminServerController;
