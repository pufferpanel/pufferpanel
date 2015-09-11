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

                User.isTOTPEnabled(enabledUser, function(err, enabled) {
                    Assert.isTrue(enabled);
                });
            });
        });

        context('when disabled', function () {

            it('should be disabled', function () {

                User.isTOTPEnabled(disabledUser, function(err, enabled) {
                    Assert.isFalse(enabled);
                });
            });
        });
    });
});