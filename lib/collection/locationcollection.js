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
var Location = Rfr('lib/model/location.js');

/**
 * Creates a new LocationCollection which relies on data from RethinkDB
 *
 * @param {Object} dbConn - Existing connection, may be undefined
 * @constructor
 */
var LocationCollection = function (dbConn) {
    this.db = dbConn ? dbConn : new Rethink({
        pool: false,
        discovery: true
    });
};

/**
 * Gets a location based on the UUID
 *
 * @param {String} uuid - UUID of location
 * @returns {Location} - Location, or undefined if no location
 */
LocationCollection.prototype.get = function (uuid) {
    var self = this;
    var result = yield self.db.table('locations').get(uuid).run();
    if (result) {
        return new Location(data.uuid, data.name);
    } else {
        return undefined;
    }
};

/**
 * Adds a new location to the database
 *
 * @param {Location} node - Location to add
 */
LocationCollection.prototype.add = function (location) {
    var self = this;
    yield self.db.table('locations').insert(location);
};

/**
 * Removes a location from the database
 *
 * @param {String} uuid - UUID of location
 */
LocationCollection.prototype.remove = function (uuid) {
    var self = this;
    yield self.db.table('locations').delete(uuid);
};