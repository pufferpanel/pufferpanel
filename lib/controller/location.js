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
var Validate = Rfr('lib/utility/validate.js');
var Location = Rfr('lib/model/location.js');

/**
 * Creates a new LocationController backed by data from the given collections
 *
 * @param {Object} collections - Data collections to use
 * @constructor
 */
var LocationController = function (collections) {
    this.collections = collections;
};

/**
 * Creates a new location with the give name
 *
 * @param {String} name - Name of location
 * @returns {Location} - New Location
 */
LocationController.prototype.create = function (name) {
    Validate.isString(name, 'name');

    var self = this;
    var location = new Location(Uuid.v4(), name);
    self.collections.location.add(location);
    return location;
};

/**
 * Changes the name of a location
 *
 * @param {String} uuid - Location UUID
 * @param {String} name - New name
 */
LocationController.prototype.changeName = function (uuid, name) {
    Validate.isUUID(uuid, 'uuid');
    Validate.isString(name, 'name');

    var self = this;
    self.collections.location.update(uuid, { name: name });
};

/**
 * Gets a location
 *
 * @param {String} uuid - UUID of location
 * @returns {Location} - Location, or undefined if it does not exist
 */
LocationController.prototype.get = function (uuid) {
    Validate.isUUID(uuid, 'uuid');

    var self = this;
    return self.collections.location.get(uuid);
};

/**
 * Deletes a location
 *
 * @param {String} uuid - UUID to delete
 */
LocationController.prototype.delete = function (uuid) {
    Validate.isUUID(uuid, 'uuid');

    var self = this;
    self.collections.location.remove(uuid);
};

module.exports = LocationController;