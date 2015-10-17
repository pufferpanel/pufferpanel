var Rfr = require('rfr');
var Fs = require('fs-extra');
var _ = require('underscore');
var Path = require('path');
var Logger = Rfr('lib/api/logger.js');
var RoutingApi = Rfr('lib/api/routing.js');

var ApiLoader = {};

var moduleFolder = Path.join(__dirname, '..', 'modules');
ApiLoader.load = function () {

    var folders = Fs.readdirSync(moduleFolder);

    _.forEach(folders, function (e) {

        var name = e;
        var version = 'unknown';
        var path = Path.join(moduleFolder, e);
        try {
            Logger.debug('Preparing module ' + e);
            var moduleJson = Rfr('modules/' + e + '/module.json');
            name = moduleJson.name || e;
            version = moduleJson.version || '???';
            Logger.info('Loading module ' + name + ' (v' + version + ')');
            if (moduleJson.views !== undefined) {
                Logger.debug('Registering views folder: ' + Path.join(path, moduleJson.views));
                RoutingApi.registerView(Path.join(path, moduleJson.views));
            }
            Rfr('modules/' + e + '/' + (moduleJson.main || 'main.js'));
        } catch (error) {
            Logger.error('Error loading module (' + name + ', v' + version + ')', error);
        }
    });
};

module.exports = ApiLoader;
