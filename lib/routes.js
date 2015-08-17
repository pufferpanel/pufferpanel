/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */

var Logger = require('./logger.js')
var Hapi = require('hapi')
var Async = require('async')
var Nunjucks = require('nunjucks')
var server = new Hapi.Server()
var Path = require('path')
var Configuration = require(Path.join(__dirname, '../configuration.json'))
var Auth = require(Path.join(__dirname, '../server/front/authentication.js'))
var Authentication = new Auth()

server.connection({ port: Configuration.server.port })

server.register(require('inert'), function (err) {
  if (err) {
    Logger.error('Failed to register \'inert\' plugin in server.')
  }
})

server.register({
  register: require('yar'),
  options: {
    cookieOptions: {
      password: Configuration.yarPassword,
      isSecure: false,
      isHttpOnly: true
    }
  }
}, function (err) {
  if (err) {
    Logger.error('Failed to register \'yar\' plugin in server.')
  }
})

server.register(require('vision'), function (err) {

  if (err) {
    Logger.error('Failed to register \'vision\' plugin in server.')
  }

  server.views({
    engines: {
      html: {
        compile: function (src, options) {
          var template = Nunjucks.compile(src, options.environment)
          return function (context) {
            return template.render(context)
          }
        },
        prepare: function (options, next) {
          options.compileOptions.environment = Nunjucks.configure(options.path, { watch: false })
          return next()
        }
      }
    },
    path: Path.join(__dirname, '../app/views')
  })

})

// Handle public folder
server.route({
  method: 'GET',
  path: '/public/{param*}',
  handler: {
    directory: {
      path: 'public'
    }
  }
})

// Routes for Authentication
server.route([
  {
    method: 'GET',
    path: '/auth/login',
    handler: function (request, reply) {
      reply.view('auth/login', {
        flash: request.session.flash('loginError')
      })
    }
  },
  {
    method: 'POST',
    path: '/auth/login',
    handler: function (request, reply) {

      Async.series([
        function (callback) {
          Authentication.validateCredentials(request.info.remoteAddress, request.payload, callback)
        }
      ], function (data, result) {

        if (result[0] === true) {
          request.session.set('pp_session', data)
          reply.redirect('/index').temporary(true)
        } else {
          request.session.flash('loginError', data, true)
          reply.redirect('/auth/login').temporary(true)
        }

      })

    }
  },
  {
    method: 'POST',
    path: '/auth/login/totp',
    handler: function (request, reply) {

      Async.series([
        function (callback) {
          Authentication.TOTPEnabled(request.payload.check, callback)
        }
      ], function (err, results) {

        if (err) {
          Logger.error('An unhandled error occured during async.series in the login TOTP handler.', err)
          return false
        }

        if (results[0] === 1) {
          reply('true')
        } else {
          reply('false')
        }

      })

    }
  },
  {
    method: 'GET',
    path: '/auth/register',
    handler: function (request, reply) {
      reply.view('auth/register')
    }
  },
  {
    method: 'GET',
    path: '/auth/reset-password',
    handler: function (request, reply) {
      reply.view('auth/reset-password')
    }
  }
])

server.start(function () {
  Logger.info('Server running at: ' + server.info.uri)
})
