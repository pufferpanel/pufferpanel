/*
 * PufferPanel - Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var _ = require('underscore');
var Rfr = require('rfr');
var Rethink = Rfr('lib/rethink.js');

/** @namespace */
var UserModels = {};

UserModels.select = function (criteria, next) {

    Rethink.table('users').filter(criteria).run().then(function (user) {

        return next(null, _.first(user));
    }).error(function (err) {

        return next(err);
    });
};

UserModels.create = function (fields, next) {
    Rethink.table('users').insert(fields).run().then(function () {

        return next(null);
    }).error(function (err) {

        return next(err);
    });
};

UserModels.update = function (id, fields, next) {

    Rethink.table('users').get(id).update(fields).run().then(function () {

        return next(null);
    }).error(function (err) {

        return next(err);
    });
};

UserModels.delete = function (id, next) {

    Rethink.table('users').get(id).delete().run().then(function () {

        return next(null);
    }).error(function (err) {

        return next(err);
    });
};

module.exports = UserModels;
