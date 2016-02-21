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
var Crypto = require('crypto');
var Bcrypt = require('bcrypt');
var Uuid = Rfr('lib/data/uuid.js');
var User = Rfr('lib/model/user.js');
var Validate = Rfr('lib/utility/validate.js');

/**
 * Creates a new UserController backed by data from the given collections
 *
 * @param {Object} collections - Data collections to use
 * @constructor
 */
var UserController = function (collections) {
    this.collections = collections;
};

/**
 * Creates a new user with the given username and email
 *
 * @param {String} username - Username for the user
 * @param {String} email - Email for the user.
 * @param {String} password - Unhashed password for the user
 * @returns {User|Boolean} - Resulting user, or false if the username or email already existed
 */
UserController.prototype.create = function (username, email, password) {
    Validate.isString(username, 'username');
    Validate.isEmail(email, 'email');
    Validate.isString(password, 'password');

    var self = this;
    var user = new User(Uuid.v4(), email, username, Bcrypt.hashSync(password), false, undefined);
    if (self.collections.user.getBy({ username: username }) || self.collections.user.getBy({ email: email })) {
        return false;
    }
    self.collections.user.add(user);
    return user;
};

/**
 * Validates if the session token is correct for the given user
 *
 * @param {Uuid} userUUID - User who claims the session key
 * @param {String} session - Session key
 * @returns {Boolean} - True if the user's session key is correct, otherwise false
 */
UserController.prototype.validateSession = function (userUUID, session) {
    Validate.isUUID(userUUID, 'uuid');
    Validate.isString(session, 'session');

    var self = this;
    var user = self.collections.user.get(userUUID);
    if (!user) {
        return false;
    }
    return user.getSession() == session;
};

/**
 * Validates a user's login credentials, then updates their session if correct.
 *
 * @param {String} email - User's email
 * @param {String} password - Unhashed password
 * @returns {String|Boolean} - Session key if auth was successful, otherwise false
 */
UserController.prototype.loginUser = function (email, password) {
    Validate.isEmail(email, 'email');
    Validate.isString(password, 'password');

    var self = this;
    if (self.validateAccountInfo(email, password)) {
        var user = self.collections.user.getBy({ email: email });
        var sessionToken = Crypto.randomBytes(64).toString('hex');
        self.collections.user.update(user.uuid, { session: sessionToken });
        return sessionToken;
    } else {
        return false;
    }
};

/**
 * Validates if a given email and password is valid for a user
 *
 * @param {String} email - Email of user
 * @param {String} password - Unhashed password
 * @returns {Boolean} - True if the information is correct, false otherwise
 */
UserController.prototype.validateAccountInfo = function (email, password) {
    Validate.isEmail(email, 'email');
    Validate.isString(password, 'password');

    var self = this;
    var user = self.collections.user.getBy({ email: email });
    if (!user) {
        return false;
    }
    var userPw = user.getPassword();
    return Bcrypt.compareSync(password, userPw);
};

/**
 * Toggles the suspended state of a user.
 * This will also clear their existing session should they have one.
 *
 * @param {Uuid} userUUID - User UUID
 * @param {String} suspended - Whether to suspend or unsuspend the account
 */
UserController.prototype.suspend = function (userUUID, suspended) {
    Validate.isUUID(userUUID, 'uuid');
    Validate.isBoolean(suspended, 'suspended');

    var self = this;
    self._updateUser(userUUID, { suspended: suspended, session: undefined });
};

/**
 * Changes a user's password
 *
 * @param {Uuid} userUUID - User UUID to alter
 * @param {String} password - The new unhashed password
 */
UserController.prototype.changePassword = function (userUUID, password) {
    Validate.isUUID(userUUID, 'uuid');
    Validate.isString(password, 'password');

    var self = this;
    var pw = Bcrypt.hashSync(password, 10);
    self._updateUser(userUUID, { password: pw });
};

/**
 * Changes a user's email
 *
 * @param {Uuid} userUUID - User UUID to alter
 * @param {String} email - The new email
 */
UserController.prototype.changeEmail = function (userUUID, email) {
    Validate.isUUID(userUUID, 'uuid');
    Validate.isEmail(email, 'email');

    var self = this;
    self._updateUser(userUUID, { email: email });
};

/**
 * Changes a user's username
 *
 * @param {Uuid} userUUID - User UUID to alter
 * @param {String} username - The new username
 */
UserController.prototype.changeUsername = function (userUUID, username) {
    Validate.isUUID(userUUID, 'uuid');
    Validate.isString(username, 'username');

    var self = this;
    self._updateUser(userUUID, { username: username });
};

/**
 * Gets information about a user
 *
 * @param {Uuid} userUUID - User UUID to get
 * @returns {User} - User specified, or undefined if no user
 */
UserController.prototype.getUser = function (userUUID) {
    Validate.isUUID(userUUID, 'uuid');

    var self = this;
    return self.collections.user.get(userUUID);
};

/**
 * Convenience method to handle updating a user
 *
 * @param {Uuid} userUUID - User UUID
 * @param {Object} data - Values to update
 * @private
 */
UserController.prototype._updateUser = function (userUUID, data) {
    Validate.isUUID(userUUID, 'uuid');
    Validate.isObject(data, 'data');

    self.collections.user.update(userUUID, data);
};

module.exports = UserController;