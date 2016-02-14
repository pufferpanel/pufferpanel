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
 * Creates a new Location with the given UUID and name. This will NOT add
 * the location to the database.
 *
 * @param {String} uuid - UUID for the location
 * @param {String} name - Name for the location
 * @constructor
 */
var Location = function (uuid, name) {
    this.uuid = uuid;
    this.name = name;
};

/**
 * Get's the UUID of this location
 *
 * @returns {String} - UUID for the location
 */
Location.prototype.getUuid = function () {
    return this.uuid;
};

/**
 * Gets the name of this location
 *
 * @returns {String} - Name of the location
 */
Location.prototype.getName = function () {
    return this.name;
};

module.exports = Location;