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
        return ServerCollection._castFromDb(data);
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
    var results = yield self.db.table('servers').filter({ uuid: uuid }).run();
    if (results) {
        return _.map(results, ServerCollection._castFromDb);
    } else {
        return [];
    }
};

/**
 * Adds a server to the database
 *
 * @param {Server} server - Server to add
 */
ServerCollection.prototype.add = function (server) {
    var self = this;
    yield self.db.table('servers').insert(server);
};

/**
 * Removes a server from the database
 *
 * @param {String} uuid - Server UUID
 */
ServerCollection.prototype.remove = function (uuid) {
    var self = this;
    yield self.db.table('servers').delete(uuid);
};

/**
 * Updates a server with new values
 *
 * @param {String} uuid - UUID of server
 * @param {Object} newValues - New values for the server
 * @returns {Server} - The updated server
 */
ServerCollection.prototype.update = function (uuid, newValues) {
    var self = this;
    yield self.db.table('servers').get(uuid).update(newValues).run();
    return self.get(uuid);
};

/**
 * @private
 */
ServerCollection._castFromDb = function (data) {
    //TODO: Determine how RethinkDB will store this and adapt logic
    var addresses = data.addresses;
    return new Server(data.uuid, data.name, data.type, addresses, data.owner, data.node, data.suspended, data.data);
};

module.exports = ServerCollection;