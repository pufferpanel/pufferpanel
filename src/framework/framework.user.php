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

/**
 * PufferPanel Core User Class File
 */
class user extends Auth\auth {

	/**
	 * @param array $_data Implements a blank array for the functions to write to.
	 */
	private $_data;

	/**
	 * @param bool $_l Defaults to true and will be changed to false if there is an error.
	 */
	private $_l;

	/**
	 * @param string $_shash Private variable used for keeping track of server we are interested in for permissions.
	 */
	private static $_shash;

	/**
	* @param int $_soid Private variable used for keeping track of server owner id for permissions.
	*/
	private static $_soid;

	/**
	* @param string $_perms Private variable used for keeping track of what permissions a user hash.
	*/
	private $_perms = null;

	/**
	 * Constructor Class responsible for filling in arrays with the data from a specified user.
	 *
	 * @param string $ip The IP address of a user who is requesting the function, or if called from the Admin CP it is the user id.
	 * @param mixed $session The value of the pp_auth_token cookie.
	 * @param mixed $hash The server hash of the requesting user which is used when they are viewing node pages.
	 * @return void
	 */
	public function __construct($ip, $session = null, $hash = null){

		$this->mysql = self::connect();

		/*
		 * Reset Values
		 */
		$this->_data = array();
		$this->_l = true;

		if(self::isLoggedIn($ip, $session, $hash) === true){

			$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sesip AND `session_id` = :sesid");
			$this->query->execute(array(':sesip' => $ip, ':sesid' => $session));

				$this->row = $this->query->fetch();
				if(is_array($this->row))
					foreach($this->row as $this->id => $this->val)
						$this->_data[$this->id] = $this->val;

		}else if(is_null($session) && is_null($hash) && is_numeric($ip)){

			$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `id` = :id");
			$this->query->execute(array(':id' => $ip));

				$this->row = $this->query->fetch();
				if(is_array($this->row))
					foreach($this->row as $this->id => $this->val)
						$this->_data[$this->id] = $this->val;

		}else
			$this->_l = false;

	}

	/**
	 * Re-runs the __construct() class with a defined ID for the admin control panel.
	 *
	 * @param string $id The ID of a user requested in the Admin CP.
	 * @return void
	 */
	public function rebuildData($id){

		$this->__construct($id);

	}

	/**
	 * Provides the corresponding value for the id provided from the MySQL Database.
	 *
	 * @param mixed $id The column value for the data you need (e.g. email).
	 * @return mixed A string is returned on success, array if nothing was passed, and if the command fails 'false' is returned.
	 */
	public function getData($id = null){

		if(is_null($id))
			if($this->_l === true)
				return $this->_data;
			else
				return false;
		else
			if($this->_l === true && is_array($this->_data) && array_key_exists($id, $this->_data))
				return $this->_data[$id];
			else
				return false;

	}

	/**
	 * Initiator class for server based on Hash.
	 *
	 * @param string $server The hash of the server we are interested in.
	 * @param int $oid Returns the owner of the server in question
	 * @return void
	 * @static
	 */
	static public function permissionsInit($server, $oid) {

		self::$_shash = $server;
		self::$_soid = $oid;

	}

	/**
	 * Collects permissions from the Database given a user id.
	 *
	 * @param bool $list Set to true to return an array of all servers that a user has permission for. Defaults to false.
	 * @return void
	 */
	public function getPermissions($list = false) {

		$this->query = $this->mysql->prepare("SELECT `permissions` FROM `users` WHERE `id` = :uid");
		$this->query->execute(array(
			':uid' => self::$_soid
		));

			if($this->query->rowCount() == 0)
				return false;
			else {

				$this->row = $this->query->fetch();
				if(is_null($this->row['permissions']) || empty($this->row['permissions']))
					return false;
				else {

					$this->json = json_decode($this->row['permissions'], true);
					return ($list === false) ? ((!is_int($this->_shash)) ? false : $this->json[$this->_shash]) : $this->json;

				}

			}

	}

	/**
	 * Lists all servers that a user has permission to view
	 *
	 * @return array
	 */
	public function listServerPermissions() {

		if(is_null($this->_perms))
			$this->_perms = $this->getPermissions(true);

		if(is_array($this->_perms))
			return array_keys($this->_perms);
		else
			return array("0" => "0");

	}

	/**
	 * Checks if a given user has permission to access a part of the Control Panel. Defaults to true if the user is the owner.
	 *
	 * @param string $permission The permission node to check aganist.
	 * @return bool Returns true if they have permission, false if not.
	 */
	public function checkPermission($permission) {

		if($this->_soid != $this->getData('id')){

			if(is_null($this->_perms))
				$this->_perms = $this->getPermissions();

			if(is_array($this->_perms)) {

				if(array_key_exists($permission, $this->_perms['permissions']))
					return true;
				else
					return false;

			}else
				return false;

		}else
			return true;

	}

	/**
	 * Check if a user has a protected permissions class.
	 *
	 * @return bool
	 */
	public function checkProtected() {

		if(is_null($this->_perms))
			$this->_perms = $this->getPermissions();

		if(is_array($this->_perms)) {

			if($this->_perms['protected'] == 1)
				return true;
			else
				return false;

		}else
			return true;

	}

}

?>
