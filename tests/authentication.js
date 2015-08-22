/**
 * PufferPanel � Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
process.env.NODE_ENV = 'test';

var Rfr = require('rfr');
var Chai = require('chai');
Chai.config.includeStack = false;
var Should = Chai.should();
var Expect = Chai.expect;
var Assert = Chai.assert;
var Authentication = Rfr('server/auth.js');

describe('Server/Auth', function () {

  describe('validateLogin', function () {

    context('when user logs in with valid credentials', function () {
      it('should return successfully', function (done) {
        Authentication.validateLogin({
          payload: {
            email: 'theoretical@email.com',
            password: 'theoretically_correct_password'
          }
        }, function (reply) {
          Expect(reply.error).to.be.undefined;
          Expect(reply.success).to.be.defined;
          Expect(reply.success).to.be.true;
          done();
        });
      });
    });

    context('when user logs in with invalid credentials', function () {
      it('should return unsuccessfully', function (done) {
        Authentication.validateLogin({
          payload: {
            email: 'bad@email.com',
            password: 'FalsePassword'
          }
        }, function (reply) {
          Assert.isUndefined(reply.success, 'success is undefined');
          Assert.isDefined(reply.error, 'error is defined');
          done();
        });
      });
    });

  });

  describe('updatePasswordHash', function () {

    var oldPw = '$2y$10$CrEJiLF5OoK/D.FgBs8Wc.Kr0C0KZaxWwOJwlYI4P98wjHP9BzXnK';
    var newPw = '$2a$10$CrEJiLF5OoK/D.FgBs8Wc.Kr0C0KZaxWwOJwlYI4P98wjHP9BzXnK';
    var validatedPw = '$2a$10$CrEJiLF5OoK/D.FgBs8Wc.Kr0C0KZaxWwOJwlYI4P98wjHP9BzXnK';

    context('when needs update', function () {
      it('should update', function () {
        Assert.strictEqual(Authentication.updatePasswordHash(oldPw), validatedPw, 'passwords are the same.');
      });
    });

    context('when already converted', function () {
      it('should not update', function () {
        Assert.strictEqual(Authentication.updatePasswordHash(newPw), validatedPw, 'passwords are the same.');
      });
    });

  });

  describe('generatePasswordHash', function () {

    var rawPw = 'admin';
    var hashRegex = /^\$2a\$10\$.{53}/g;

    context('when generates', function () {
      it('should be hashed', function () {
        Assert.isTrue(hashRegex.test(Authentication.generatePasswordHash(rawPw)), 'the password hash is valid.');
      });
    });
  });

});
