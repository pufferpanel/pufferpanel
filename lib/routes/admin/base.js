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
var Logger = Rfr('lib/api/logger.js');
var Fs = require('fs-extra');
var Boom = require('boom');
var Util = require('util');
var Routes = {
    settings: Rfr('lib/routes/admin/settings.js'),
    servers: Rfr('lib/routes/admin/servers.js')
};

Routes.getIndex = function (request, reply) {

    Fs.readFile('./.git/HEAD', function (err, data) {

        if (err) {
            Logger.error('An error occured while attempting to process .git information.', err);
            return reply(Boom.badImplementation());
        }

        if (data.indexOf('ref: ') > -1) {

            var ref = data.toString().split(' ');
            Fs.readFile(Util.format('./.git/%s', ref[1].trim()), function (err, moreData) {

                if (err) {
                    Logger.error('An error occured while attempting to process more detailed .git information.', err);
                    return reply(Boom.badImplementation());
                }

                return reply.view('admin/index', {
                    version: ref[1].trim(),
                    sha: moreData.toString().trim()
                });

            });

        } else {

            return reply.view('admin/index', {
                version: 'master',
                sha: data.toString().trim()
            });

        }

    });

};

module.exports = Routes;
