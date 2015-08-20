/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Async = require('async');
var Database = requireFromRoot('lib/database');
var Logger = requireFromRoot('lib/logger');

var Routes = {
  get: {
    index: function (request, reply) {

      var userServers;
      Async.series([
        function (callback) {
          //TODO: Pass this off to a helper function in deeper code
          Database.get('servers', 'owner', request.auth.credentials.id, function (servers, err) {
            if (err !== undefined) {
              Logger.err(err);
              userServers = [];
            } else {
              userServers = servers;
            }
          });
        },
        function (callback) {
          reply.view('base/index', {
            servers: userServers,
            user: request.auth.credentials
          });
        }
      ]);

    }
  }
};

module.exports = Routes;
