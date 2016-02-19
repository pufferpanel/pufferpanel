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
var Location = Rfr('lib/model/location.js');
var LocationCollection = Rfr('tests/dep/collection/location.js');
var UserController = Rfr('lib/controller/user.js');
var ServerController = Rfr('lib/controller/server.js');
var LocationController = Rfr('lib/controller/location.js');
var NodeController = Rfr('lib/controller/node.js');
var Common = Rfr('tests/common.js');

var controllers = Common.controllers;
var collections = Common.collections;
var data = Common.data;

describe('Controller/Location', function () {
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
            var oldCount = collections.location.db.length;

            Should.doesNotThrow(function () {
                location = controllers.location.create('asdf');
            }, 'Expected location to not error');

            Should.exist(location, 'Expected location to exist');
            Should.exist(controllers.location.get(location.getUUID()), 'Expected location to exist in controller');
            Should.equal(oldCount + 1, collections.location.db.length, 'Expected collection to increase by one');
        });
    });

    describe('#changeName', function () {
        var validName = 'Test3';
        var validUUID = data[0].getUUID();

        describe('when uuid is invalid', function () {
            it('should error when uuid is null', function () {
                Should.throws(function () {
                    controllers.location.changeName(null, validName);
                });
            });

            it('should error when uuid is undefined', function () {
                Should.throws(function () {
                    controllers.location.changeName(undefined, validName);
                });
            });

            it('should error when uuid is empty', function () {
                Should.throws(function () {
                    controllers.location.changeName('', validName);
                });
            });

            it('should error when uuid is whitespace', function () {
                Should.throws(function () {
                    controllers.location.changeName(' ', validName);
                });
            });

            it('should error when uuid is not a string', function () {
                Should.throws(function () {
                    controllers.location.changeName({ asdf: 'asdf' }, validName);
                });
            });

            it('should error when uuid is not a valid UUID', function () {
                Should.throws(function () {
                    controllers.location.changeName('123-123-123', validName);
                });
            });
        });

        describe('when uuid is valid', function () {
            it('should error when name is null', function () {
                Should.throws(function () {
                    controllers.location.changeName(validUUID, null);
                });
            });

            it('should error when name is undefined', function () {
                Should.throws(function () {
                    controllers.location.changeName(validUUID, undefined);
                });
            });

            it('should error when name is empty', function () {
                Should.throws(function () {
                    controllers.location.changeName(validUUID, '');
                });
            });

            it('should error when name is whitespace', function () {
                Should.throws(function () {
                    controllers.location.changeName(validUUID, ' ');
                });
            });

            it('should error when name is not a string', function () {
                Should.throws(function () {
                    controllers.location.changeName(validUUID, { asdf: "asdf" });
                });
            });

            it('should not error when name is valid', function () {
                var newName = 'NewName';

                Should.doesNotThrow(function () {
                    location = controllers.location.changeName(validUUID, newName);
                }, 'Expected change to not error');

                var location = controllers.location.get(validUUID);
                Should.exist(location, 'Expected location to exist');
                Should.exist(controllers.location.get(location.getUUID()), 'Expected location to exist in controller');
                Should.equal(newName, location.getName(), 'Expected location name to have updated');
            });
        });
    });

    describe("#get", function () {
        it('should error when UUID is null', function () {
            Should.throws(function () {
                controllers.location.get(null);
            });
        });

        it('should error when UUID is undefined', function () {
            Should.throws(function () {
                controllers.location.get(undefined);
            });
        });

        it('should error when UUID is empty', function () {
            Should.throws(function () {
                controllers.location.get('');
            });
        });

        it('should error when UUID is whitespace', function () {
            Should.throws(function () {
                controllers.location.get(' ');
            });
        });

        it('should error when UUID is not a string', function () {
            Should.throws(function () {
                controllers.location.get({ asdf: "asdf" });
            });
        });

        it('should not error when UUID is valid', function () {
            var location = undefined;
            var validUUID = data[0].getUUID();

            Should.doesNotThrow(function () {
                location = controllers.location.get(validUUID);
            }, 'Expected change to not error');

            Should.exist(location, 'Expected location to exist');
        });

        it('should return undefined when UUID is valid and no location exists', function () {
            var location = undefined;
            var validButNotExistUUID = '12345678-1234-4531-ABCD-111111111111';

            Should.doesNotThrow(function () {
                location = controllers.location.get(validButNotExistUUID);
            }, 'Expected change to not error');

            Should.not.exist(location, 'Expected location to exist');
        });
    });

    describe("#delete", function () {
        it('should error when UUID is null', function () {
            Should.throws(function () {
                controllers.location.delete(null);
            });
        });

        it('should error when UUID is undefined', function () {
            Should.throws(function () {
                controllers.location.delete(undefined);
            });
        });

        it('should error when UUID is empty', function () {
            Should.throws(function () {
                controllers.location.delete('');
            });
        });

        it('should error when UUID is whitespace', function () {
            Should.throws(function () {
                controllers.location.delete(' ');
            });
        });

        it('should error when UUID is not a string', function () {
            Should.throws(function () {
                controllers.location.delete({ asdf: "asdf" });
            });
        });

        it('should not error when UUID is valid', function () {
            var oldCount = collections.location.db.length;
            var validUUID = data[0].getUUID();

            Should.doesNotThrow(function () {
                controllers.location.delete(validUUID);
            }, 'Expected delete to not error');

            var location = controllers.location.get(validUUID);
            Should.not.exist(location, 'Expected location to not exist');
            Should.equal(oldCount - 1, collections.location.db.length, 'Expected collection to decrease by one');
        });

        it('should not error when UUID is valid and no location exists', function () {
            var validUUID = '98765432-1234-4851-ADC1-ABC123456789';
            var oldCount = collections.location.db.length;
            var location = {};

            Should.doesNotThrow(function () {
                controllers.location.delete(validUUID);
            }, 'Expected delete to not error');

            var location = controllers.location.get(validUUID);
            Should.not.exist(location, 'Expected location to not exist');
            Should.equal(oldCount, collections.location.db.length, 'Expected collection to not change');
        });
    });
});