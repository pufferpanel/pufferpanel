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
var Authentication = Rfr('lib/controllers/authentication.js');
var UserModels = Rfr('tests/models/users');
Chai.config.includeStack = false;
var Assert = Chai.assert;

describe('Controller/Authentication', function () {

    beforeEach(function () {

        UserModels.reset();
    });

    describe('generatePasswordHash', function () {

        var rawPw = 'admin';
        var hashRegex = /^\$2a\$10\$.{53}/g;

        context('when generates', function () {

            it('should be hashed', function () {

                Assert.isTrue(hashRegex.test(Authentication.generatePasswordHash(rawPw)));
            });
        });
    });

    describe('validateAccountPassword', function () {

        var userId = 'ABCDEFGH-1234-5678-9012-IJKLMNOPQRST';
        var goodPassword = 'Dinosaur1';
        var badPassword = 'NotDinosaur1';

        context('when password is correct', function () {

            it('should return null', function () {

                Authentication.validateAccountPassword(userId, goodPassword, function (err) {
                    Assert.isNull(err);
                });

            });

        });

        context('when password is incorrect', function () {

            it('should return false', function () {

                Authentication.validateAccountPassword(userId, badPassword, function (err) {
                    Assert.isFalse(err);
                });

            });

        });

    });

    describe('loginUser', function () {

        var email = 'admin@example.com';
        var badEmail = 'donotuse@example.com';
        var password = 'Dinosaur1';
        var badPassword = 'wrong';
        var ipAddress = '127.0.0.1';

        context('when email and password are correct', function () {

            it('should correctly log in user', function () {

                Authentication.loginUser(email, password, null, ipAddress, function (err, data) {
                    Assert.isTrue(!err);
                    Assert.isNotString(data);
                    Assert.isObject(data);
                    Assert.property(data, 'id');
                    Assert.property(data, 'sessionToken');
                    Assert.property(data, 'sessionIp');
                });
            });
        });

        context('when email is correct and password is incorrect', function () {

            it('should fail to log in user', function () {

                Authentication.loginUser(email, badPassword, null, ipAddress, function (err, data) {

                    Assert.isTrue(!err);
                    Assert.isString(data);
                });
            });
        });

        context('when email does not exist', function () {

            it('should fail to log in user', function () {

                Authentication.loginUser(badEmail, password, null, ipAddress, function (err, data) {

                    Assert.isTrue(!err);
                    Assert.isString(data);
                });
            });
        });
    });
});
