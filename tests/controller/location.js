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
var LocationCollection = Rfr('tests/dep/collection/location.js');
var UserController = Rfr('lib/controller/user.js');
var ServerController = Rfr('lib/controller/server.js');
var LocationController = Rfr('lib/controller/location.js');
var NodeController = Rfr('lib/controller/node.js');

var collections = {
    location: new LocationCollection()
};

var controllers = {
    location: new LocationController(collections),
    user: new UserController(collections),
    server: new ServerController(collections),
    node: new NodeController(collections)
};

beforeEach(function () {
    collections.location._reset();
});

describe('Location Controller', function () {
    describe('#create', function () {
        it('should error when username is null', function () {
            Should.throws(function () {
                controllers.location.create(null)
            });
        });

        it('should error when username is undefined', function () {
            Should.throws(function () {
                controllers.location.create(undefined);
            });
        });

        it('should error when username is empty', function () {
            Should.throws(function () {
                controllers.location.create('');
            });
        });

        it('should error when username is whitespace', function () {
            Should.throws(function () {
                controllers.location.create(' ');
            });
        });

        it('should error when username is not a string', function () {
            Should.throws(function () {
                controllers.location.create({ asdf: 'notavalue' });
            });
        });

        it('should not error when username is valid', function () {
            var location = undefined;

            Should.doesNotThrow(function () {
                location = controllers.location.create('asdf');
            }, 'Expected location to not error');

            Should.exist(location, 'Expected location to exist');
            Should.exist(controllers.location.get(location.getUUID()), 'Expected location to exist in controller');
        });
    });
});