var Fs = require('fs-extra');

Fs.ensureFileSync('config.json');
if (/^\s*$/.test(Fs.readFileSync('config.json'))) {
  Fs.writeFileSync('config.json', '{}');
}

Config = require('../config.json');

module.export = Config;
