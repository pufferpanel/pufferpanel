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
 * PufferPanel Core Server management class.
 */
class Server extends User {

	use Components\Page;

	/**
	 * @param bool $found_server
	 */
	protected $found_server = true;

	/**
	 * @param bool $found_node
	 */
	protected $found_node = true;

	/**
	 * @param array $server
	 */
	protected $server;

	/**
	 * @param array $node
	 */
	protected $node;

	/**
	 * Constructor class for building server data.
	 *
	 * @param string|int $reference This can either be a string (server hash) or a numeric value (server id).
	 * @return false|null
	 */
	public function __construct($reference = null) {

		Authentication::__construct();
		parent::__construct();

		if(is_null($reference) && !isset($_COOKIE['pp_server_hash'])) {
			return false;
		} else if(is_null($reference)) {
			$reference = $_COOKIE['pp_server_hash'];
		}

		if(is_numeric($reference)) {
			$this->_rebuildData($reference);
		} else {
			$this->_buildData($reference);
		}
		
		parent::initalizePermissions($this->getData('hash'), $this->getData('owner_id'));

	}

	/**
	 * Re-runs the __construct() class with a defined ID for the admin control panel.
	 *
	 * @param int $id This value should be the ID of the server you are getting information for.
	 * @return void
	 */
	public function rebuildData($id) {

		self::__construct($id);

	}

	/**
	 * Provides the corresponding value for the id provided from the MySQL Database.
	 *
	 * @param string $id The column value for the data you need (e.g. server_name).
	 * @return mixed A string is returned on success, array if nothing was passed, and if the command fails 'false' is returned.
	 */
	public function getData($id = null) {


		if(is_null($id) && !is_null($this->server)) {

			$reflect = new ReflectionClass($this->server);
			$data = $reflect->getProperty('_data');
			$data->setAccessible(true);

			return ($this->found_server) ? $data->getValue($this->server) : false;

		} else {
			return ($this->found_server && isset($this->server->{$id})) ? $this->server->{$id} : false;
		}

	}

	/**
	 * Returns data about the node in which the server selected is running.
	 *
	 * @param string $id The column value for the data you need (e.g. ip).
	 * @return mixed A string is returned on success, array if nothing was passed, and if the command fails 'false' is returned.
	 */
	public function nodeData($id = null) {

		if(is_null($id) && !is_null($this->node)) {

			$reflect = new ReflectionClass($this->node);
			$data = $reflect->getProperty('_data');
			$data->setAccessible(true);

			return ($this->found_node) ? $data->getValue($this->node) : false;

		} else {
			return ($this->found_node && isset($this->node->{$id})) ? $this->node->{$id} : false;
		}

	}

	/**
	 * Handles incoming requests to access a server and redirects to the correct location and sets a cookie.
	 *
	 * @param string $hash The hash of the server you are redirecting to.
	 * @return boolean
	 */
	public function nodeRedirect($hash) {

		$query = ORM::forTable('servers')->where(array('hash' => $hash, 'active' => 1));

		if(!$this->isAdmin()) {
			$query = $query->where_raw('`owner_id` = ? OR `hash` IN(?)', array(parent::getData('id'), join(',', parent::listServerPermissions())));
		}

		return (!$query->findOne()) ? false : true;

	}

	/**
	 * Rebuilds server data using a specified ID. Useful for Admin CP applications.
	 *
	 * @param int $server_id The server ID.
	 * @return mixed Returns an array on success or false on failure.
	 */
	private function _rebuildData($server_id) {

		$this->server = ORM::forTable('servers')->findOne($server_id);

		if(!$this->server) {
			$this->found_server = false;
			$this->found_node = false;
			return;
		}

		$this->node = ORM::forTable('nodes')->findOne($this->server->node);

		if(!$this->node) {
			$this->found_node = false;
			return;
		}

	}

	/**
	 * Builds server data using a specified ID, Hash, and Root Administrator Status.
	 *
	 * @param string $hash The server hash.
	 * @return mixed Returns an array on success or false on failure.
	 */
	private function _buildData($hash) {
		
		$query = ORM::forTable('servers')->where(array(
					'hash' => $hash,
					'active' => 1
				));

		if(!$this->isAdmin()) {

			$query = $query->where_raw('`owner_id` = ? OR `hash` IN(?)', array(
					parent::getData('id'),
					join(',', parent::listServerPermissions())
				));
			
		}

		$this->server = $query->findOne();

		if(!$this->server) {
			$this->found_server = false;
			$this->found_node = false;
			return;
		}

		$this->node = ORM::forTable('nodes')->where(array('id' => $this->server->node))->findOne();

		if(!$this->node) {
			$this->found_node = false;
			return;
		}

	}

	/**
	 * Returns an array of users with access to currently active server and their current status.
	 *
	 * @return array
	 */
	public function listAffiliatedUsers() {

		$affiliated = json_decode($this->getData('subusers'), true);
		$userdata = array();

		if(is_array($affiliated) && !empty($affiliated)) {

			foreach($affiliated as $id => $status) {

				if($status != "verified") {

					$selectUser = ORM::forTable('account_change')
						->select('content')
						->where(array('key' => $status, 'verified' => 0))
						->findOne();

					if($selectUser) {

						$content = json_decode($selectUser->content, true);
						$userdata[$id] = array(
							"status" => "pending",
							"revoke" => $status,
							"permissions" => $content[$this->getData('hash')]['perms']
						);

					}

				} else {

					$selectUser = ORM::forTable('users')->selectMany('permissions', 'email', 'uuid')->where('id', $id)->findOne();

					if($selectUser) {

						$permissions = json_decode($selectUser->permissions, true);
						$userdata[$selectUser->email] = array(
							"status" => "verified",
							"id" => $selectUser->uuid,
							"permissions" => $permissions[$this->getData('hash')]['perms']
						);

					}

				}

			}

		}

		return $userdata;

	}

}