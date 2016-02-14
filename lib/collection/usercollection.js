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
var User = Rfr('lib/model/user.js');

/**
 * Creates a new UserCollection which relies on data from RethinkDB
 *
 * @param {Object} dbConn - Existing connection, may be undefined
 * @constructor
 */
var UserCollection = function (dbConn) {
    this.db = dbConn ? dbConn : new Rethink({
        pool: false,
        discovery: true
    });
};

/**
 * Gets a user based on the given UUID
 *
 * @param {String} uuid - The user's UUID
 * @returns {User} - The resulting user, or undefined if no user exists
 */
UserCollection.prototype.get = function (uuid) {
    var self = this;
    var data = yield self.db.table('users').get(uuid).run();
    if (data) {
        return new User();
    } else {
        return undefined;
    }
};

/**
 * Adds a user to the database
 *
 * @param {User} user - User to add
 */
UserCollection.prototype.add = function (user) {
    var self = this;
    yield self.db.table('users').insert(user);
};

module.exports = UserCollection;