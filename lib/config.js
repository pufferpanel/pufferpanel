var Fs = require('fs');
var FsE = require('fs-extra');

FsE.ensureFileSync('config.json');
if (Fs.readFileSync('config.json') == '') {
  Fs.writeFileSync('config.json', '{}');
}

Config = require('../config.json');

module.export = Config;