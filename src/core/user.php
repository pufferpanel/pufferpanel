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
use \ORM, \ReflectionClass;

/**
 * PufferPanel Core User Class File
 */
class User extends Authentication {

	/**
	 * @param object $user
	 */
	protected $user;

	/**
	 * @param string $permissions_server
	 */
	private static $permissions_server;

	/**
	 * @param int $permissions_owner
	 */
	private static $permissions_owner;

	/**
	 * Constructor Class responsible for filling in arrays with the data from a specified user.
	 *
	 * @param mixed $rebuild If passed it should be the ID of the user to rebuild the data as. If passed as false it will build data as the logged in user.
	 * @return void
	 */
	public function __construct($rebuild = false){

		parent::__construct();

		if(parent::isLoggedIn() && !$rebuild) {
			$this->user = $this->select;
		} else if($rebuild) {
			$this->user = ORM::forTable('users')->findOne($rebuild);
		} else {
			$this->user = false;
		}

	}

	/**
	 * Re-runs the __construct() class with a defined ID for the admin control panel.
	 *
	 * @param string $id The ID of a user requested in the Admin CP.
	 * @return bool
	 */
	public function rebuildData($id){

		self::__construct($id);

		return (!$this->user) ? false : true;

	}

	/**
	 * Provides the corresponding value for the id provided from the MySQL Database.
	 *
	 * @param mixed $id The column value for the data you need (e.g. email).
	 * @return mixed A string is returned on success, array if nothing was passed, and if the command fails 'false' is returned.
	 */
	public function getData($id = null){

		if($this->user && !is_null($this->user)) {

			if(is_null($id)) {

				$reflect = new ReflectionClass($this->user);
				$data = $reflect->getProperty('_data');
				$data->setAccessible(true);

				return $data->getValue($this->user);

			} else {
				return (isset($this->user->{$id})) ? $this->user->{$id} : false;
			}

		} else {
			return false;
		}

	}

}