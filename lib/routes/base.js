/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Logger = Rfr('lib/logger.js');
var Servers = Rfr('lib/controllers/server.js');

var Routes = {
  get: {
    index: function (request, reply) {

      Servers.getServersFor(request.auth.credentials.id, function (err, servers) {

        if (err !== undefined) {
          Logger.error(err);
        }

        reply.view('base/index', {
          servers: servers || [],
          user: request.auth.credentials
        });

      });

    }
  }
};

module.exports = Routes;
