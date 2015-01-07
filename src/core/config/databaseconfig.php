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

class DatabaseConfig implements ConfigInterface {

	private $table;
	private $columnKey;
	private $columnValue;

	public function __construct($table = 'settings', $columnKey = 'key', $columnValue = 'value') {
		$this->table = $table;
		$this->columnKey = $columnKey;
		$this->columnValue = $columnValue;
	}

	public function config($base = null) {
		if($base == null) {
			throw new \Exception('Cannot get null key from database');
		}
		$val = ORM::forTable($this->table)->where($this->columnKey, $base)->select($this->columnValue)->findOne()->asArray();
		return $val[$this->columnValue];
	}

}