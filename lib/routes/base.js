/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Async = require('async');
var Rethink = requireFromRoot('lib/rethink');
var Logger = requireFromRoot('lib/logger');

var Routes = {
  get: {
    index: function (request, reply) {

      var userServers;
      Async.series([
        function (callback) {
          Rethink.table('servers').filter(Rethink.row('owner').eq(request.auth.credentials.id)).run().then(function (servers) {
            userServers = servers;
            callback();
          }).error(function (err) {
            Logger.error(err);
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
