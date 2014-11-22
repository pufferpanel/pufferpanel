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
class Servers {

	protected $_serversData = array();

	/**
	* Constructor Class
	* @return void
	*/
	public function __construct() { }

	/**
	* Collects and returns data about a single server.
	*
	* @param string $hash Hash of server to return data about.
	* @return array
	*/
	public function getServer($hash) {

		$this->server = ORM::forTable('servers')->select('servers.*')->select('users.uuid')->join('users', array('users.id', '=', 'servers.owner_id'))->where('hash', $hash)->findOne();

		if(is_null($this->server->id)) {
			return false;
		} else {

			return array(
				"id" => (int) $this->server->id,
				"node" => (int) $this->server->node,
				"owner" => $this->server->uuid,
				"name" => $this->server->name,
				"server_jar" => $this->server->server_jar,
				"active" => (int) $this->server->active,
				"ram" => (int) $this->server->max_ram,
				"disk" => (int) $this->server->disk_space,
				"cpu" => (int) $this->server->cpu_limit,
				"ip" => $this->server->server_ip,
				"port" => (int) $this->server->server_port,
				"ftp_user" => $this->server->ftp_user
			);

		}

	}

	/**
	* Collects and returns data about all servers in the system.
	*
	* @return array
	*/
	public function getServers() {

		$this->servers = ORM::forTable('servers')->select('servers.*')->select('users.uuid')->join('users', array('users.id', '=', 'servers.owner_id'))->findMany();

		foreach($this->servers as &$this->server) {

			$this->_serversData = array_merge($this->_serversData, array(
				$this->server->hash => array(
					"id" => (int) $this->server->id,
					"owner" => $this->server->uuid,
					"name" => $this->server->name,
					"node" => (int) $this->server->node,
					"active" => (int) $this->server->active
				)
			));

		}

		return $this->_serversData;

	}

}