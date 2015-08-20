var Configuration = requireFromRoot('lib/config');

var config = Configuration.rethinkdb || {
  host: 'localhost',
  port: 28015,
  database: 'pufferpanel'
};

var Rethink = {};
var connection = require('rethinkdbdash')({
  host: config.host,
  port: config.port,
  db: config.database,
  silent: true,
  buffer: 10
});

Rethink.get = function (table, row, criteria, callback) {
  connection.table(table).filter(connection.row(row).eq(criteria)).run().then(function (result) {
    return callback(result, null);
  }).error(function (err) {
      return callback(null, err);
    }
  );
};

Rethink.set = function (table, id, data, callback) {
  connection.table(table).get(id).update(data).run().error(function (err) {
    return callback(err);
  });
};

var getRawConnection = function () {
  return connection;
};

Rethink.getRawConnection = getRawConnection;

module.exports = Rethink;
