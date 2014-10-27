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
 * PufferPanel Core User Class File
 */
class User extends Authentication {

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
	* @param int $_uid Private variable used for keeping track of current user id for permissions.
	*/
	private static $_uid;

	/**
	* @param int $_oid Private variable used for keeping track of server owner id for permissions.
	*/
	private static $_oid;

	/**
	* @param string $_perms Private variable used for keeping track of what permissions a user hash.
	*/
	private $_perms = null;

	private static $_row;

	private static $_json;

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
	 * @param int $uid Returns the current user's ID
	 * @param int $oid Returns the server owner's ID
	 * @return void
	 * @static
	 */
	public static function permissionsInit($server, $uid, $oid = null) {

		self::$_shash = $server;
		self::$_uid = $uid;
		self::$_oid = $oid;

	}

	/**
	 * Collects permissions from the Database given a user id.
	 *
	 * @param bool $list Set to true to return an array of all servers that a user has permission for. Defaults to false.
	 * @return void
	 */
	public function getPermissions($list = false) {

		if(!isset($this->_permissionJson)){

			$this->permissionQuery = $this->mysql->prepare("SELECT `permissions` FROM `users` WHERE `id` = :uid");
			$this->permissionQuery->execute(array(
				':uid' => self::$_uid
			));

			if($this->permissionQuery->rowCount() == 0)
				$this->permissionRow = null;
			else
				$this->permissionRow = $this->permissionQuery->fetch();

		}

		if(!is_null($this->permissionRow))
			if(array_key_exists('permissions', $this->permissionRow) && !empty($this->permissionRow['permissions']))
				$this->_permissionJson = json_decode($this->permissionRow['permissions'], true);
			else
				$this->_permissionJson = null;
		else
			$this->_permissionJson = null;

		return (is_null($this->_permissionJson) || !isset($this->_permissionJson)) ? false : (($list === false) ? ((!array_key_exists(self::$_shash, $this->_permissionJson)) ? false : $this->_permissionJson[self::$_shash]) : $this->_permissionJson);

	}

	/**
	 * Lists all servers that a user has permission to view
	 *
	 * @return array
	 */
	public function listServerPermissions() {

		$this->_perms = $this->getPermissions(true);
		if(is_array($this->_perms))
			return array_keys($this->_perms);
		else
			return array("0" => "0");

	}

	/**
	 * Returns a list of all permissions nodes avaliable
	 *
	 * @return array
	 * @static
	 */
	private static function listAvaliablePermissions() {

		return array(
			'console' => array('view', 'commands', 'power'),
			'files' => array('view', 'edit', 'save', 'download', 'delete'),
			'manage' => array('view', 'rename' => array('view', 'jar'), 'ftp' => array('view', 'details', 'password')),
			'users' => array('view')
		);

	}

	/**
	 * Returns permissions for a user in a twig friendly format
	 *
	 * @param array $array
	 * @return array
	 */
	public function twigListPermissions($array = null) {

		$this->buildPermissions = array();
		$this->allPerms = self::listAvaliablePermissions();

		foreach($this->allPerms as $permission => $submodule)
			foreach($submodule as $id => $subpermission)
				if(!is_array($subpermission))
					$this->buildPermissions[$permission][$subpermission] = $this->hasPermission($permission.".".$subpermission, $array);
				else
					foreach($subpermission as $subid => $subsubpermission)
						$this->buildPermissions[$permission][$id][$subsubpermission] = $this->hasPermission($permission.".".$id.".".$subsubpermission, $array);

		return $this->buildPermissions;

	}

	/**
	 * Checks if a given user has permission to access a part of the Control Panel. Defaults to true if the user is the owner.
	 *
	 * @param string $permission The permission node to check aganist.
	 * @param array $array
	 * @return bool Returns true if they have permission, false if not.
	 */
	public function hasPermission($permission, $array = null) {

		if(!is_null($array))
			$this->_perms = $array;
		else
			if(is_null($this->_perms))
				$this->_perms = $this->getPermissions();

		if(self::$_oid != $this->getData('id') && is_null($array))
			if(is_array($this->_perms))
				return in_array($permission, $this->_perms);
			else
				return false;
		else
			if(!is_array($array))
				return true;
			else
				if(is_array($this->_perms))
					return in_array($permission, $this->_perms);
				else
					return false;

	}

}

?>
