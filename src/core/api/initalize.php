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
namespace PufferPanel\Core\API;

/**
* PufferPanel API Initalization Class
*/
class Initalize {

	const PP_NAMESPACE = 'PufferPanel\Core\API\\';

	/**
	 * Constructor class
	 *
	 * @return void
	 */
	public function __construct() { }

	/**
	 * Dynamically loads a given class into the API as needed.
	 *
	 * @param string $class The name of the API class to initalize.
	 * @return object
	 */
	public function loadClass($class) {

		require __DIR__.'/'.strtolower($class).'.php';
		$controller = self::PP_NAMESPACE.$class;

		return new $controller;

	}

	/**
	 * Dynamically loads a given class into the API as needed from an external source.
	 *
	 * @param string $class The name of the API class to initalize.
	 * @param string $namespace If a namespace is used for the class it must be defined. Set to "null" for a global namespace.
	 * @param string $location The location to the file containing the class.
	 * @return object
	 */
	public function loadExternalClass($class, $namespace, $location) {

		require $location;
		$controller = $namespace.$class;

		return new $controller;

	}

}