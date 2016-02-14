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
 * Creates a new Node object. This does NOT add the node to the database.
 *
 * @param {String} uuid - UUID for the node
 * @param {String} name - Name for the node
 * @param {String} location - The UUID for the location this node is at
 * @param {String} ip - The IP for this node
 * @param {Number} port - The port for this node
 * @param {Array} addresses - Addresses this node may allocate
 * @constructor
 */
var Node = function (uuid, name, location, ip, port, addresses) {
    this.uuid = uuid;
    this.name = name;
    this.location = location;
    this.ip = ip;
    this.port = port;
    this.addresses = addresses;
};

/**
 * Gets the node's UUID
 *
 * @returns {String} - UUID for the node
 */
Node.prototype.getUUID = function () {
    return this.uuid;
};

/**
 * Get the node's name
 *
 * @returns {String} - Name for the node
 */
Node.prototype.getName = function () {
    return this.name;
};

/**
 * Get the location's UUID where this node is located at.
 *
 * @returns {String} - UUID of the location
 */
Node.prototype.getLocationUUID = function () {
    return this.location;
};

/**
 * Get's the addresses this node will allocate to servers.
 * Array returned will contain Objects in the following format:
 * {
 *      ip: "1.2.3.4",
 *      port: 1234,
 *      allocated: false
 * }
 *
 * @returns {Array} - Array of IP blocks
 */
Node.prototype.getAddresses = function () {
    return this.addresses;
};

/**
 * Gets the IP for the node
 *
 * @returns {String} - IP of node
 */
Node.prototype.getIP = function () {
    return this.ip;
};

/**
 * Get's the port the node uses for communication
 *
 * @returns {Number} - Port number
 */
Node.prototype.getPort = function () {
    return this.port;
};

module.exports = Node;