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
var Node = Rfr('lib/model/node.js');
var Common = Rfr('tests/common.js');

var controllers = Common.controllers;
var collections = Common.collections;
var data = Common.data;

describe('Controller/Node', function () {
    describe('#create', function () {
        var validName = 'NewNode1';
        var validLocation = '12346578-1234-4321-ABCD-1234567890AB';
        var validIP = '192.168.1.1';
        var validPort = 5656;
        var validAddresses = [{
            ip: '192.168.1.1',
            port: 25565
        }];

        describe('when name is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.create(null, validLocation, validIP, validPort, validAddresses);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.create(undefined, validLocation, validIP, validPort, validAddresses);
                });
            });

            it('should fail when empty', function () {
                Should.throws(function () {
                    controllers.node.create('', validLocation, validIP, validPort, validAddresses);
                });
            });

            it('should fail when whitespace', function () {
                Should.throws(function () {
                    controllers.node.create(' ', validLocation, validIP, validPort, validAddresses);
                });
            });

            it('should fail when not a string', function () {
                Should.throws(function () {
                    controllers.node.create({ asdf: 'asdf' }, validLocation, validIP, validPort, validAddresses);
                });
            });
        });

        describe('when location is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.create(validName, null, validIP, validPort, validAddresses);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.create(validName, undefined, validIP, validPort, validAddresses);
                });
            });

            it('should fail when empty', function () {
                Should.throws(function () {
                    controllers.node.create(validName, '', validIP, validPort, validAddresses);
                });
            });

            it('should fail when whitespace', function () {
                Should.throws(function () {
                    controllers.node.create(validName, ' ', validIP, validPort, validAddresses);
                });
            });

            it('should fail when not a string', function () {
                Should.throws(function () {
                    controllers.node.create(validName, { asdf: 'asdf' }, validIP, validPort, validAddresses);
                });
            });

            it('should fail when string but not UUID', function () {
                Should.throws(function () {
                    controllers.node.create(validName, 'asdf', validIP, validPort, validAddresses);
                });
            });
        });

        describe('when ip is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, null, validPort, validAddresses);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, undefined, validPort, validAddresses);
                });
            });

            it('should fail when empty', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, '', validPort, validAddresses);
                });
            });

            it('should fail when whitespace', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, ' ', validPort, validAddresses);
                });
            });

            it('should fail when not a string', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, { asdf: 'asdf' }, validPort, validAddresses);
                });
            });

            it('should fail when string but not IP', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, 'asdf', validPort, validAddresses);
                });
            });
        });

        describe('when port is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, validIP, null, validAddresses);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, validIP, undefined, validAddresses);
                });
            });

            it('should fail when not a number', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, validIP, 'asdf', validAddresses);
                });
            });

            it('should fail if outside port range', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, validIP, -10, validAddresses);
                });
            });
        });

        describe('when addresses is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, validIP, validPort, null);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, validIP, validPort, undefined);
                });
            });

            it('should fail when not array', function () {
                Should.throws(function () {
                    controllers.node.create(validName, validLocation, validIP, validPort, { asdf: 'asdf' });
                });
            });
        });

        describe('when data is valid', function () {
            it('should create node', function () {
                var oldCount = collections.node.db.length;

                Should.doesNotThrow(function () {
                    controllers.node.create(validName, validLocation, validIP, validPort, validAddresses);
                });

                Should.equal(oldCount + 1, collections.node.db.length, 'Expected node count to increase by one');
            });
        });
    });
});