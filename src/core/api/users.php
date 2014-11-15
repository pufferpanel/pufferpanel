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
use \ORM;

/**
* PufferPanel API User Class
*/
class Users {

	protected $usersData = array();

	/**
	* Constructor Class
	* @return void
	*/
	public function __construct() { }

	/**
	 * Collects and returns data about a single user. Also provides detailed informaation about which servers the user owns.
	 *
	 * @param string $uuid UUID of user to return data about.
	 * @return array
	 */
	public function getUser($uuid) {

		$this->user = ORM::forTable('users')->rawQuery("SELECT users.*, GROUP_CONCAT(servers.hash) AS s_hash FROM users LEFT JOIN servers ON servers.owner_id = users.id WHERE users.uuid = :uuid AND servers.active = 1 LIMIT 1", array('uuid' => $uuid))->findOne();

		if(is_null($this->user->id))
			return false;
		else {

			return array(
				"id" => (int) $this->user->id,
				"email" => $this->user->email,
				"username" => $this->user->username,
				"admin" => (int) $this->user->root_admin,
				"servers" => (!empty($this->user->s_hash)) ? explode(',', $this->user->s_hash) : array()
			);

		}

	}

	/**
	 * Collects and returns data about all users in the system.
	 *
	 * @return array
	 */
	public function getUsers() {

		$this->users = ORM::forTable('users')->findMany();

		foreach($this->users as &$this->user){

			$this->usersData = array_merge($this->usersData, array(
				$this->user->uuid => array(
					"id" => (int) $this->user->id,
					"email" => $this->user->email,
					"username" => $this->user->username,
					"admin" => (int) $this->user->root_admin
				)
			));

		}

		return $this->usersData;

	}

}