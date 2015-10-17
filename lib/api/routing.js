var Path = require('path');

var RoutingApi = {};

var routes = [];
var views = [
    Path.join(__dirname, '../../app/views')
];

/**
 * Registers a routing path using an object definition
 * @param {Object} route
 */
RoutingApi.register = function (route) {

    routes[route.path] = {
        method: route.method,
        callback: route.handler,
        auth: route.auth
    };
};

/**
 *  Registers a folder to the views path
 *
 * @param {String} path - Folder path to register
 */
RoutingApi.registerView = function (path) {

    views.unshift(path);
};

RoutingApi._getRoutes = function () {

    return routes;
};

RoutingApi._getViews = function () {

    return views;
};

module.exports = RoutingApi;
