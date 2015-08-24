/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Routes = {
  Authentication: Rfr('lib/routes/authentication.js'),
  Base: Rfr('lib/routes/base.js'),
  Server: Rfr('lib/routes/server.js')
};

module.exports = Routes;
