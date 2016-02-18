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
var Uuid = require('uuid');
var _ = require('underscore');
var Server = Rfr('lib/model/server.js');
var Validate = Rfr('lib/utility/validate.js');

/**
 * Creates a new ServerController backed by data from the given collections
 *
 * @param {Object} collections - Data collections to use
 * @constructor
 */
var ServerController = function (collections) {
    this.collections = collections;
};

/**
 * Creates a new server and adds it to the database
 *
 * @param {String} name - Server name
 * @param {String} ownerUUID - Owner's UUID
 * @param {String} nodeUUID - Node's UUID
 * @param {String} plugin - Server plugin type
 * @param {String} ip - Main IP
 * @param {Number} port - Main port
 * @param {Object} data - Extra plugin data
 * @returns {Server} - Created server
 */
ServerController.prototype.create = function (name, ownerUUID, nodeUUID, plugin, ip, port, data) {
    Validate.isString(name, 'name');
    Validate.isUUID(ownerUUID, 'ownerUUID');
    Validate.isUUID(nodeUUID, 'nodeUUID');
    Validate.isString(plugin, 'plugin');
    Validate.isIP(ip, 'ip');
    Validate.isNumber(port, 'port');
    Validate.isObject(data, 'data');

    var self = this;
    self._validateNode(nodeUUID);
    self._validateUser(ownerUUID);
    var server = new Server(Uuid.v4(), name, type, [{
        ip: ip, port: port
    }], ownerUUID, nodeUUID, false, data);
    self.collections.server.add(server);
    return server;
};

/**
 * Gets information about a specific server
 *
 * @param {String} uuid - UUID of server
 * @returns {Server} - Server object, or undefined if no server
 */
ServerController.prototype.get = function (uuid) {
    Validate.isUUID(uuid, 'uuid');

    var self = this;
    var server = self.collections.server.get(uuid);
    return server;
};

/**
 * Gets all servers owned by the specified user
 *
 * @param {String} userUUID - User UUID
 * @returns {Server[]} - Array of servers this user owns
 */
ServerController.prototype.getForUser = function (userUUID) {
    Validate.isUUID(userUUID, 'userUUID');

    var self = this;
    self._validateUser(userUUID);
    var servers = self.collections.server.getByOwner(userUUID);
    return servers;
};

/**
 * Deletes a server from the database
 *
 * @param {String} uuid - Server UUID
 */
ServerController.prototype.delete = function (uuid) {
    Validate.isUUID(uuid, 'uuid');

    var self = this;
    self.collections.server.remove(uuid);
};

/**
 * Updates a server's name
 *
 * @param {String} uuid - Server UUID
 * @param {String} name - New name
 */
ServerController.prototype.updateName = function (uuid, name) {
    Validate.isUUID(uuid, 'uuid');
    Validate.isString(name, 'name');

    var self = this;
    self.collections.server.update(uuid, { name: name });
};

/**
 * Updates the owner of a server
 *
 * @param {String} serverUUID - Server UUID
 * @param {String} userUUID - New owner UUID
 */
ServerController.prototype.updateOwner = function (serverUUID, userUUID) {
    Validate.isUUID(serverUUID, 'serverUUID');
    Validate.isUUID(userUUID, 'userUUID');

    var self = this;
    self._validateUser(userUUID);
    self.collections.server.update(serverUUID, { owner: userUUID });
};

/**
 * Adds a new bindable address to the server
 *
 * @param {String} serverUUID - Server UUID
 * @param {String} ip - IP of address
 * @param {Number} port - Port of address
 * @param {Boolean} primary - True if this is the main address, false otherwise
 */
ServerController.prototype.addAddress = function (serverUUID, ip, port, primary) {
    Validate.isUUID(serverUUID, 'uuid');
    Validate.isIP(ip, 'ip');
    Validate.isNumber(port, 'port');
    Validate.isBoolean(primary, 'primary');

    var self = this;
    var server = self.collections.server.get(serverUUID);
    if (server) {
        var addresses = server.getAddresses();
        if (primary) {
            _.each(addresses, function (address) {
                address.primary = false;
            });
        }
        addresses.push({ ip: ip, port: port, primary: primary });
        self.collections.server.update(server.getUUID(), addresses);
    } else {
        throw new Error('No server with the given UUID');
    }
};

/**
 * Adds a new bindable address to the server
 *
 * @param {String} serverUUID - Server UUID
 * @param {String} ip - IP of address
 * @param {Number} port - Port of address
 * @param {Boolean} primary - True if this is the main address, false otherwise
 */
ServerController.prototype.removeAddress = function (serverUUID, ip, port) {
    Validate.isUUID(serverUUID, 'uuid');
    Validate.isIP(ip, 'ip');
    Validate.isNumber(port, 'port');
    Validate.isBoolean(primary, 'primary');

    var self = this;
    var server = self.collections.server.get(serverUUID);
    if (server) {
        var addresses = _.reject(server.getAddresses(), function (address) {
            return address.ip === ip && address.port === port;
        });
        self.collections.server.update(server.getUUID(), addresses);
    } else {
        throw new Error('No server with the given UUID');
    }
};

/**
 * Changes the plugin data for a server.
 * This will *override* existing data.
 * @param {String} serverUUID - Server UUID
 * @param {Object} data - Plugin data
 */
ServerController.prototype.changePluginData = function (serverUUID, data) {
    Validate.isUUID(serverUUID, 'uuid');
    Validate.isObject(data, 'data');

    self.collections.server.update(serverUUID, { data: data });
};

ServerController.prototype._validateUser = function (uuid) {
    var user = self.collections.user.get(uuid);
    if (!user) {
        throw new Error('No user with the provided UUID', { uuid: uuid });
    }
};

ServerController.prototype._validateNode = function (uuid) {
    var node = self.collections.node.get(uuid);
    if (!node) {
        throw new Error('No node with the provided UUID', { uuid: uuid });
    }
};

module.exports = ServerController;