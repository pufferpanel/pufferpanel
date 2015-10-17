var Rfr = require('rfr');
var ServerApi = Rfr('lib/api/servers.js');

ServerApi.registerType('minecraft', {
    '': {
        method: 'GET',
        handler: function (request, response, server) {

            response.view('minecraft/index.html');
        }
    },
    'index': {
        method: 'GET',
        handler: function (request, response, server) {

            response.view('minecraft/index.html');
        }
    }
});