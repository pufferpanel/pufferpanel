var Authentication = require('../server/auth/authentication.js');
require('should');

describe('Server/Auth/Authentication', function () {
  describe('updatePasswordHash', function () {

    var oldPw = '$2y$10$CrEJiLF5OoK/D.FgBs8Wc.Kr0C0KZaxWwOJwlYI4P98wjHP9BzXnK';
    var newPw = '$2a$10$CrEJiLF5OoK/D.FgBs8Wc.Kr0C0KZaxWwOJwlYI4P98wjHP9BzXnK';
    var validatedPw = '$2a$10$CrEJiLF5OoK/D.FgBs8Wc.Kr0C0KZaxWwOJwlYI4P98wjHP9BzXnK';

    context('when needs update', function () {
      it('should update', function () {
        Authentication.updatePasswordHash(oldPw).should.equal(validatedPw);
      });
    });

    context('when already converted', function () {
      it('should not update', function () {
        Authentication.updatePasswordHash(newPw).should.equal(validatedPw);
      })
    })

  });

  describe('generatePasswordHash', function () {

    var rawPw = 'admin';
    var hashRegex =/^\$2a\$10\$.{53}/g;

    context('when generates', function () {
      it('should be hashed', function () {
        hashRegex.test(Authentication.generatePasswordHash(rawPw)).should.be.true();
      });
    })
  });

});