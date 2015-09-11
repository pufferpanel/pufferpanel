/*
 * PufferPanel ? Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var _ = require('underscore');

/** @namespace */
var TestServerModels = {};

var _fakeData = [];

TestServerModels.getByOwner = function (ownerId, next) {

    return next(undefined, _.findWhere(_fakeData, { owner: ownerId }));
};

TestServerModels.create = function (fields, next) {

    _fakeData.push(fields);
    return next();
};

TestServerModels.update = function (id, fields, next) {

    var server = _.findWhere(_fakeData, { id: id });
    if (server === undefined) {
        return next(new Error('No such server'));
    }
    _.extend(server, fields);
    return next();
};

TestServerModels.delete = function (id, next) {

    _fakeData = _.reject(_fakeData, function (server) {

        return server.id === id;
    });
    return next();
};

TestServerModels.reset = function () {

    _fakeData = [{}];
};

module.exports = TestServerModels;
