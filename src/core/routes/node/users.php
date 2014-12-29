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
namespace PufferPanel\Core\Router\Node;
use \ORM;

class Users extends \PufferPanel\Core\Email {

	protected $_server;

	protected static $_error = null;

	/**
	 * Constructor class for \PufferPanel\Core\Router\Node_Users
	 *
	 * @return void
	 */
	public function __construct($server) {

		$this->_server = $server;

	}

	/**
	 * Handles adding a new subuser to a server.
	 *
	 * @return bool
	 */
	public function addSubuser($data = array()) {

		self::_setError('This function is not currently functional.');

		return false;

	}

	/**
	 * Controller for setting the last error that occured. Called when addSubuser() returns false.
	 *
	 * @return void
	 */
	protected final static function _setError($error) {

		self::$_error = $error;

	}

	/**
	 * Controller for retrieving the last error that occured. Used when addSubuser() returns false.
	 *
	 * @param bool $wrapped If set to false will return the error not wrapped in <div> error tags.
	 * @return string
	 */
	public function retrieveLastError($wrapped = true) {

		return (!$wrapped) ? self::$_error : '<div class="alert alert-danger">'.self::$_error.'</div>';

	}

}