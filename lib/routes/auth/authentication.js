/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Authentication = requireFromRoot('server/auth/authentication');

var Routes = {
  post: {
    login: function (request, reply) {

      Authentication.loginUser(request.payload.email, request.payload.password, request.payload.totp_token, request.payload.remoteAddress, function (err, success, data) {
        if (err === undefined && success) {
          request.auth.session.set(data);
          reply.redirect('/index');
        } else {
          request.session.flash('loginError', data, true);
          reply.redirect('/auth/login');
        }
      });

    },
    totp: function (request, reply) {

      Authentication.isTOTPEnabled(request.payload.check, function (data) {
        if (data === true) {
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
