/*
 * PufferPanel � Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var _ = require('underscore');

/** @namespace */
var TestUserModels = {};

var _fakeData = [];

TestUserModels.select = function (criteria, next) {

    return next(undefined, _.findWhere(_fakeData, criteria));
};

TestUserModels.create = function (fields, next) {

    _fakeData.push(fields);
    return next();
};

TestUserModels.update = function (id, fields, next) {

    var user = _.findWhere(_fakeData, { id: id });
    if (user === undefined) {
        return next(new Error('No such user'));
    }
    _.extend(user, fields);
    return next();
};

TestUserModels.delete = function (id, next) {

    _fakeData = _.reject(_fakeData, function (user) {

        return user.id === id;
    });
    return next();
};

TestUserModels.reset = function () {

    _fakeData = [{
        admin: true,
        email: 'admin@example.com',
        id: 'ABCDEFGH-1234-5678-9012-IJKLMNOPQRST',
        language: 'en',
        notifications: {
            loginFailure: false,
            loginSuccess: true
        },
        password: '$2a$10$y5Eekya3T6pMZzhJXFzHTudW.IrKVXhX3KPI2LB3c4UYgYdTBIqEC', //Dinosaur1
        scope: 'admin',
        sessionIp: '',
        sessionToken: '',
        totpSecret:  'ABCDEFGHIJKLMNOP12345',
        useTotp: false,
        username: 'admin'
    }, {
        admin: false,
        email: 'example@example.com',
        id: 'ABCDEFGH-1234-5678-9012-ABCDEFGHIJKL',
        language: 'en',
        notifications: {
            loginFailure: false,
            loginSuccess: true
        },
        password: '$2a$10$Jng3U.6P9FOBdLj0Vcmnn.Ob1pSHkw0qa20ZUi2hYRzxLN4G3mFmy',
        scope: 'user',
        sessionIp: '',
        sessionToken: '',
        totpSecret: 'WEFEMEFFPVTJ7GDIMS6TLPCBU4',
        useTotp: true,
        username: 'example'
    }];
};

module.exports = TestUserModels;
