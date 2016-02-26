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
var Common = Rfr('tests/dep/common.js');
var Uuid = Rfr('lib/data/uuid.js');

var controllers = Common.controllers;
var collections = Common.collections;
var data = Common.data;

describe('Controller/Node', function () {
    describe('#create', function () {
        var validName = 'NewNode1';
        var validLocation = new Uuid('12346578-1234-4321-ABCD-1234567890AB');
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
                    controllers.node.create({asdf: 'asdf'}, validLocation, validIP, validPort, validAddresses);
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

            it('should fail when not UUID', function () {
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
                    controllers.node.create(validName, validLocation, {asdf: 'asdf'}, validPort, validAddresses);
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
                    controllers.node.create(validName, validLocation, validIP, validPort, {asdf: 'asdf'});
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

    describe('#get', function () {
        it('should fail when null', function () {
            Should.throws(function () {
                controllers.node.get(null);
            });
        });

        it('should fail when undefined', function () {
            Should.throws(function () {
                controllers.node.get(undefined);
            });
        });

        it('should fail when not UUID', function () {
            Should.throws(function () {
                controllers.node.get('asdf');
            });
        });

        it('should return undefined if no node', function () {
            var node = undefined;
            var validUUID = new Uuid('3cc22a4d-6f12-4ca1-9cc7-2cbd6b1ce893');

            Should.doesNotThrow(function () {
                node = controllers.node.get(validUUID);
            });

            Should.not.exist(node, 'Did not expect a node');
        });

        it('should return the node when defined', function () {
            var node = undefined;
            var validUUID = collections.node.db[0].getUUID();

            Should.doesNotThrow(function () {
                node = controllers.node.get(validUUID);
            });

            Should.exist(node, 'Expected a node');
        });
    });

    describe('#getAtLocation', function () {
        it('should fail when null', function () {
            Should.throws(function () {
                controllers.node.getAtLocation(null);
            });
        });

        it('should fail when undefined', function () {
            Should.throws(function () {
                controllers.node.getAtLocation(undefined);
            });
        });

        it('should fail when not UUID', function () {
            Should.throws(function () {
                controllers.node.getAtLocation('asdf');
            });
        });

        it('should fail if no location', function () {
            var validUUID = new Uuid('3cc22a4d-6f12-4ca1-9cc7-2cbd6b1ce893');

            Should.throws(function () {
                controllers.node.getAtLocation(validUUID);
            });
        });

        it('should return the nodes when defined', function () {
            var nodes = undefined;
            var validUUID = data.location[1].getUUID();

            Should.doesNotThrow(function () {
                nodes = controllers.node.getAtLocation(validUUID);
            });

            Should.equal(nodes.length, 2, 'Expected 2 nodes at the location');
        });
    });

    describe('#addAddress', function () {
        var validUUID = new Uuid('7d7d8d97-3a1c-4a8d-be5f-7577fac267af');
        var validIP = '192.168.1.10';
        var validPort = 25656;

        describe('when node is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.addAddress(null, validIP, validPort, false);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.addAddress(undefined, validIP, validPort, false);
                });
            });

            it('should fail when not UUID', function () {
                Should.throws(function () {
                    controllers.node.addAddress('asdfasdf', validIP, validPort, false);
                });
            });
        });

        describe('when node is valid', function () {
            describe('and ip is invalid', function () {
                it('should fail when null', function () {
                    Should.throws(function () {
                        controllers.node.addAddress(validUUID, null, validPort, false);
                    });
                });

                it('should fail when undefined', function () {
                    Should.throws(function () {
                        controllers.node.addAddress(validUUID, undefined, validPort, false);
                    });
                });

                it('should fail when not string', function () {
                    Should.throws(function () {
                        controllers.node.addAddress(validUUID, {asdf: 'asdf'}, validPort, false);
                    });
                });

                it('should fail when empty', function () {
                    Should.throws(function () {
                        controllers.node.addAddress(validUUID, '', validPort, false);
                    });
                });

                it('should fail when whitespace', function () {
                    Should.throws(function () {
                        controllers.node.addAddress(validUUID, ' ', validPort, false);
                    });
                });

                it('should fail when not IP', function () {
                    Should.throws(function () {
                        controllers.node.addAddress(validUUID, 'asdfasdf', validPort, false);
                    });
                });

                it('should fail when domain instead of IP', function () {
                    Should.throws(function () {
                        controllers.node.addAddress(validUUID, 'google.com', validPort, false);
                    });
                })
            });

            describe('and ip is valid', function () {
                describe('and  port is invalid', function () {
                    it('should fail when null', function () {
                        Should.throws(function () {
                            controllers.node.addAddress(validUUID, validIP, null, false);
                        });
                    });

                    it('should fail when undefined', function () {
                        Should.throws(function () {
                            controllers.node.addAddress(validUUID, validIP, undefined, false);
                        });
                    });

                    it('should fail when not a number', function () {
                        Should.throws(function () {
                            controllers.node.addAddress(validUUID, validIP, 'asdf', false);
                        });
                    });

                    it('should fail if outside port range', function () {
                        Should.throws(function () {
                            controllers.node.addAddress(validUUID, validIP, -10, false);
                        });
                    });
                });

                describe('and port is valid', function () {
                    it('should not error', function () {
                        //TODO: Test if address was correctly added
                        Should.doesNotThrow(function () {
                            controllers.node.addAddress(validUUID, validIP, validPort, false);
                        });
                    });
                });
            });
        });
    });

    describe('#removeAddress', function () {
        var validUUID = new Uuid('0f4fb36f-7795-4765-8268-94d102810fa5');
        var validIP = '10.0.0.2';
        var validPort = 25656;

        describe('when node is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.removeAddress(null, validIP, validPort, false);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.removeAddress(undefined, validIP, validPort, false);
                });
            });

            it('should fail when not UUID', function () {
                Should.throws(function () {
                    controllers.node.removeAddress('asdfasdf', validIP, validPort, false);
                });
            });
        });

        describe('when node is valid', function () {
            describe('and ip is invalid', function () {
                it('should fail when null', function () {
                    Should.throws(function () {
                        controllers.node.removeAddress(validUUID, null, validPort, false);
                    });
                });

                it('should fail when undefined', function () {
                    Should.throws(function () {
                        controllers.node.removeAddress(validUUID, undefined, validPort, false);
                    });
                });

                it('should fail when not string', function () {
                    Should.throws(function () {
                        controllers.node.removeAddress(validUUID, {asdf: 'asdf'}, validPort, false);
                    });
                });

                it('should fail when empty', function () {
                    Should.throws(function () {
                        controllers.node.removeAddress(validUUID, '', validPort, false);
                    });
                });

                it('should fail when whitespace', function () {
                    Should.throws(function () {
                        controllers.node.removeAddress(validUUID, ' ', validPort, false);
                    });
                });

                it('should fail when not IP', function () {
                    Should.throws(function () {
                        controllers.node.removeAddress(validUUID, 'asdfasdf', validPort, false);
                    });
                });

                it('should fail when domain instead of IP', function () {
                    Should.throws(function () {
                        controllers.node.removeAddress(validUUID, 'google.com', validPort, false);
                    });
                })
            });

            describe('and ip is valid', function () {
                describe('and  port is invalid', function () {
                    it('should fail when null', function () {
                        Should.throws(function () {
                            controllers.node.removeAddress(validUUID, validIP, null, false);
                        });
                    });

                    it('should fail when undefined', function () {
                        Should.throws(function () {
                            controllers.node.removeAddress(validUUID, validIP, undefined, false);
                        });
                    });

                    it('should fail when not a number', function () {
                        Should.throws(function () {
                            controllers.node.removeAddress(validUUID, validIP, 'asdf', false);
                        });
                    });

                    it('should fail if outside port range', function () {
                        Should.throws(function () {
                            controllers.node.removeAddress(validUUID, validIP, -10, false);
                        });
                    });
                });

                describe('and port is valid', function () {
                    it('should not error', function () {
                        //TODO: Test if remove correctly removed the address
                        Should.doesNotThrow(function () {
                            controllers.node.removeAddress(validUUID, validIP, validPort, false);
                        });
                    });
                });
            });
        });
    });

    describe('#updateAddress', function () {
        var validUUID = new Uuid('0f4fb36f-7795-4765-8268-94d102810fa5');
        var validIP = '10.0.0.2';
        var validPort = 25656;

        describe('when node is invalid', function () {
            it('should fail when null', function () {
                Should.throws(function () {
                    controllers.node.updateAddress(null, validIP, validPort, false);
                });
            });

            it('should fail when undefined', function () {
                Should.throws(function () {
                    controllers.node.updateAddress(undefined, validIP, validPort, false);
                });
            });

            it('should fail when not UUID', function () {
                Should.throws(function () {
                    controllers.node.updateAddress('asdfasdf', validIP, validPort, false);
                });
            });
        });

        describe('when node is valid', function () {
            describe('and ip is invalid', function () {
                it('should fail when null', function () {
                    Should.throws(function () {
                        controllers.node.updateAddress(validUUID, null, validPort, false);
                    });
                });

                it('should fail when undefined', function () {
                    Should.throws(function () {
                        controllers.node.updateAddress(validUUID, undefined, validPort, false);
                    });
                });

                it('should fail when not string', function () {
                    Should.throws(function () {
                        controllers.node.updateAddress(validUUID, {asdf: 'asdf'}, validPort, false);
                    });
                });

                it('should fail when empty', function () {
                    Should.throws(function () {
                        controllers.node.updateAddress(validUUID, '', validPort, false);
                    });
                });

                it('should fail when whitespace', function () {
                    Should.throws(function () {
                        controllers.node.updateAddress(validUUID, ' ', validPort, false);
                    });
                });

                it('should fail when not IP', function () {
                    Should.throws(function () {
                        controllers.node.updateAddress(validUUID, 'asdfasdf', validPort, false);
                    });
                });

                it('should fail when domain instead of IP', function () {
                    Should.throws(function () {
                        controllers.node.updateAddress(validUUID, 'google.com', validPort, false);
                    });
                })
            });

            describe('and ip is valid', function () {
                describe('and  port is invalid', function () {
                    it('should fail when null', function () {
                        Should.throws(function () {
                            controllers.node.updateAddress(validUUID, validIP, null, false);
                        });
                    });

                    it('should fail when undefined', function () {
                        Should.throws(function () {
                            controllers.node.updateAddress(validUUID, validIP, undefined, false);
                        });
                    });

                    it('should fail when not a number', function () {
                        Should.throws(function () {
                            controllers.node.updateAddress(validUUID, validIP, 'asdf', false);
                        });
                    });

                    it('should fail if outside port range', function () {
                        Should.throws(function () {
                            controllers.node.updateAddress(validUUID, validIP, -10, false);
                        });
                    });
                });

                describe('and port is valid', function () {
                    it('should not error', function () {
                        //TODO: Test if address updated correctly
                        Should.doesNotThrow(function () {
                            controllers.node.updateAddress(validUUID, validIP, validPort, false);
                        });
                    });
                });
            });
        });
    });
});