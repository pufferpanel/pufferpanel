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

var UuidGen = require('uuid');
var Rfr = require('rfr');
var Validate = Rfr('lib/utility/validate.js');

/**
 * Generates a UUID type from a string
 *
 * @param {String} value - Version 4 UUID string
 * @constructor
 */
var Uuid = function Uuid(value) {
    Validate.isUUID(value, 'value');
    this.value = value;
};

Uuid.prototype.toString = function () {
    return this.value;
};

/**
 * Generates a random UUID
 *
 * @returns {Uuid} - New UUID
 */
Uuid.generate = function () {
    return new Uuid(UuidGen.v4());
};

module.exports = Uuid;