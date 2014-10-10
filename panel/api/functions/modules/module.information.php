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

class apiModuleGetInformation {

 	use \Database\database, functions;

 	private $node = array();

	public function __construct() {

		$this->mysql = self::connect();

		/*
		 * Load the Data
		 */
		$this->_getNodes();

		/*
		 * Set Values for Functions
		 */
		$this->filterNode = (array_key_exists('filter_node', self::getStoredData()['data'])) ? self::getStoredData()['data']['filter_node'] : null;
		$this->filterIP = (array_key_exists('filter_ip', self::getStoredData()['data'])) ? self::getStoredData()['data']['filter_ip'] : null;

		/*
		 * Determine Function to Call
		 */
		switch(self::getStoredData()['data']['function']){

			case 'list_nodes':
				self::throwResponse($this->_listNodes(), true);
				break;
			case 'list_ips':
				self::throwResponse($this->_listIPs($this->filterNode), true);
				break;
			case 'list_ports':
				self::throwResponse($this->_listPorts($this->filterNode, $this->filterIP), true);
				break;
			default:
				self::throwResponse("No function specified in API request.", false);
				break;

		}

	}

	/*
	 * Generate Data for Other Functions
	 */
	private function _getNodes() {

		/*
		 * List All Nodes
		 */
		$this->query = $this->mysql->prepare("SELECT * FROM `nodes`");
		$this->query->execute(array());

			while($this->data = $this->query->fetch())
				$this->node[$this->data['id']] = $this->data;

	}

	/*
	 * List Nodes
	 */
	private function _listNodes() {

		$this->return = array();

		foreach($this->node as $id => $data)
			$this->return = array_merge($this->return, array(array(
				"id" => $data['id'],
				"node" => $data['node'],
				"fqdn" => $data['fqdn'],
				"ip" => $data['ip']
			)));

		return $this->return;

	}

	/*
	 * List Avaliable IPs
	 */
	private function _listIPs($node = null) {

		$this->return = null;

		/*
		 * List all IPs on all Nodes
		 */
		if(is_null($node)) {

			self::throwResponse("No node ID was passed in the request.", false);

		}else{

			/*
			 * Validate Node
			 */
			if(!array_key_exists($node, $this->node))
				self::throwResponse("Unable to call listIP() due to an invalid node being passed.", false);

			return json_decode($this->node[$node]['ips']);

		}

	}

	/*
	 * List Avaliable Ports
	 */
	private function _listPorts($node = null, $ip = null) {

		$this->ports = array();

		/*
		 * List all Ports & IPs on Node
		 */
		if(!is_null($node) && is_null($ip)){

			$this->ports = array_merge($this->ports, array("node" => $this->node[$node]['id'], "ports" => json_decode($this->node[$node]['ports'], true)));

			return $this->ports;

		}else if(!is_null($node) && !is_null($ip)){

			$this->list = json_decode($this->node[$node]['ports'], true);

			$this->ports = array("node" => $this->node[$node]['id'], "ports" => array($ip => $this->list[$ip]));

			return $this->ports;

		}else{

			self::throwError("You must pass a node in order to list the ports.", false);

		}

	}

}

 ?>
