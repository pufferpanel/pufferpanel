<?php
/*
PufferPanel - A Minecraft Server Management Panel
Copyright (c) 2013 Dane Everitt

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see http://www.gnu.org/licenses/.
*/
namespace PufferPanel\Core\Config;

interface ConfigInterface {

	/**
	 * Constructor class for implementing configuration files from JSON.
	 *
	 * @param string $path Configuration file relative to BASE_DIR
	 * @param bool $array
	 * @return void
	 */
	public function __construct($path, $array = false);

	/**
	 * Gets a specific key from a config, or the entire config if key is null.
	 *
	 * @param type $base Key
	 */
	public function config($base = null);

}