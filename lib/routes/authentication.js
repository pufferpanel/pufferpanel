/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Authentication = Rfr('server/auth.js');
var R = Rfr('lib/rethink.js');

var Routes = {
  post: {
    login: function (request, reply) {

      Authentication.validateLogin(request, function (data) {
        if (data.success !== undefined) {
          request.auth.session.set(data.session);
          reply.redirect('/index');
        } else {
          request.session.flash('loginError', data.error, true);
          reply.redirect('/auth/login');
        }
      });

    },
    totp: function (request, reply) {
      Authentication.isTOTPEnabled(request.payload.check, function (data) {
        reply(data.toString());
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
