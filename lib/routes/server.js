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
var ServerApi = Rfr('lib/api/servers.js');
var ServerController = Rfr('lib/api/controllers/server.js');

var Routes = {
    handler: function (request, reply) {

        //TODO: Get server from controller
        var server = {
            plugin: 'minecraft'
        };

        var route = ServerApi._parseRoute(server.plugin, request.params.path);
        return route.handler(request, reply, server);
    }
};

module.exports = Routes;
