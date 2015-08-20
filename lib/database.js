/**
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Database = {};

var connectionModule = 'server/database/rethink';

Database.setConnection = function (path) {
  connectionModule = path;
};

Database.getConnection = function () {
  return requireFromRoot(connectionModule);
};

//The follow functions must be implemented by any Database connection
Database.get = function (table, column, criteria, callback) {
  Database.getConnection().get(table, column, criteria, callback);
};

Database.set = function (table, id, data, callback) {
  Database.getConnection().set(table, id, data, callback);
};

Database.getRawConnection = function () {
  return Database.getRawConnection();
};

module.exports = Database;
