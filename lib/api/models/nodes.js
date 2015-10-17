/*
 * PufferPanel - Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Rethink = Rfr('lib/rethink.js');

/** @namespace */
var NodeModels = {};

NodeModels.select = function (criteria, next) {

    Rethink.table('nodes').filter(criteria).run().then(function (nodes) {

        return next(null, nodes);
    }).error(function (err) {

        return next(err);
    });

};

NodeModels.create = function (fields, next) {

    Rethink.table('nodes').insert(fields).run().then(function () {

        return next(null);
    }).error(function (err) {

        return next(err);
    });
};

NodeModels.update = function (id, fields, next) {

    Rethink.table('nodes').get(id).update(fields).run().then(function () {

        return next(null);
    }).error(function (err) {

        return next(err);
    });
};

NodeModels.delete = function (id, next) {

    Rethink.table('nodes').get(id).delete().run().then(function () {

        return next(null);
    }).error(function (err) {

        return next(err);
    });
};

module.exports = NodeModels;
