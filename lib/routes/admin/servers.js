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
var Logger = Rfr('lib/logger.js');
var Boom = require('boom');
var Util = require('util');
var AdminServerController = Rfr('lib/controllers/admin/servers.js');
var UserVisibleError = Rfr('lib/errors/UserVisibleError.js');
var Routes = {};

Routes.getListAllServers = function (request, reply) {

    AdminServerController.listAllServers(function (err, servers) {

        if (err) {
            return reply(Boom.badImplementation());
        }

        return reply.view('admin/servers/index', {
            servers: servers
        });

    });

};

Routes.getAddNewServer = function (request, reply) {

};

module.exports = Routes;
