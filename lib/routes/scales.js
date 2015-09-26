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

var Routes = {
    mapUri: function (request, callback) {
        // Demo Data to Test with Locally
        callback(null, Util.format('https://192.168.0.12:5656/%s', request.params.path), {
            'X-Access-Server': '08a980a8-107a-4a67-a800-d49839d17c3c',
            'X-Access-Token': 'c9dd4b91-b4fc-4fff-b151-36eab7852c09'
        });
    },
    onResponse:function (err, res, request, reply, settings, ttl) {
        Wreck.read(res, { json: true }, function (err, payload) {
            reply(payload);
        });
    }
};

module.exports = Routes;
