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

	use \PufferPanel\Core\Components\Functions;

	protected $_serversData = array();

	private $__addServerFields = array('name', 'node', 'owner', 'ip', 'port', 'memory', 'disk', 'cpu');

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

		if(!$this->server || is_null($this->server->id)) {
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

	/**
	 * Adds a new server to PufferPanel and sets up the GSD information.
	 *
	 * @param array $data An array containing all of the data to be used when creating the server.
	 * @return array|int
	 * 		- Returns the UUID of the server that was just created.
	 *		- Returns an integer if the operation failed which is then matched to an error in the API.
	 */
	public function addServer(array $data) {

		$this->data = $data;
		if(count(array_intersect_key(array_flip($__addServerFields), $this->data)) !== count($__addServerFields)) {
			return 1;
		}

		if(!filter_var($this->data['owner'], FILTER_VALIDATE_EMAIL)) {
			return 7;
		}

		$this->node = ORM::forTable('nodes')->findOne($this->data['node']);
		$this->user = ORM::forTable('users')->select('id')->where('uuid', $this->data['owner'])->findOne();
		$this->ips = json_decode($this->node->ips, true);
		$this->ports = json_decode($this->node->ports, true);

		if(!$this->node || !$this->user) {
			return 2;
		}

		if(!@fsockopen($this->node->ip, 8003, $this->num, $this->error, 3)) {
			return 3;
		}

		if(!preg_match('/^[\w-]{4,35}$/', $this->data['name'])) {
			return 4;
		}

		if(!array_key_exists($this->data['ip'], $this->ips) || !array_key_exists($this->data['port'], $this->ports[$this->data['ip']])) {
			return 5;
		}

		if($this->ports[$this->data['ip']][$this->data['port']] == 0) {
			return 6;
		}

		if(!is_numeric($this->data['memory']) || !is_numeric($this->data['disk']) || !is_numeric($this->data['cpu'])) {
			return 8;
		}

		$this->iv = $this->generate_iv();
		$this->data['ftp_password_raw'] = self::keygen(12);
		$this->data['ftp_password'] = $this->encrypt($this->data['ftp_password_raw'], $this->iv);
		$this->data['ftp_username'] = self::generateFTPUsername($this->data['name']);
		$this->serverHash = $this->generateUniqueUUID('servers', 'hash');
		$this->daemonSecret = $this->generateUniqueUUID('servers', 'gsd_secret');

		/*
		* Build Call
		*/
		$this->data['gsd_id'] = $this->_sendToDaemon(array(
			"name" => $this->serverHash,
			"user" => $this->data['ftp_username'],
			"overide_command_line" => "",
			"path" => $this->node->gsd_server_dir.$this->data['ftp_username'],
			"variables" => array(
				"-jar" => 'server.jar',
				"-Xmx" => $this->data['memory']."M"
			),
			"gameport" => $this->data['port'],
			"gamehost" => $this->data['ip'],
			"plugin" => "minecraft",
			"autoon" => false
		), array(
			'ip' => $this->node->ip,
			'listen' => $this->node->gsd_listen,
			'secret' => $this->node->gsd_secret
		));

		/*
		* Update MySQL
		*/
		$this->server = ORM::forTable('servers')->create();
		$this->server->set(array(
			'gsd_id' => $this->data['gsd_id'],
			'hash' => $this->serverHash,
			'gsd_secret' => $this->node->gsd_secret, // @TODO change this to be a unique secret for each server once GSD/GSC is fixed/implemented
			'encryption_iv' => $this->iv,
			'node' => $this->data['node'],
			'name' => $this->data['name'],
			'modpack' => '0000-0000-0000-0',
			'server_jar' => 'server.jar',
			'owner_id' => $this->user->id,
			'max_ram' => $this->data['memory'],
			'disk_space' => $this->data['disk'],
			'cpu_limit' => $this->data['cpu'],
			'date_added' => time(),
			'server_ip' => $this->data['ip'],
			'server_port' => $this->data['port'],
			'ftp_user' => $this->data['ftp_username'],
			'ftp_pass' => $this->data['ftp_password']
		));
		$this->server->save();

		$this->ips[$this->data['ip']]['ports_free']--;
		$this->ports[$this->data['ip']][$this->data['port']]--;

		$this->node->ips = json_encode($this->ips);
		$this->node->ports = json_encode($this->ports);
		$this->node->save();

		return array(
			'hash' => $this->serverHash,
			'owner' => $this->data['owner'],
			'ftp' => array(
				'username' => $this->data['ftp_username'],
				'password' => $this->data['ftp_password_raw']
			)
		);

	}

	/**
	 * Sends data to GSD to add the server.
	 *
	 * @param array $data
	 * @param array $node
	 * @return void
	 */
	protected function _sendToDaemon(array $data, array $node) {

		$data = json_encode($data);

		$ch = curl_init();
		curl_setopt($ch, CURLOPT_URL, 'http://'.$node['ip'].':'.$node['listen'].'/gameservers');
		curl_setopt($ch, CURLOPT_HTTPHEADER, array(
			'X-Access-Token: '.$node['secret']
		));
		curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 5);
		curl_setopt($ch, CURLOPT_TIMEOUT, 5);
		curl_setopt($ch, CURLOPT_POST, true);
		curl_setopt($ch, CURLOPT_POSTFIELDS, "settings=".$data);
		curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
		$content = json_decode(curl_exec($ch), true);
		curl_close($ch);

		return $content;

	}

}