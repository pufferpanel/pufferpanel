var Assert = require('assert');
var Rfr = require('rfr');
var Location = Rfr('lib/model/location.js');
var Node = Rfr('lib/model/node.js');
var LocationCollection = Rfr('tests/dep/collection/location.js');
var NodeCollection = Rfr('tests/dep/collection/node.js');
var UserController = Rfr('lib/controller/user.js');
var ServerController = Rfr('lib/controller/server.js');
var LocationController = Rfr('lib/controller/location.js');
var NodeController = Rfr('lib/controller/node.js');

var collections = {
    location: new LocationCollection(),
    node: new NodeCollection()
};

var controllers = {
    location: new LocationController(collections),
    user: new UserController(collections),
    server: new ServerController(collections),
    node: new NodeController(collections)
};

var data = {
    location: [
        new Location('12346578-1234-4321-ABCD-1234567890AB', 'Test1'),
        new Location('87654321-4321-4321-8765-BA0123456789', 'Test2')
    ],
    node: [
        new Node('7d7d8d97-3a1c-4a8d-be5f-7577fac267af', 'Node1', '12346578-1234-4321-ABCD-1234567890AB', '127.0.0.1', 5656, [{
            ip: '127.0.0.1',
            port: 25565,
            allocated: false
        }]),
        new Node('1d6e94c2-b6ff-411f-86a8-625fb4dc7a78', 'Node2', '87654321-4321-4321-8765-BA0123456789', '127.0.0.1', 5656, [{
            ip: '127.0.0.1',
            port: 25566,
            allocated: false
        }]),
        new Node('0f4fb36f-7795-4765-8268-94d102810fa5', 'Node3', '87654321-4321-4321-8765-BA0123456789', '10.0.0.1', 5656, [{
            ip: '10.0.0.2',
            port: 25566,
            allocated: false,
            ip: '10.0.0.3',
            port: 25565,
            allocated: true
        }])
    ]
};

beforeEach(function () {
    collections.location._reset(data.location);
    collections.node._reset(data.node);
});

module.exports.collections = collections;
module.exports.controllers = controllers;
module.exports.data = data;