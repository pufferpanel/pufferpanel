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
var Node = Rfr('lib/model/node.js');

/**
 * Creates a new NodeCollection which relies on data from RethinkDB
 *
 * @param {Object} dbConn - Existing connection, may be undefined
 * @constructor
 */
var NodeCollection = function (dbConn) {
    this.db = dbConn ? dbConn : new Rethink({
        pool: false,
        discovery: true
    });
};

/**
 * Gets a node based on the UUID
 *
 * @param {String} uuid - UUID of node
 * @returns {Node} - Node, or undefined if no node
 */
NodeCollection.prototype.get = function (uuid) {
    var self = this;
    var result = yield self.db.table('nodes').get(uuid).run();
    if (result) {
        return NodeCollection._castFromDb(result);
    } else {
        return undefined;
    }
};

/**
 * Gets all nodes at a location
 *
 * @param location - UUID of location
 * @returns {Array} - Array of Nodes at the given location
 */
NodeCollection.prototype.getByLocation = function (location) {
    var self = this;
    var nodes = self.db.table('nodes').filter({ location: location }).run();
    if (nodes) {
        return _.map(nodes, NodeCollection._castFromDb);
    } else {
        return [];
    }
};

/**
 * Adds a new node to the database
 *
 * @param {Node} node - Node to add
 */
NodeCollection.prototype.add = function (node) {
    var self = this;
    yield self.db.table('nodes').insert(node);
};

/**
 * Removes a node from the database
 *
 * @param {String} uuid - UUID of node
 */
NodeCollection.prototype.remove = function (uuid) {
    var self = this;
    yield self.db.table('nodes').delete(uuid);
};

/**
 * Updates a node with new values
 *
 * @param {String} uuid - UUID of node
 * @param {Object} newValues - New values for the node
 * @returns {Node} - The updated node
 */
NodeCollection.prototype.update = function (uuid, newValues) {
    var self = this;
    yield self.db.table('nodes').get(uuid).update(newValues).run();
    return self.get(uuid);
};

/**
 * @private
 */
NodeCollection._castFromDb = function (data) {
    //TODO: Determine how RethinkDB will return this data
    var addresses = data.addresses;
    return new Node(data.uuid, data.name, data.location, data.ip, data.port, addresses);
};

module.exports = NodeCollection;