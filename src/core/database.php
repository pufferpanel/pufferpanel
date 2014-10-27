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

/**
 * Core PufferPanel Database Class.
 * Designed to keep track of the number of queries and allow for advanced debugging if necessary.
 */
class Database_Initiator extends \PDO {

	use Components\Database;

	/**
	* @param int $queryCounter
	* @static
	*/
	private static $queryCounter = 0;

	/**
	* @param string $query
	* @return object
	*/
	public function query($query) {

	    ++self::$queryCounter;
	    return parent::query($query);

	}

	/**
	 * @param string $statement
	 * @param array $options
	 * @return object
	 */
	public function prepare($statement, $options = array()){

	    ++self::$queryCounter;
	    return parent::prepare($statement, $options);

	}

	/**
	 * @param string $statement
	 * @return object
	 */
	public function exec($statement) {

	    ++self::$queryCounter;
	    return parent::exec($statement);

	}

	/**
	 * @return int Returns the total number of queries executed on a page.
	 * @static
	 */
	public static function getCount(){

		return self::$queryCounter;

	}

}
?>
