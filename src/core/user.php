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
	 * @param array $permissions
	 */
	protected $permissions;

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
	 * @return void
	 */
	public function rebuildData($id){

		self::__construct($id);

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

	/**
	 * Initiator class for server based on Hash.
	 *
	 * @param string $server The hash of the server we are interested in.
	 * @param int $oid Sets the server owner.
	 * @return void
	 * @static
	 */
	public static function initalizePermissions($server, $oid) {

		self::$permissions_server = $server;
		self::$permissions_owner = $oid;

	}

	/**
	 * Collects permissions from the Database given a user id.
	 *
	 * @param bool $list Set to true to return an array of all servers that a user has permission for. Defaults to false.
	 * @return void
	 */
	final protected function _getPermissions($list = false) {

		if(self::getData('permissions') && !empty(self::getData('permissions'))) {

			$this->permissions = json_decode(self::getData('permissions'), true);

		}

		if(!$list && !is_null($this->permissions)) {

			if(array_key_exists(self::$permissions_server, $this->permissions)) {
				$this->permissions = $this->permissions[self::$permissions_server]['perms'];
			}

		}

	}

	/**
	 * Lists all servers that a user has permission to view
	 *
	 * @return array
	 */
	public function listServerPermissions() {

		$this->_getPermissions(true);
		return (is_array($this->permissions)) ? array_keys($this->permissions) : array("0" => "0");

	}

	/**
	 * Returns a list of all permissions nodes avaliable
	 *
	 * @return array
	 * @static
	 */
	final protected static function _listAvaliablePermissions() {

		return array(
			'console' => array('view', 'commands', 'power'),
			'files' => array('view', 'edit', 'save', 'download', 'delete', 'create', 'upload', 'zip'),
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

		$buildPermissions = array();

		foreach(self::_listAvaliablePermissions() as $permission => $submodule) {

			foreach($submodule as $id => $subpermission) {

				if(!is_array($subpermission)) {
					$buildPermissions[$permission][$subpermission] = $this->hasPermission($permission.".".$subpermission, $array);
				} else {

					foreach($subpermission as $subid => $subsubpermission) {
						$buildPermissions[$permission][$id][$subsubpermission] = $this->hasPermission($permission.".".$id.".".$subsubpermission, $array);
					}

				}

			}

		}

		return $buildPermissions;

	}

	/**
	 * Checks if a given user has permission to access a part of the Control Panel. Defaults to true if the user is the owner.
	 *
	 * @param string $permission The permission node to check aganist.
	 * @param array $array
	 * @return bool Returns true if they have permission, false if not.
	 */
	public function hasPermission($permission, $array = null) {

		if(!is_null($array) && is_array($array)) {
			$this->permissions = $array;
		} else {
			$this->_getPermissions();
		}

		if(self::$permissions_owner != self::getData('id') && is_null($array)) {

			if($this->isAdmin()) {
				return true;
			} else {
				return (is_array($this->permissions)) ? in_array($permission, $this->permissions) : false;
			}

		} else {

			if(!is_array($array)) {
				return true;
			} else {
				return (is_array($this->permissions)) ? in_array($permission, $this->permissions) : false;
			}

		}

	}

}