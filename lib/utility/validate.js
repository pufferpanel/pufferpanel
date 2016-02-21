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

var ArgValidate = require('argument-validator');
var StrValidate = require('validator');
var Rfr = require('rfr');
var Uuid = Rfr('lib/data/uuid.js');

var Validate = {};

/**
 * Validate if a value is a String
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not a string
 */
Validate.isString = function (value, name) {
    Validate.isNotNull(value, name);
    ArgValidate.string(value, name);
};

/**
 * Validate if a value is a Boolean
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not a boolean
 */
Validate.isBoolean = function (value, name) {
    Validate.isNotNull(value, name);
    ArgValidate.boolean(value, name);
};

/**
 * Validate if a value is a Date
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not a date
 */
Validate.isDate = function (value, name) {
    Validate.isNotNull(value, name);
    ArgValidate.isDate(value, name);
};

/**
 * Validate if a value is a valid email address
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not an email
 */
Validate.isEmail = function (value, name) {
    Validate.isNotNull(value, name);
    Validate.isString(value, name);
    if (!StrValidate.isEmail(value)) {
        throw new Error(name + ' is not a valid email', { data: value });
    }
};

/**
 * Validate if a value is an IP
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not an IP
 */
Validate.isIP = function (value, name) {
    Validate.isNotNull(value, name);
    Validate.isString(value, name);
    if (!StrValidate.isIP(value)) {
        throw new Error(name + ' is not a valid IP', { data: value });
    }
};

/**
 * Validate if a value is a FQDN
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not a FQDN
 */
Validate.isFQDN = function (value, name) {
    Validate.isNotNull(value, name);
    Validate.isString(value, name);
    if (!StrValidate.isFQDN(value)) {
        throw new Error(name + ' is not a valid FQDN', { data: value });
    }
};

/**
 * Validate if a value is a Number
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not a number
 */
Validate.isNumber = function (value, name) {
    Validate.isNotNull(value, name);
    ArgValidate.number(value, name);
};

/**
 * Validate if a value is a UUID
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not a UUID
 */
Validate.isUUID = function (value, name) {
    Validate.isNotNull(value, name);

    if (value.constructor.name == 'Uuid') {
        return;
    }

    Validate.isString(value, name);
    if (!StrValidate.isUUID(value, 4)) {
        throw new Error(name + ' is not a valid UUID', { data: value });
    }
};

/**
 * Validate if a value is not null
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is null or undefined
 */
Validate.isNotNull = function (value, name) {
    ArgValidate.notNull(value, name);
};

/**
 * Validate if a value is an array
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not an array
 */
Validate.isArray = function (value, name) {
    ArgValidate.array(value, name);
};

/**
 * Validate if a value is an Object
 *
 * @param {*} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not an Object
 */
Validate.isObject = function (value, name) {
    ArgValidate.object(value, name);
};

/**
 * Validate if a value is a Number and within 0-65535
 *
 * @param {String} value - Value to test
 * @param {String} name - Name of value
 * @throws Will throw if value is not a Number or is outside the port range
 */
Validate.isPort = function (value, name) {
    Validate.isNumber(value, name);
    if (value <= 0 || value > 65535) {
        throw new Error('Port is outside range of 0-65535');
    }
};

/**
 * Validate if a value is a specific type
 *
 * @param {*} value - Value to test
 * @param {String} type - Type value should be
 * @param {String} name - Name of value
 * @throws Will throw if value is not the specified type
 */
Validate.isType = function (value, type, name) {
    console.log(Object.prototype.toString.call(value));
    if (typeof value !== type) {
        throw new Error('Expected ' + type + ', found ' + (typeof value) + ' for ' + name);
    }
};

module.exports = Validate;