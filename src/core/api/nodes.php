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
class Nodes {

	/**
	 * @param array $nodesData Empty array to initalize listig of all nodes.
	 */
	protected $nodesData = array();

	/**
	 * @param array $_requiredAddFields Array containing the names of all fields required to be included when adding a node.
	 */
	protected $_requiredAddFields = array('node', 'ip', 'ips', 'ports');

	/**
	 * @param array $_optionalAddFields Array containing the names of all fields that are optional when adding a node.
	 */
	protected $_optionalAddFields = array('fqdn', 'gsd_listen', 'gsd_console', 'gsd_server_dir');

	/**
	 * Constructor Class
	 * @return void
	 */
	public function __construct() { }

	/**
	 * Intermediate class for handling node data listing data calls.
	 *
	 * @param string $uuid UUID of user to return data about.
	 * @return array
	 */
	public function listNodes($id = null) {

		if(is_null($id) || !is_numeric($id))
			return $this->allNodeData();
		else
			return $this->singleNodeData($id);

	}

	/**
	 * Collects and returns data about a single node.
	 *
	 * @param string $id ID of the node to return data about.
	 * @return array
	 */
	protected function singleNodeData($id) {

		$this->node = ORM::forTable('users')->rawQuery("SELECT nodes.*, GROUP_CONCAT(servers.hash) AS servers FROM nodes LEFT JOIN servers ON servers.node = nodes.id WHERE nodes.id = :id LIMIT 1", array('id' => $id))->findOne();

		if(is_null($this->node->id))
			return false;
		else {

			return array(
				"id" => (int) $this->node->id,
				"node" => $this->node->node,
				"fqdn" => $this->node->fqdn,
				"ip" => $this->node->ip,
				"gsd_listen" => $this->node->gsd_listen,
				"gsd_console" => $this->node->gsd_console,
				"gsd_server_dir" => $this->node->gsd_server_dir,
				"ports" => json_decode($this->node->ports, true),
				"servers" => (!empty($this->node->servers)) ? explode(',', $this->node->servers) : array()
			);

		}

	}

	/**
	 * Collects and returns data about all servers in the system.
	 *
	 * @return array
	 */
	protected function allNodeData() {

		$this->nodes = ORM::forTable('nodes')->findMany();

		foreach($this->nodes as &$this->node){

			$this->nodesData[$this->node->id] = array(
				"node" => $this->node->node,
				"fqdn" => $this->node->fqdn,
				"ip" => $this->node->ip,
				"ports" => json_decode($this->node->ips, true),
			);

		}

		return $this->nodesData;

	}

	/**
	 * Collects and returns data about all servers in the system.
	 *
	 * @param array $data An array of all data to use to add the node to the database.
	 * @return bool
	 */
	public function addNode($data = array()) {



	}

}