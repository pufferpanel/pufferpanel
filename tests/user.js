/**
 * PufferPanel ? Reinventing the way game servers are managed.
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
var User = Rfr('lib/controllers/user.js');
var UserModels = Rfr('tests/models/users.js');
Chai.config.includeStack = false;
var Assert = Chai.assert;

describe('Controller/User', function () {

    beforeEach(function () {

        UserModels.reset();
    });

    describe('getData', function () {

        context('when run', function () {

            var userId = 'ABCDEFGH-1234-5678-9012-IJKLMNOPQRST';

            it('should return user data object', function () {

                User.getData(userId, function (err, user) {
                    Assert.isNull(err);
                    Assert.isObject(user);
                    Assert.property(user, 'id');
                    Assert.property(user, 'email');
                    Assert.property(user, 'password');
                });

            });

        });

    });

    describe('generateTOTP', function () {

        context('when generated', function () {

            var userId = 'ABCDEFGH-1234-5678-9012-IJKLMNOPQRST';
            var totpRegex = /^[A-Z0-9]{26}/g;

            it('should be valid', function () {

                User.generateTOTP(userId, function (err, totp) {
                    var totpSecret = totp.secret;
                    Assert.isTrue(totpRegex.test(totpSecret));
                });
            });
        });
    });

    describe('isTOTPEnabled', function () {

        var enabledUser = 'example@example.com';
        var disabledUser = 'admin@example.com';

        context('when enabled', function () {

            it('should be enabled', function () {

                User.isTOTPEnabled(enabledUser, function (err, enabled) {
                    Assert.isTrue(enabled);
                });
            });
        });

        context('when disabled', function () {

            it('should be disabled', function () {

                User.isTOTPEnabled(disabledUser, function (err, enabled) {
                    Assert.isFalse(enabled);
                });
            });
        });
    });

    describe('updatePassword', function () {

        var userId = 'ABCDEFGH-1234-5678-9012-IJKLMNOPQRST';
        var user = 'admin@example.com';
        var userCurrentPassword = 'Dinosaur1';
        var newUserPassword = 'Pterodactyl1';

        context('when called', function () {

            it('should update password', function () {

                User.updatePassword(userId, userCurrentPassword, newUserPassword, function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isUndefined(response);
                });

            });

        });

        context('when sent an invalid password', function () {

            it('should fail', function () {

                User.updatePassword(userId, 'invalid', newUserPassword, function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isString(response);
                });

            });

        });

    });

    describe('updateNotifications', function () {

        var userId = 'ABCDEFGH-1234-5678-9012-IJKLMNOPQRST';
        var userPassword = 'Dinosaur1';

        context('when sent valid data', function () {

            it('should be successful', function () {

                User.updateNotifications(userId, userPassword, true, true, function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isUndefined(response);
                });

            });

        });

        context('when sent invalid data', function () {

            it('should not be successful', function () {

                User.updateNotifications(userId, userPassword, 'invalid', true, function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isArray(response);
                });

            });

        });

        context('when sent an invalid password', function () {

            it('should not be successful', function () {

                User.updateNotifications(userId, 'invalidpassword', true, true, function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isString(response);
                });

            });

        });

    });

    describe('updateEmail', function () {

        var userId = 'ABCDEFGH-1234-5678-9012-IJKLMNOPQRST';
        var userPassword = 'Dinosaur1';

        context('when sent valid data', function () {

            it('should be successful', function () {

                User.updateEmail(userId, userPassword, 'newemail@example.com', function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isUndefined(response);
                });

            });

        });

        context('when sent invalid email', function () {

            it('should not be successful', function () {

                User.updateEmail(userId, userPassword, 'invalid', function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isArray(response);
                });

            });

        });

        context('when sent an invalid password', function () {

            it('should not be successful', function () {

                User.updateEmail(userId, 'invalidpassword', 'newemail@example.com', function (err, response) {
                    Assert.isTrue(!err);
                    Assert.isString(response);
                });

            });

        });

    });

});
