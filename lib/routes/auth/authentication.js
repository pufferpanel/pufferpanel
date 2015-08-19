/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Path = require('path');
var Async = require('async');
var Logger = require(Path.join(__dirname, '../../logger.js'));
var Auth = require(Path.join(__dirname, '../../../server/auth/authentication.js'));
var Authentication = new Auth();

var Routes = {
  post: {
    login: function (request, reply) {
      Async.series([
        function (callback) {
          Authentication.validateCredentials(request, callback);
        }
      ], function (data, result, err) {

        if (err !== null) {
          Logger.error(err);
        }

        if (result[0] === true) {
          request.auth.session.set(data);
          reply.redirect('/index');
        } else {
          request.session.flash('loginError', data, true);
          reply.redirect('/auth/login');
        }

      });
    },
    totp: function (request, reply) {
      Async.series([
        function (callback) {
          Authentication.TOTPEnabled(request.payload.check, callback);
        }
      ], function (err, results) {

        if (err) {
          Logger.error('An unhandled error occured during async.series in the login TOTP handler.', err);
          return false;
        }

        if (results[0] === 1) {
          reply('true');
        } else {
          reply('false');
        }

      });
    },
    register: function (request, reply) {

    },
    password: function (request, reply) {

    }
  },
  get: {
    login: function (request, reply) {
      reply.view('auth/login', {
        flash: request.session.flash('loginError')
      });
    },
    logout: function (request, reply) {
      request.auth.session.clear();
      reply.redirect('/auth/login');
    },
    register: function (request, reply) {
      reply.view('auth/register', {
        flash: request.session.flash('registerError'),
        token: request.params.token
      });
    },
    password: function (request, reply) {
      reply.view('auth/password', {
        flash: request.session.flash('passwordError'),
        noshow: false
      });
    }
  }
};

module.exports = Routes;
