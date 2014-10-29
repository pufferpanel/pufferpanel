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
use \ORM as ORM;

/**
 * PufferPanel Core Server management class.
 */
class Server extends User {

	use Components\Page;

	/**
	 * @param array $_data Implements a blank array for the functions to write to.
	 */
	private $_data;

	/**
	 * @param array $_ndata Implements a blank array for the functions to write to. This variable is used for the node part of the code.
	 */
	private $_ndata;

	/**
	 * @param array $_s Defaults to true and will be changed to false if there is an error. This variable is used for the server portion of the code.
	 */
	private $_s;

	/**
	 * @param array $_n Defaults to true and will be changed to false if there is an error. This variable is used for the node part of the code.
	 */
	private $_n;

	/**
	 * Constructor class for building server data.
	 *
	 * @param string $hash The server hash.
	 * @param int $userid The ID of the user who is requesting the server information or in some cases the ID of the server (if called via rebuildData function).
	 * @param int $isroot The root administrator status of the user requesting the server information.
	 * @return void
	 */
	public function __construct($hash = null, $userid = null, $isroot = null){

		/*
		 * Reset Values
		 */
		$this->_data = array();
		$this->_ndata = array();
		$this->_s = true;
		$this->_n = true;

		/*
		 * Make Calls
		 */
		User::initalizePermissions($hash, $userid, null);
		if(!is_null($userid) && is_numeric($userid) && !is_null($hash))
			$this->_buildData($hash, $userid, $isroot);
		else if(!is_null($userid) && is_null($hash))
			$this->_rebuildData($userid);
		else
			$this->_s = false;

		//Re-assign values with owner ID
		User::initalizePermissions($hash, $userid, $this->getData('owner_id'));

	}

	/**
	 * Re-runs the __construct() class with a defined ID for the admin control panel.
	 *
	 * @param int $id This value should be the ID of the server you are getting information for.
	 * @return void
	 */
	public function rebuildData($id){

		$this->__construct(null, $id);

	}

	/**
	 * Provides the corresponding value for the id provided from the MySQL Database.
	 *
	 * @param string $id The column value for the data you need (e.g. server_name).
	 * @return mixed A string is returned on success, array if nothing was passed, and if the command fails 'false' is returned.
	 */
	public function getData($id = null){

		if(is_null($id))
			if($this->_s === true)
				return $this->server;
			else
				return false;
		else
			if($this->_s === true && isset($this->server->{$id}))
				return $this->server->{$id};
			else
				return false;

	}

	/**
	 * Returns data about the node in which the server selected is running.
	 *
	 * @param string $id The column value for the data you need (e.g. ip).
	 * @return mixed A string is returned on success, array if nothing was passed, and if the command fails 'false' is returned.
	 */
	public function nodeData($id = null) {

		if(is_null($id))
			if($this->_n === true)
				return $this->node;
			else
				return false;
		else
			if($this->_n === true && isset($this->node->{$id}))
				return $this->node->{$id};
			else
				return false;

	}

	/**
	 * Handles incoming requests to access a server and redirects to the correct location and sets a cookie.
	 *
	 * @param string $hash The hash of the server you are redirecting to.
	 * @param array $perms
	 * @param int $userid
	 * @param int $rootAdmin
	 * @return void
	 */
	public function nodeRedirect($hash, $userid, $rootAdmin) {

		$this->__construct($hash, $userid, $rootAdmin);
		if($rootAdmin == 1){

			$this->query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = ? AND `active` = '1'");
			$this->query->execute(array($hash));

		}else{

			$this->hashes = array_map(array($this->mysql, 'quote'), parent::listServerPermissions());
			$this->query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :oid OR `hash` IN(".join(',', $this->hashes).") AND `hash` = :hash AND `active` = '1'");
			$this->query->execute(array(
				':oid' => $userid,
				':hash' => $hash
			));

		}

			if($this->query->rowCount() == 1){

				$this->row = $this->query->fetch();
				setcookie('pp_server_hash', $this->row['hash'], 0, '/');
				$this->redirect('node/index.php');

			}else
				$this->redirect('servers.php?error=error');

	}

	/**
	 * Rebuilds server data using a specified ID. Useful for Admin CP applications.
	 *
	 * @param int $sid The server ID.
	 * @return mixed Returns an array on success or false on failure.
	 */
	private function _rebuildData($sid){

		$this->server = ORM::forTable('servers')->findOne($sid);

		if($this->server === false)
			$this->_s = false;

		/*
		 * Grab Node Information
		 */
		if($this->_s !== false){

			$this->node = ORM::forTable('nodes')->findOne($this->_data->node);

			if($this->node === false)
				$this->_n = false;

		}else
			$this->_n = false;

	}

	/**
	 * Builds server data using a specified ID, Hash, and Root Administrator Status.
	 *
	 * @param string $hash The server hash.
	 * @param int $userid The ID of the user who is requesting the server information.
	 * @param int $isroot The root administrator status of the user requesting the server information.
	 * @return mixed Returns an array on success or false on failure.
	 */
	private function _buildData($hash, $userid, $isroot){

		if($isroot == '1')
			$this->server = ORM::forTable('servers')->where(array('hash' => $hash, 'active' => 1))->findOne();
		else
			$this->server = ORM::forTable('servers')->where(array('hash' => $hash, 'active' => 1))->where_raw('`owner_id` = ? OR `hash` IN(?)', array($userid, join(',', parent::listServerPermissions())))->findOne();

		if($this->server === false)
			$this->_s = false;

		/*
		 * Grab Node Information
		 */
		if($this->_s !== false){

			$this->_n = true;

			$this->node = ORM::forTable('nodes')->where(array('id' => $this->server->node))->findOne();

			if($this->node === false)
				$this->_n = false;

		}else
			$this->_n = false;

	}

}

?>
