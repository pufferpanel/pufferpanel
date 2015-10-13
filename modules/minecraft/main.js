var Rfr = require('rfr');
var Path = require('path');
var ServerApi = Rfr('lib/api/servers.js');
var RoutingApi = Rfr('lib/api/routing.js');

RoutingApi.registerView(Path.join(__dirname, 'views'));
ServerApi.registerType('minecraft', {
    'index': {
        method: 'GET',
        handler: function (request, response, server) {
            response.view('minecraft/index.html');
        }
    }
});
