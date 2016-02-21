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
var _ = require('underscore');
var Validate = Rfr('lib/utility/validate.js');
var Node = Rfr('lib/model/node.js');
var Uuid = Rfr('lib/data/uuid.js');

/**
 * Creates a new NodeController backed by data from the given collections
 *
 * @param {Object} collections - Data collections to use
 * @constructor
 */
var NodeController = function (collections) {
    this.collections = collections;
};

/**
 * Creates a new node and adds it to the database
 *
 * @param {String} name - Name of node
 * @param {Uuid} locationUUID - UUID for location
 * @param {String} ip - IP of node
 * @param {Number} port - Port of node
 * @param {Object[]} addresses - Array of addresses which can be allocated to servers
 */
NodeController.prototype.create = function (name, locationUUID, ip, port, addresses) {
    Validate.isString(name, 'name');
    Validate.isUUID(locationUUID, 'location');
    Validate.isIP(ip, 'ip');
    Validate.isPort(port, 'port');
    Validate.isArray(addresses, 'addresses');

    var self = this;
    var node = new Node(Uuid.generate(), name, locationUUID, ip, port, addresses);
    self.collections.node.add(node);
};

/**
 * Gets the details for a node
 *
 * @param {Uuid} uuid - UUID of node
 * @returns {Node} - Node object, or undefined if no node
 */
NodeController.prototype.get = function (uuid) {
    Validate.isUUID(uuid, 'uuid');

    var self = this;
    return self.collections.node.get(uuid);
};

/**
 * Gets the details for all nodes at a location
 *
 * @param {Uuid} locationUUID - Location UUID
 * @returns {Node[]} - Array of nodes at that location
 */
NodeController.prototype.getAtLocation = function (locationUUID) {
    Validate.isUUID(locationUUID, 'uuid');

    var self = this;
    self._validateLocation(locationUUID);
    return self.collections.node.getByLocation(locationUUID);
};

/**
 * Adds an address that may be used by servers
 *
 * @param {Uuid} nodeUUID - Node UUID
 * @param {String} ip - IP address
 * @param {Number} port - Port
 */
NodeController.prototype.addAddress = function (nodeUUID, ip, port) {
    Validate.isUUID(nodeUUID, 'uuid');
    Validate.isIP(ip, 'ip');
    Validate.isPort(port, 'port');

    var self = this;
    var node = self.collections.node.get(nodeUUID);
    if (node) {
        var addresses = node.getAddresses();
        addresses.push({ ip: ip, port: port, allocated: false });
        self.collections.node.update(node.getUUID(), { addresses: addresses });
    } else {
        throw new Error('No node with the given UUID', { uuid: nodeUUID });
    }
};

/**
 * Removes an address that may be used by servers
 *
 * @param {Uuid} nodeUUID - Node UUID
 * @param {String} ip - IP address
 * @param {Number} port - Port
 */
NodeController.prototype.removeAddress = function (nodeUUID, ip, port) {
    Validate.isUUID(nodeUUID, 'uuid');
    Validate.isIP(ip, 'ip');
    Validate.isPort(port, 'port');

    var self = this;
    var node = self.collections.node.get(nodeUUID);
    if (node) {
        var addresses = _.reject(node.getAddresses(), function (address) {
            return address.ip === ip && address.port === port;
        });
        self.collections.node.update(node.getUUID(), { addresses: addresses });
    } else {
        throw new Error('No node with the given UUID', { uuid: nodeUUID });
    }
};

/**
 * Updates an address that may be used by servers
 *
 * @param {Uuid} nodeUUID - Node UUID
 * @param {String} ip - IP address
 * @param {Number} port - Port
 * @param {Boolean} allocated - If the address is in use or not
 */
NodeController.prototype.updateAddress = function (nodeUUID, ip, port, allocated) {
    Validate.isUUID(nodeUUID, 'uuid');
    Validate.isIP(ip, 'ip');
    Validate.isPort(port, 'port');
    Validate.isBoolean(allocated, 'allocated');

    var self = this;
    var node = self.collections.node.get(nodeUUID);
    if (node) {
        var addresses = node.getAddresses();
        _.each(addresses, function (address) {
            if (address.ip === ip && address.port === port) {
                address.allocated = allocated;
            }
        });
        self.collections.node.update(node.getUUID(), { addresses: addresses });
    } else {
        throw new Error('No node with the given UUID', { uuid: nodeUUID });
    }
};

/**
 * @private
 */
NodeController.prototype._validateLocation = function (uuid) {
    var self = this;
    var user = self.collections.location.get(uuid);
    if (!user) {
        throw new Error('No location with the provided UUID', { uuid: uuid });
    }
};

module.exports = NodeController;