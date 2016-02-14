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

/**
 * Creates a new Server object. This does NOT add it to the database.
 *
 * @param {String} uuid - The UUID for this server
 * @param {String} name - The name for this server
 * @param {String} type - The type of server this is
 * @param {Array} addresses - Bind addresses this server can use
 * @param {String} ownerUuid - The owner's UUID
 * @param {String} nodeUuid - The node's UUID
 * @param {Boolean} suspended - Whether or not the server is suspended
 * @param {Object} data - Data related to this server
 * @constructor
 */
var Server = function (uuid, name, type, addresses, ownerUuid, nodeUuid, suspended, data) {
    this.uuid = uuid;
    this.name = name;
    this.type = type;
    this.owner = ownerUuid;
    this.node = nodeUuid;
    this.suspended = suspended;
    this.data = data;
    this.addresses = address;
};

/**
 * Get's this server's UUID
 *
 * @returns {String} - UUID for the server
 */
Server.prototype.getUUID = function () {
    return this.uuid;
};

/**
 * Get's this server's name
 *
 * @returns {String} - Name for the server
 */
Server.prototype.getName = function () {
    return this.name;
};

/**
 * Get's the type of this server
 *
 * @returns {String} - Type of server
 */
Server.prototype.getType = function () {
    return this.type;
};

/**
 * Get's the owner's UUID of this server
 *
 * @returns {String} - Owner's UUID
 */
Server.prototype.getOwnerUUID = function () {
    return this.owner;
};

/**
 * Get's the suspended status for this server
 *
 * @returns {Boolean} - True if the server is suspended, otherwise false
 */
Server.prototype.isSuspended = function () {
    return this.suspended;
};

/**
 * Get's the data for this server, which includes server-specific values.
 * Examples may include (but not always):
 * - Ram allocation
 * - Max slots
 * - Map selection
 *
 * @returns {Object} - Server's data object
 */
Server.prototype.getData = function () {
    return this.data;
};

/**
 * Get's the node's UUID this server is on
 *
 * @returns {String} - UUID for node
 */
Server.prototype.getNodeUUID = function () {
    return this.node;
};

/**
 * Get's the addresses this server may use.
 * Array is a collection of Objects following the format:
 * {
 *     ip: "1.2.3.4",
 *     port: 1234,
 *     primary: false
 * }
 *
 * Primary indicates if the IP object is the main IP for this server.
 *
 * @returns {Array} - Bindable addresses
 */
Server.prototype.getAddresses = function () {
    return this.addresses;
};

module.exports = Server;