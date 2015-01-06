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
namespace PufferPanel\Core\Router;
use \ORM;

class Router_Controller {

	const PP_NAMESPACE = 'PufferPanel\Core\Router\\';

	protected $_class;

	protected $_pass;

	/**
	 * Constructor class for Routes that handles loading class files on the fly.
	 *
	 * @return void
	 */
	public function __construct($class, $pass = false) {

		$this->_class = $class;
		$this->_pass = $pass;

	}

	/**
	 * Handles loading the controller for a given class file.
	 *
	 * @return object
	 */
	public function loadClass() {

		require __DIR__.'/routes/'.strtolower(str_replace('\\', '/', $this->_class)).'.php';
		$controller = self::PP_NAMESPACE.$this->_class;

		return (!$this->_pass) ? new $controller : new $controller($this->_pass);

	}

}