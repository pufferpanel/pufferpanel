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

	use \PufferPanel\Core\Components\Authentication;

	/**
	 * @param array $__allowedUpdateColumns An array containing all of the fields that can be updated through a PUT request.
	 */
	private $__allowedUpdateColumns = array('whmcs_id', 'username', 'email', 'password', 'language', 'root_admin', 'notify_login_s', 'notify_login_f');

	/**
	 * @param array $__requiredUserFields An array containing all of the fields that must be sent in a request for creating a user.
	 */
	private $__requiredUserFields = array('username', 'email', 'password', 'admin', 'options');

	/**
	 * @param array $_usersData
	 */
	protected $_usersData = array();

	/**
	 * Constructor Class
	 * @return void
	 */
	public function __construct() { }

	/**
	 * Collects and returns data about a single user. Also provides detailed informaation about which servers the user owns.
	 *
	 * @param string $uuid
	 *		The UUID of user to return data about.
	 * @return array
	 */
	public function getUser($uuid) {

		$this->user = ORM::forTable('users')->rawQuery("SELECT users.*, GROUP_CONCAT(servers.hash) AS s_hash FROM users LEFT JOIN servers ON servers.owner_id = users.id AND servers.active = 1 WHERE users.uuid = :uuid LIMIT 1", array('uuid' => $uuid))->findOne();

		if(!$this->user || is_null($this->user->email)) {
			return false;
		} else {

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

		foreach($this->users as &$this->user) {

			$this->_usersData = array_merge($this->_usersData, array(
				$this->user->uuid => array(
					"id" => (int) $this->user->id,
					"email" => $this->user->email,
					"username" => $this->user->username,
					"admin" => (int) $this->user->root_admin
				)
			));

		}

		return $this->_usersData;

	}

	/**
	 * Create a new user for PufferPanel.
	 *
	 * @param array $data
	 *		An array containing all of the posted JSON data.
	 * @return string|bool
	 * 		- Returns the UUID of the user that was just created.
	 *		- Returns an integer if the operation failed which is then matched to an error in the API.
	 */
	public function createUser(array $data) {

		$this->data = $data;

		if(count(array_intersect_key(array_flip($this->__requiredUserFields), $this->data)) !== count($this->__requiredUserFields)) {
			return 1;
		}

		if(!preg_match('/^[\w-]{4,35}$/', $this->data['username'])) {
			return 2;
		}

		if(!filter_var($this->data['email'], FILTER_VALIDATE_EMAIL)) {
			return 3;
		}

		if(!filter_var($this->data['admin'], FILTER_VALIDATE_INT, array('options' => array('min_range' => 0, 'max_range' => 1)))) {
			$this->data['admin'] = 0;
		}

		$this->findExisting = ORM::forTable('users')->where_any_is(array(array('username' => $this->data['username']), array('email' => $this->data['email'])))->findOne();
		if($this->findExisting) {
			return 4;
		}

		if(array_key_exists('password', $this->data['options']) && $this->data['options']['password'] === true) {

			$this->data['password_raw'] = self::keygen(12);
			$this->data['password'] = $this->hash($this->data['password_raw']);

		} else {

			$this->data['password'] = $this->hash($this->data['password']);

		}

		$this->data['uuid'] = $this->generateUniqueUUID('users', 'uuid');
		$this->user = ORM::forTable('users')->create();
		$this->user->set(array(
			'uuid' => $this->data['uuid'],
			'username' => $this->data['username'],
			'email' => $this->data['email'],
			'password' => $this->data['password'],
			'language' => 'en',
			'register_time' => time()
		));
		$this->user->save();

		return array('uuid' => $this->data['uuid']);

	}

	/**
	 * Updates information about a user given an array.
	 *
	 * @param string $uuid
	 * 		The UUID of the user whom you are trying to update.
	 * @param array $data
	 * 		An array containing the coulmn names to update as keys and their value being the data you wish to replace it with.
	 * @return boolean
	 * 		- Returns true if the operation was successful and all of the data was updated.
	 * 		- Returns an integer if the operation failed which is then matched to an error in the API.
	 */
	public function updateUser($uuid, array $data) {

		$this->uuid = $uuid;
		$this->data = $data;

		foreach($this->data as $key => $value) {

			if(!in_array($key, $this->__allowedUpdateColumns)) {
				return false;
			}

			switch($key) {

				case 'whmcs_id':
					if(!is_numeric($this->data['whmcs_id'])) {
						return false;
					}
					break;
				case 'username':
					if(!preg_match('/^[\w-]{4,35}$/', $this->data['username'])) {
						return false;
					}
					break;
				case 'email':
					if(!filter_var($this->data['email'], FILTER_VALIDATE_EMAIL)) {
						return false;
					}
					break;
				case 'password':
					if(strlen($this->data['password']) < 8) {
						return false;
					} else {
						$this->data['password'] = $this->hash($this->data['password']);
					}
					break;
				case 'root_admin':
					if(!filter_var($this->data['root_admin'], FILTER_VALIDATE_INT, array('options' => array('min_range' => 0, 'max_range' => 1)))) {
						return false;
					}
					break;
				case 'notify_login_s':
					if($this->data['notify_login_s'] > 1 || $this->data['notify_login_s'] < 0) {
						return false;
					}
					break;
				case 'notify_login_f':
					if($this->data['notify_login_f'] > 1 || $this->data['notify_login_f'] < 0) {
						return false;
					}
					break;
				default:
					return false;

			}

		}

		$this->user = ORM::forTable('users')->where('uuid', $this->uuid)->findOne();
		$this->user->set($this->data);
		$this->user->save();

		return true;

	}

}