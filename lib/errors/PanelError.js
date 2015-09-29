/*
 * PufferPanel - Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Util = require('util');


/**
 * Custom error class that does not allow user to see the error.
 * @param  {string,array}  data     The error data to return, can be a string or an array (or object).
 * @param  {object}        extra    An object containing additional data if needed.
 * @return {object}                 Returns an instance of the Error object named PanelError.
 */
module.exports = function PanelError (data, extra) {

    Error.captureStackTrace(this, this.constructor);
    this.name = this.constructor.name;

    this.message = data;
    this.extra = extra;

};

Util.inherits(module.exports, Error);
