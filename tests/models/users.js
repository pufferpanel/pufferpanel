/*
 * PufferPanel — Reinventing the way game servers are managed.
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
        return next(new Error("No such user"));
    }
    _.extend(user, fields);
    return next();
};

TestUserModels.delete = function (id, next) {

    _fakeData = _.reject(_fakeData, function (user) {

        return user.id === id;
    })
};

TestUserModels.reset = function () {

    _fakeData = [{
        id: 1,
        email: 'admin@example.com',
        password: '$2a$10$Jng3U.6P9FOBdLj0Vcmnn.Ob1pSHkw0qa20ZUi2hYRzxLN4G3mFmy'
    }];
};

module.exports = TestUserModels;
