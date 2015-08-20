/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Async = require('async');
var Logger = requireFromRoot('lib/logger');
var Authentication = requireFromRoot('lib/service/authentication');

var postLogin = function (request, reply) {
  Authentication.loginUser(request.payload.email, request.payload.password, request.payload.totp_token, request.payload.remoteAddress, function (result, data) {
    if (result) {
      request.auth.session.set(data.user);
      reply.redirect('/index');
    } else {
      request.session.flash('loginError', data, true);
      reply.redirect('/auth/login');
    }
  });
};

var postTotp = function (request, reply) {
  Async.series([
    function (callback) {
      Authentication.isTOTPEnabled(request.payload.check, callback);
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
};

var postRegister = function (requset, reply) {
};

var postPassword = function (request, reply) {
};

var getLogin = function (request, reply) {
  reply.view('auth/login', {
    flash: request.session.flash('loginError')
  });
};

var getLogout = function (request, reply) {
  request.auth.session.clear();
  reply.redirect('/auth/login');
};

var getRegister = function (request, reply) {
  reply.view('auth/register', {
    flash: request.session.flash('registerError'),
    token: request.params.token
  });
};

var getPassword = function (request, reply) {
  reply.view('auth/password', {
    flash: request.session.flash('passwordError'),
    noshow: false
  });
};

var Routes = {
  post: {
    login: postLogin,
    totp: postTotp,
    register: postRegister,
    password: postPassword
  },
  get: {
    login: getLogin,
    logout: getLogout,
    register: getRegister,
    password: getPassword
  }
};

module.exports = Routes;
