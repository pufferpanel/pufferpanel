var Assert = require('assert');
var Rfr = require('rfr');
var Location = Rfr('lib/model/location.js');
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

var data = [
    new Location('12346578-1234-4321-ABCD-1234567890AB', 'Test1'),
    new Location('87654321-4321-4321-8765-BA0123456789', 'Test2')
];

beforeEach(function () {
    collections.location._reset(data);
});

module.exports.collections = collections;
module.exports.controllers = controllers;
module.exports.data = data;