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
namespace PufferPanel\Core\Components;
use \Exception;

class Config {

	/**
	 * Global config for the panel
	 * @param object $config
	 * @static
	 */
	private static $globalConfig;

	/**
	 * Local instance of this config
	 */
	private $config;

	/**
	 * Returns config variables from the config.json file.
	 *
	 * @param string $base
	 * @param bool $array
	 */
	final public static function getGlobal($base = null, $array = false) {
		if(is_null(self::$globalConfig)) {
			self::$globalConfig = new Config('config.json');
		}

		return (is_null($base)) ? self::$globalConfig : self::$globalConfig->{$base};

	}

	public function __construct($path) {
		if(!file_exists(BASE_DIR.$path)) {
			throw new Exception("The config file ".$path." does not exist.");
		}

		$this->config = json_decode(file_get_contents(BASE_DIR.'config.json'));

		if(json_last_error() != "JSON_ERROR_NONE") {
			throw new Exception("An error occured when trying decode the config.json file. ".json_last_error());
		}
	}

}