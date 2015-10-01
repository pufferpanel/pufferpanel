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
var TestActionsModel = {};

var _fakeData = [];

TestActionsModel.reset = function () {

    _fakeData = [{}];
};

TestActionsModel.select = function (data, next) {

    return next(null, _.findWhere(_fakeData, criteria));

};

TestActionsModel.create = function (data, next) {

    _fakeData.push(fields);
    return next();

};

TestActionsModel.deleteId = function (id, next) {

    _fakeData = _.reject(_fakeData, function (server) {

        return server.id === id;
    });
    return next();

};

module.exports = TestActionsModel;
