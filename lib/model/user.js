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
 * Creates a new User object for tracking. This will NOT create a new user in the system.
 *
 * @param {String} uuid - UUID for the user
 * @param {String} email - Email for the user
 * @param {String} username - Username for the user
 * @param {String} password - Hashed password for the user
 * @param {Boolean} suspended - If the user is currently suspended
 * @constructor
 */
var User = function (uuid, email, username, password, suspended, session) {
    this.uuid = uuid;
    this.email = email;
    this.username = username;
    this.password = password;
    this.suspended = suspended;
    this.session = session;
};

/**
 * Gets the user's UUID
 *
 * @returns {String} - The UUID for the user
 */
User.prototype.getUUID = function () {
    return this.uuid;
};

/**
 * Gets the user's username
 *
 * @returns {String} - The username for this user
 */
User.prototype.getUsername = function () {
    return this.username;
};

/**
 * Get's the user's email
 *
 * @returns {String} - The email for this user
 */
User.prototype.getEmail = function () {
    return this.email;
};

/**
 * Get's whether the user is currently suspended
 *
 * @returns {Boolean} - True if the user is suspended, otherwise false
 */
User.prototype.isSuspended = function () {
    return this.suspended;
};

/**
 * Get's this user's hashed password
 *
 * @returns {String} - User's password
 */
User.prototype.getPassword = function () {
    return this.password;
};

/**
 * Get's this user's session key
 *
 * @returns {String} - Session key
 */
User.prototype.getSession = function () {
    return this.session;
};

module.exports = User;