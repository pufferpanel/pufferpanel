<?php
/*
	PufferPanel - A Game Server Management Panel
	Copyright (c) 2015 Dane Everitt

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
use \ORM, \Unirest, \ReflectionClass;

/**
* PufferPanel Core Permission Implementation Class
*/
class Permissions extends User {

	protected $permissions;

	protected $permission_list = array();

	protected $subuser;

	private $server;

	/**
	 * Constructor class for Permissions
	 *
	 * @param int $server The server ID that we are checking for.
	 */
	public function __construct($server = null) {

		parent::__construct();

		$this->permissions = ORM::forTable('permissions')->where(array(
			'user' => User::getData('id'),
			'server' => $server
		))->findMany();

		self::buildPermissionList();

		$this->subuser = ORM::forTable('subusers')->where(array(
			'user' => User::getData('id'),
			'server' => $server
		))->findOne();

		if($server !== null) {

			$this->server = ORM::forTable('servers')->select('owner_id')->findOne($server);

		}

	}

	/**
	 * Returns a list of servers that a user has permission to access.
	 *
	 * @return array
	 */
	public function listServers() {

		$select = ORM::forTable('servers')
					->raw_query('(SELECT id FROM servers WHERE owner_id = :userid) UNION (SELECT server FROM subusers WHERE user = :userid)',
						array('userid' => User::getData('id'))
					)
					->findMany();

		$servers = [];
		foreach($select as &$select) {
			$servers = array_merge($servers, array($select['server']));
		}

		return $servers;

	}

	/**
	 * Builds a list of permissions that a user has access to.
	 *
	 * @return void
	 */
	protected function buildPermissionList() {

		foreach($this->permissions as &$permission) {
			$this->permission_list = array_merge($this->permission_list, array($permission['permission']));
		}

	}

	/**
	 * Returns the object value for different rows of a subuser.
	 *
	 * @return string
	 */
	public function get($param) {

		return (isset($this->subuser->{$param})) ? $this->subuser->{$param} : false;

	}

	/**
	 * Returns true or false depending on if the user has a specified permission.
	 *
	 * @return bool
	 */
	public function has($permission) {

		if(User::getData('id') == $this->server->owner_id || User::getData('root_admin') == 1) {
			return true;
		}

		return in_array($permission, $this->permission_list);

	}

}