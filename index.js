/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Path = require('path')
var Yargs = require('yargs').argv
var Logger = require('./lib/logger.js')

// Setup Logging Information
if (Yargs.debug) {
  Logger.debugStatus = true
  Logger.debug('PufferPanel running in debug mode.')
}

if (Yargs.verbose) {
  Logger.verboseStatus = true
  Logger.verbose('PufferPanel running in verbose mode.')
}

// Include HapiJS Routing Mechanisms
require(Path.join(__dirname, 'lib/routes.js'))
