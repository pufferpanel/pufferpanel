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

      if (Authentication.validateCredentials(request.payload.email, request.payload.password)) {
        request.session.set('pp_authenticated', 'token-value')
      } else {
        request.session.flash('loginError', 'Sorry, your email or password was invalid.', true)
      }

      reply.redirect('/auth/login').temporary(true)

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
  console.log('Server running at: ', server.info.uri)
})
