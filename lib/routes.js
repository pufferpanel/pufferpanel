/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Logger = require('./logger.js');
var Hapi = require('hapi');
var Async = require('async');
var Nunjucks = require('nunjucks');
var server = new Hapi.Server();
var Path = require('path');
var Configuration = require(Path.join(__dirname, '../configuration.json'));
var Auth = require(Path.join(__dirname, '../server/auth/authentication.js'));
var Authentication = new Auth();
var Routes = require(Path.join(__dirname, 'routes/core.js'));

server.connection({ port: Configuration.server.port });

server.register([
  {
    register: require('inert')
  },
  {
    register: require('yar'),
    options: {
      cookieOptions: {
        password: Configuration.yarPassword,
        isSecure: false,
        isHttpOnly: true
      }
    }
  },
  {
    register: require('hapi-auth-cookie')
  },
  {
    register: require('vision')
  }
], function (err) {

  Logger.info('Testing');
  if (err) {
    Logger.error('Failed to load one or more required server plugins.', err);
  }

  // Hapi Auth Cookie
  server.auth.strategy('session', 'cookie', {
    password: Configuration.yarPassword,
    cookie: 'pp_session',
    redirectTo: '/auth/login',
    isSecure: false
  });

  // Hapi Vision
  server.views({
    engines: {
      html: {
        compile: function (src, options) {
          var template = Nunjucks.compile(src, options.environment);
          return function (context) {
            return template.render(context);
          };
        },
        prepare: function (options, next) {
          options.compileOptions.environment = Nunjucks.configure(options.path, { watch: false });
          return next();
        }
      }
    },
    path: Path.join(__dirname, '../app/views')
  });

});

// Handle public folder
server.route({
  method: 'GET',
  path: '/public/{param*}',
  handler: {
    directory: {
      path: 'public'
    }
  }
});

server.route([
  {
    method: 'GET',
    path: '/index',
    config: {
      handler: Routes.Base.get.index,
      auth: 'session'
    }
  }
]);

// Routes for Authentication
server.route([
  {
    method: 'GET',
    path: '/auth/login',
    config: {
      handler: Routes.Authentication.get.login,
      auth: {
        mode: 'try',
        strategy: 'session'
      },
      plugins: {
        'hapi-auth-cookie': {
          redirectTo: false
        }
      }
    }
  },
  {
    method: 'POST',
    path: '/auth/login',
    config: {
      handler: Routes.Authentication.post.login,
      auth: {
        mode: 'try',
        strategy: 'session'
      },
      plugins: {
        'hapi-auth-cookie': {
          redirectTo: false
        }
      }
    }
  },
  {
    method: 'POST',
    path: '/auth/login/totp',
    config: {
      handler: Routes.Authentication.post.totp,
      auth: {
        mode: 'try',
        strategy: 'session'
      },
      plugins: {
        'hapi-auth-cookie': {
          redirectTo: false
        }
      }
    }
  }
]);

server.start(function () {
  Logger.info('Server running at: ' + server.info.uri);
});
