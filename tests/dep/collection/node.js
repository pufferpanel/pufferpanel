/*
 Copyright 2016 Joshua Taylor <lordralex@gmail.com>

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
 */

var _ = require('underscore');

var NodeCollection = function () {
    this.db = [];
};

NodeCollection.prototype.get = function (uuid) {
    var self = this;
    var result = _.find(self.db, function (node) {
        return node.uuid == uuid;
    });
    return result;
};

NodeCollection.prototype.getByLocation = function (location) {
    var self = this;
    var results = _.filter(self.db, function (node) {
        return node.location == location;
    });
    return results;
};

NodeCollection.prototype.add = function (node) {
    var self = this;
    self.db.push(node);
};

NodeCollection.prototype.remove = function (uuid) {
    var self = this;
    self.db = _.reject(self.db, function (node) {
        return node.uuid == uuid;
    });
};

NodeCollection.prototype.update = function (uuid, newValues) {
    var self = this;
    _.each(self.db, function (node) {
        if (node.uuid == uuid) {
            _.extend(node, newValues);
        }
    });
    return self.get(uuid);
};

NodeCollection.prototype._reset = function (data) {
    var self = this;
    self.db = _.clone(data);
};

module.exports = NodeCollection;