var _ = require('underscore');

var ServersApi = {};

var servers = [];
var defaultRouter = [];

defaultRouter[''] = {
    method: 'GET',
    handler: function (request, response, server) {

        response.view('server/index.html');
    }
};
defaultRouter['index'] = {
    method: 'GET',
    handler: function (request, response, server) {

        response.view('server/index.html');
    }
};

ServersApi._getRegisteredTypes = function () {

    return _.keys(servers);
};

ServersApi._getRoutes = function (type) {

    return servers[type];
};

ServersApi._parseRoute = function (type, route) {

    if (route === undefined) {
        route = '';
    }
    var routes = servers[type].router;
    return routes[route];
};

/**
 * Registers a new server type, assigning default routes.
 *
 * @param {String} type - Server type
 * @param {Object?} router - Default routes
 */
ServersApi.registerType = function (type, router) {

    if (router === undefined) {
        router = {};
    }
    var serverType = type.toLowerCase();
    if (_.contains(servers, serverType)) {
        throw new Error('Cannot register the same type twice (' + serverType + ')');
    }
    servers[serverType] = {
        router: _.extend(ServersApi._createDefaultRouteClone(), router)
    };
};

/**
 * Registers routing information for a given server type.
 * This will extend the existing routing structure, but will override
 * any duplicate routes.
 *
 * @param {String} type - Server type
 * @param {Object} router - Route object
 */
ServersApi.registerRoute = function (type, router) {

    var serverType = type.toLowerCase();
    var server = servers[serverType];
    if (router === undefined) {
        throw new Error('Tried to register router for unknown type (' + serverType + ')');
    }
    server.router = _.extend(server.router, router);
};

ServersApi._createDefaultRouteClone = function () {

    var cloned = [];
    _.forEach(_.keys(defaultRouter), function (e) {
        cloned[e] = _.clone(defaultRouter[e]);
    });
    return cloned;
};

module.exports = ServersApi;
