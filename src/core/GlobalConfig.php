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
namespace PufferPanel\Core;
use \Exception;

class GlobalConfig implements ConfigInterface {

	/**
	 * @param object $config
	 */
	private static $config;

	/**
	 * {@inheritdoc}
	 */
	public function __construct($path, $array = false) {

		if(!file_exists(BASE_DIR.$path)) {
			throw new Exception("The config file ".$path." does not exist.");
		}

		self::$config = json_decode(file_get_contents(BASE_DIR.'config.json'), $array);

		if(json_last_error() != "JSON_ERROR_NONE") {
			throw new Exception("An error occured when trying decode the config.json file. ".json_last_error());
		}

	}

	/**
	 * {@inheritdoc}
	 */
	public static function config($base = null) {

		if(is_null(self::$config)) {
			new self('config.json');
		}

		return (is_null($base)) ? self::$config : self::$config->{$base};

	}

}