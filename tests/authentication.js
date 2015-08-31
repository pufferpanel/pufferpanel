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
var Authentication = Rfr('lib/controllers/authentication.js');

describe('Controller/Authentication', function () {

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
