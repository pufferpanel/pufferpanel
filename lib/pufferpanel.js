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
var Rethink = require('rethinkdbdash');
var UserCollection = Rfr('lib/collection/user.js');
var UserController = Rfr('lib/controller/user.js');
var ServerCollection = Rfr('lib/collection/server.js');
var ServerController = Rfr('lib/controller/server.js');
var LocationCollection = Rfr('lib/collection/location.js');
var LocationController = Rfr('lib/controller/location.js');
var NodeCollection = Rfr('lib/collection/node.js');
var NodeController = Rfr('lib/controller/node.js');
var Web = Rfr('lib/web/server.js');
var Logger = Rfr('lib/logger/logger.js');
var ModuleLoader = Rfr('lib/modules/loader.js');

var PufferPanel = {};

PufferPanel.logger = Logger;

var dbConn = new Rethink({});

PufferPanel.collections = {
    user: new UserCollection(dbConn),
    location: new LocationCollection(dbConn),
    node: new NodeCollection(dbConn),
    server: new ServerCollection(dbConn)
};

PufferPanel.controllers = {
    user: new UserController(PufferPanel.collections.user),
    server: new ServerController(PufferPanel.collections.server)
};

PufferPanel.routes = {};

PufferPanel.web = Web;

PufferPanel.modules = ModuleLoader;

module.exports = PufferPanel;

PufferPanel.web.initialize();

PufferPanel.web.start();