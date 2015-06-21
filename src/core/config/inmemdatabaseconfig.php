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
use \ORM;

class InMemoryDatabaseConfig implements ConfigInterface {

	private $cache = array();

	public function __construct($table = 'settings', $columnKey = 'key', $columnValue = 'value') {
		
		$val = ORM::forTable($table)->select($columnKey, 'key')->select($columnValue, 'value')->findMany();

		$temp = array();
		foreach($val as $configOption) {
			$temp[$configOption->key] = $configOption->value;
		}

		$this->cache = (object) $temp;
	}

	public function config($base = null) {
		
		return (is_null($base)) ? $this->cache : $this->cache->{$base};

	}

}