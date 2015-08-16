/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
require('date-utils')
var fs = require('fs-extra')
var Winston = require('winston')
var CurrentTime = new Date()
var Log = {
  debugStatus: false,
  verboseStatus: false
}

// Make sure that the logs directory exists
fs.ensureDir('./logs', function (err) {
  if (err != null) {
    console.log(err)
  }
})

var Logger = new (Winston.Logger)({
  transports: [
    new (Winston.transports.File)({
      filename: './logs/' + CurrentTime.toYMD('.') + '-' + CurrentTime.toFormat('HH24.MI.SS') + '.log',
      level: 'info',
      json: false,
      timestamp: function () {
        var CurrentTime = new Date()
        return CurrentTime.toFormat('HH24:MI:SS')
      }
    })
  ]
})

Log.debug = function (message, data) {

  if (this.debugStatus === true) {

    Logger.transports.file.level = 'debug'
    Logger.debug(message, { meta: data })

    var CurrentTime = new Date()
    console.log('[' + CurrentTime.toFormat('HH24:MI:SS') + '][DEBUG] ' + message)

  }

}

Log.verbose = function (message, data) {

  if (this.verboseStatus === true || this.debugStatus === true) {

    Logger.transports.file.level = 'verbose'
    Logger.verbose(message, { meta: data })

    var CurrentTime = new Date()
    console.log('[' + CurrentTime.toFormat('HH24:MI:SS') + '][VERBOSE] ' + message)

  }

}

Log.info = function (message, data) {

  Logger.transports.file.level = 'info'
  Logger.info(message, { meta: data })

  var CurrentTime = new Date()
  console.log('[' + CurrentTime.toFormat('HH24:MI:SS') + '][INFO] ' + message)

}

Log.warn = function (message, data) {

  Logger.transports.file.level = 'warn'
  Logger.warn(message, { meta: data })

  var CurrentTime = new Date()
  console.log('[' + CurrentTime.toFormat('HH24:MI:SS') + ']\x1b[30m\x1b[43m[WARN]\x1b[0m ' + message)

}

Log.error = function (message, data) {

  Logger.transports.file.level = 'error'
  Logger.error(message, { meta: data })

  var CurrentTime = new Date()
  console.log('[' + CurrentTime.toFormat('HH24:MI:SS') + ']\x1b[1m\x1b[37m\x1b[41m[ERROR]\x1b[0m ' + message)

}

module.exports = Log
