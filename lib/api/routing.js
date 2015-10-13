var Path = require('path');

var RoutingApi = {};

var routes = [];
var views = [
    Path.join(__dirname, '../../app/views')
];

/**
 * Registers a routing path with the given callback.
 *
 * @param {String|Array} method - HTTP methods to respond to
 * @param path - Route path
 * @param {callback} call - Callback to call to handle route execution
 */
RoutingApi.register = function (method, path, call, auth) {

    routes[path] = {
        method: method,
        callback: call,
        auth: auth
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
