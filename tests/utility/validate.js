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

var Assert = require('assert');
var Should = require('should');
var Rfr = require('rfr');
var Validate = Rfr('lib/utility/validate.js');

describe('Utilty/Validate', function () {
    describe('isString', function () {
        it('should error if null', function () {
            Should.throws(function () {
                Validate.isString(null, 'asdf');
            });
        });

        it('should error if undefined', function () {
            Should.throws(function () {
                Validate.isString(undefined, 'asdf');
            });
        });

        it('should error if not string', function () {
            Should.throws(function () {
                Validate.isString({ asdf: 'asdf' }, 'asdf');
            });
        });

        it('should error if empty', function () {
            Should.throws(function () {
                Validate.isString('', 'asdf');
            });
        });

        it('should error if whitespace', function () {
            Should.throws(function () {
                Validate.isString(' ', 'asdf');
            });
        });

        it('should not error if string', function () {
            Should.doesNotThrow(function () {
                Validate.isString('asdf', 'asdf');
            });
        });
    });
});