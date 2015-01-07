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
		$this->table = $table;
		$this->columnKey = $columnKey;
		$this->columnValue = $columnValue;

		$val = ORM::forTable($this->table)->select($columnKey, $columnValue)->findMany();
		foreach($val as $configOption) {
			$this->cache[$configOption->{$columnKey}] = $configOption->{$columnValue};
		}
	}

	public function config($base = null) {
		
		if(isset($this->cache[$base])) {
			return $this->cache[$base];
		} else {
			return null;
		}
		
	}

}