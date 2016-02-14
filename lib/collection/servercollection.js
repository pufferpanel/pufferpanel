/*
 Copyright 2016 Joshua Taylor <lordralex@gmail.com>

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 */

var Rfr = require('rfr');
var Rethink = require('rethinkdbdash');
var _ = require('underscore');
var Server = Rfr('lib/model/server.js');

/**
 * Creates a new UserCollection which relies on data from RethinkDB
 *
 * @param {Object} dbConn - Existing connection, may be undefined
 * @constructor
 */
var ServerCollection = function (dbConn) {
    this.db = dbConn ? dbConn : new Rethink({
        pool: false,
        discovery: true
    });
};

/**
 * Gets a server based on UUID
 *
 * @param {String} uuid - UUID for server
 * @returns {Object} - Server object, or undefined if no server exists
 */
ServerCollection.prototype.get = function (uuid) {
    var self = this;
    var data = yield self.db.table('servers').get(uuid).run();
    if (data) {
        return self._castFromDb(data);
    } else {
        return undefined;
    }
};

/**
 * Gets all servers owned by a specific user
 *
 * @param {String} uuid - UUID for user
 * @returns {Array} - All servers owned by the user, will be empty if no servers.
 */
ServerCollection.prototype.getByOwner = function (uuid) {
    var self = this;
    var results = yield self.db.table('servers').filter({uuid: uuid}).run();
    if (results) {
        var servers = _.map(results, function (data) {
            return self._castFromDb(data);
        });
        return servers;
    } else {
        return [];
    }
};

/**
 * Creates a new server and adds it to the database
 *
 * @param {String} uuid - UUID for the server
 * @param {String} name - Name for the server
 * @param {Object} address - The address for the server, following { ip, port }
 * @param {String} type - The type of server
 * @param {String} owner - The owner's UUID
 * @param {String} node - The node's UUID
 * @param {Object} data - Any extra data for the server
 * @returns {Server} - Newly created server
 */
ServerCollection.prototype.create = function (uuid, name, address, type, owner, node, data) {
    var self = this;
    var server = new Server(uuid, name, type, owner, node, false, data);
    yield self.db.table('servers').insert(server);
    return server;
};

/**
 * Deletes a server from the database
 *
 * @param {String} uuid - Server UUID
 */
ServerCollection.prototype.delete = function (uuid) {
    var self = this;
    yield self.db.table('servers').delete(uuid);
};

/**
 * @private
 */
ServerCollection.prototype._castFromDb = function (data) {
    //TODO: Determine how RethinkDB will store this and adapt logic
    var addresses = data.addresses;
    return new Server(data.uuid, data.name, data.type, addresses, data.owner, data.node, data.suspended, data.data);
};

module.exports = ServerCollection;