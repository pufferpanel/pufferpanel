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

	use \PufferPanel\Core\Components\Authentication;

	/**
	 * @param array $nodesData
	 */
	protected $_nodesData = array();

	/**
	* @param array $_IPArray
	*/
	protected $_IPArray = array();

	/**
	* @param array $_PortArray
	*/
	protected $_PortArray = array();

	/**
	 * @param array $_requiredAddFields Array containing the names of all fields required to be included when adding a node.
	 */
	protected $_requiredAddFields = array('node', 'ip', 'ips');

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
	 * Collects and returns data about a single node.
	 *
	 * @param string $id ID of the node to return data about.
	 * @return array
	 */
	public function getNode($id) {

		$this->node = ORM::forTable('users')->rawQuery("SELECT nodes.*, GROUP_CONCAT(servers.hash) AS servers FROM nodes LEFT JOIN servers ON servers.node = nodes.id WHERE nodes.id = :id LIMIT 1", array('id' => $id))->findOne();

		if(is_null($this->node->id))
			return false;
		else {

			return array(
				"id" => (int) $this->node->id,
				"node" => $this->node->node,
				"fqdn" => $this->node->fqdn,
				"ip" => $this->node->ip,
				"gsd_listen" => (int) $this->node->gsd_listen,
				"gsd_console" => (int) $this->node->gsd_console,
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
	public function getNodes() {

		$this->nodes = ORM::forTable('nodes')->findMany();

		foreach($this->nodes as &$this->node) {

			$this->_nodesData[$this->node->id] = array(
				"node" => $this->node->node,
				"fqdn" => $this->node->fqdn,
				"ip" => $this->node->ip,
				"ports" => json_decode($this->node->ips, true),
			);

		}

		return $this->_nodesData;

	}

	/**
	 * Collects and returns data about all servers in the system.
	 *
	 * @param array $data An array of all data to use to add the node to the database.
	 * @return bool
	 */
	public function addNode($data = array()) {

		foreach($_requiredAddFields as $this->_id) {

			if(!in_array($this->_id, $data))
				return 1;
			else {

				if($this->_id == 'node') {

					if(!preg_match('/^[\w.-]{1,15}$/', $data['node']))
						return 2;
					else
						$this->node = $data['node'];

				}

				if($this->_id == 'ip'){

					if(!filter_var(gethostbyname($data['fqdn']), FILTER_VALIDATE_IP))
						return 3;
					else
						$this->ip = $data['ip'];

				}

				if($this->_id == 'ips') {

					$this->_buildPortArray($data['ips']);

					if($this->_IPArray === false)
						return 4;

					$this->ips = $this->_IPArray;
					$this->ports = $this->_PortArray;

				}

			}

		}

		$this->fqdn = (array_key_exists('fqdn', $data) && filter_var(gethostbyname($data['fqdn']), FILTER_VALIDATE_IP)) ? $data['fqdn'] : $data['ip'];
		$this->gsd_listen = (array_key_exists('gsd_listen', $data) && is_numeric($data['gsd_listen'])) ? $data['gsd_listen'] : 8003;
		$this->gsd_console = (array_key_exists('gsd_console', $data) && is_numeric($data['gsd_console'])) ? $data['gsd_console'] : 8031;
		$this->gsd_server_dir = (array_key_exists('gsd_server_dir', $data) && preg_match('/^([\/][\d\w.\-\/]+[\/])$/', $data['gsd_server_dir'])) ? $data['gsd_server_dir'] : '/home/';

		// Add to Database
		$this->insert = ORM::forTable('nodes')->create();
		$this->insert->set(array(
			'node' => $this->node,
			'fqdn' => $this->fqdn,
			'ip' => $this->ip,
			'gsd_secret' => $this->generateUniqueUUID('nodes', 'gsd_secret'),
			'gsd_listen' => $this->gsd_listen,
			'gsd_console' => $this->gsd_console,
			'gsd_server_dir' => $this->gsd_server_dir,
			'ips' => json_encode($this->_IPArray),
			'ports' => json_encode($this->_PortArray)
		));
		$this->insert->save();

		return array(
			"id" => $this->insert->id(),
			"node" => $this->node
		);

	}

	/**
	* Builds IP and Port array for a node given raw data.
	*
	* @param string $raw Raw string of IPs and Ports to build into an array
	* @return bool
	*/
	protected function _buildPortArray($raw) {

		$raw = str_replace(" ", "", $raw);

		$lines = explode("\r\n", $raw);
		foreach($lines as $id => $values) {

			list($ip, $ports) = explode('|', $values);

			$portList = array();
			$this->_this->_IPArray = array_merge($this->_IPArray, array($ip => array()));
			$this->_PortArray = array_merge($this->_PortArray, array($ip => array()));

			if(!strpos($ports, ",")) {

				// Possible a Range, or a Single Port
				if(!strpos($ports, "-")) {

					$portList[$ports] = 1;

				} else {

					// Range of Ports
					$portRange = explode('-', $ports);
					$r1 = intval($portRange[0]);
					$r2 = intval($portRange[1]);

					$i = 0;
					while($r1 <= $r2) {

						$portList[$i] = $r1;
						$i++;
						$r1++;

					}

				}

			} else {

				// Possible Mix of Ranges and Single Ports
				if(!strpos($ports, "-")) {

					// No ranges
					$portList = array_merge($portList, explode(",", $ports));

				} else {

					// Singles Mixed with Range
					foreach(explode(",", $ports) as $id => $range){

						if(!strpos($range, "-")) {

							$portList = array_merge($portList, array($range));

						} else {

							// Range of Ports
							$portRange = explode('-', $range);
							$r1 = intval($portRange[0]);
							$r2 = intval($portRange[1]);

							$i = 0;
							while($r1 <= $r2) {

								$portList = array_merge($portList, array($r1));
								$i++;
								$r1++;

							}

						}

					}

				}

			}

			for($l=0; $l<count($portList); $l++) {

				if(!array_key_exists($l, $this->_PortArray[$ip])) {
					$this->_PortArray[$ip][$portList[$l]] = 1;
				}

			}

			if(count($this->_PortArray[$ip]) > 0) {
				$this->_IPArray[$ip] = array_merge($this->_IPArray[$ip], array("ports_free" => count($this->_PortArray[$ip])));
			} else {
				$this->_IPArray = false;
			}

		}

	}

}