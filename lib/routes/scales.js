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
var Logger = Rfr('lib/logger.js');
var Wreck = require('wreck');
var Util = require('util');
var ServerController = Rfr('lib/controllers/server.js');

var Routes = {
    mapUri: function (request, callback) {

        ServerController.scalesProxyRequest(request.params.server, request.auth.credentials.id, function (err, server, node) {

            if (err) {
                Logger.error('An error occured while attempting to run a proxy request.', err);
                return callback('err');
            }

            if (typeof request.params.path === 'undefined') {
                request.params.path = '';
            }

            return callback(null, Util.format('https://%s:%s/%s', node.fqdn, node.daemon.listen, request.params.path), {
                'X-Access-Server': server.id,
                'X-Access-Token': server.daemon.secret
            });

        });

    },
    onResponse: function (err, res, request, reply, settings, ttl) {

        if (err) {
            if (err.isBoom) {
                return reply(JSON.stringify({
                    error: true,
                    errorCode: err.output.payload.statusCode,
                    message: 'Unable to complete request. Please try again.'
                })).type('application/json');
            }
        }

        Wreck.read(res, { json: true }, function (err, payload) {
            reply(payload);
        });
    }
};

module.exports = Routes;
