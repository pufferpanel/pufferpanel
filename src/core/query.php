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
use \ORM, \Unirest;

/**
 * PufferPanel Core GSD Implementation Class
 */
class Query extends Server {

	use Components\GSD;

	protected $server = false;

	protected $node = false;

	protected $query = false;

	protected $query_data;

	/**
	 * Builds server data using a specified ID, Hash, and Root Administrator Status.
	 *
	 * @param int $serverid The ID of the server that will be implementing this class.
	 * @return array Returns an array on success or false on failure.
	 */
	public function __construct() {

		Server::__construct();
		if(is_numeric($this->getData('id'))) {

			/*
			 * Load Information into Script
			 */
			$this->server = ORM::forTable('servers')->findOne($this->getData('id'));

			/*
			 * Load Node Information into Script
			 */
			$this->node = ORM::forTable('nodes')->findOne($this->server->node);

		}

	}

	/**
	 * Gets the status of any specified server given an IP address.
	 *
	 * @param string $ip The IP address of the main GSD server to check aganist.
	 * @param int $port The port of the main GSD server to check aganist.
	 * @param int $id The GSD ID of the server to check.
	 * @param string $secret The GSD secret of the server to check, or the god token for the node.
	 * @return bool Returns an true if server is on, false if off or invalid data was recieved.
	 */
	public function check_status($ip, $port, $id, $secret){

		try {

			Unirest\Request::timeout(1);
			$request = Unirest\Request::get(
				"http://".$ip.":".$port."/gameservers/".$id,
				array(
					'X-Access-Token' => $secret
				)
			);

		} catch(\Exception $e) {
			return false;
		}

		return (isset($request->body->status)) ? $request->body->status : false;

	}

	/**
	* Gets the status of any specified node.
	*
	* @param string $ip The IP address of the main GSD server to check aganist.
	* @param int $port The port of the main GSD server to check aganist.
	* @return bool
	*/
	public function checkNodeStatus($ip, $port = 8003){

		try {

			Unirest\Request::timeout(1);
			$request = Unirest\Request::get(
				"http://".$ip.":".$port."/gameservers/"
			);

			return true;

		} catch(\Exception $e) {
			return false;
		}

	}

	/**
	 * Gets the status of any specified server given an IP address.
	 *
	 * @return bool Returns an true if server is on, false if off or invalid data was recieved.
	 */
	public function online() {

		if(!$this->server || !$this->node) {
			return false;
		}

		try {

			Unirest\Request::timeout(1);
			$request = Unirest\Request::get(
				"http://".$this->node->ip.":".$this->node->port."/gameservers/".$this->server->gsd_id,
				array(
					'X-Access-Token' => $this->node->gsd_secret
				)
			);

			/*
				* Valid Data was Returned
				*/
			if(!isset($request->body->status) || $request->body->status == 0) {
				return false;
			}

			$this->query_data = $request->body;
			return true;

		} catch(\Exception $e) {
			return false;
		}

	}

	/**
	 * Returns the process ID of the last server query.
	 *
	 * @deprecated
	 * @return int
	 */
	public function pid() {

		return null;

	}

	/**
	* Returns data about the last server query.
	*
	* @param string $element A specific part of the JSON to return, if not provided the entire array is returned.
	* @return mixed
	*/
	public function retrieve($element = null) {

		if($this->online()) {

			if(is_null($element)) {
				return $this->query_data;
			}

			if(isset($this->query_data->{$element})) {
				return $this->query_data->{$element};
			}

		}

		return null;

	}

	/**
	 * Returns the last x lines from the server log for a server.
	 *
	 * @param int $lines The number of lines from the server log to return.
	 * @return string
	 */
	public function serverLog($lines = 750) {

		if($this->online()) {

			try {

				$response = Unirest\Request::get(
					"http://".$this->node->ip.":".$this->node->gsd_listen."/gameservers/".$this->server->gsd_id."/log/".$lines,
					array(
						"X-Access-Token" => $this->server->gsd_secret
					)
				);
				return $response->body;

			} catch(\Exception $e) {

				\Tracy\Debugger::log($e);
				return "Daemon error occured.";

			}

		} else {
			return "Server is currently offline.";
		}

	}

	/**
	 * Generates a server.properties file from a given template if it does not already exist.
	 *
	 * @return bool|string
	 */
	public function generateServerProperties() {

		$this->response = $this->_getServerProperties();

		if(!$this->response) {
			\Tracy\Debugger::log($this->node);
			return "Unable to connect to the GSD Daemon running on the node.";
		}

		if(!in_array($this->response->code, array(200, 500))) {

			switch($this->response->code) {

				case 403:
					return "Authentication error encountered.";
					break;
				default:
					return "[HTTP/{$this->response->code}] Invalid response was recieved. ({$this->response->raw_body})";
					break;

			}

		}

		if($this->response->code == 500 || !isset($this->response->body->contents) || empty($this->response->body->contents)) {

			/*
			* Create server.properties
			*/
			if(!file_exists(APP_DIR.'templates/server.properties.tpl') || empty(file_get_contents(APP_DIR.'templates/server.properties.tpl'))) {

				return "No Template Avaliable for server.properties";

			}

			try {

				$this->put = Unirest\Request::put(
					"http://".$this->node->ip.":".$this->node->gsd_listen."/gameservers/".$this->server->gsd_id."/file/server.properties",
					array(
						"X-Access-Token" => $this->server->gsd_secret
					),
					array(
						"contents" => sprintf(file_get_contents(APP_DIR.'templates/server.properties.tpl'), $this->server->server_port, $this->server->server_ip)
					)
				);

			} catch(\Exception $e) {

				\Tracy\Debugger::log($e);
				return "An error occured when trying to write a server.properties file.";

			}

			if(!empty($this->put->body)) {
				return "Unable to process request to create server.properties file.";
			}

		}

		return true;

	}

	/**
	* Turns on a server.
	*
	* @return bool
	*/
	public function powerOn() {

		try {

			$request = Unirest\Request::get(
				"http://".$this->node->ip.":".$this->node->gsd_listen."/gameservers/".$this->server->gsd_id."/on",
				array(
					"X-Access-Token" => $this->server->gsd_secret
				)
			);

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			return false;

		}

		return ($request->body != "ok") ? false : true;

	}

	/**
	 * Returns the contents of server.properties as a \Unirest object or false if the file doesn't exist.
	 *
	 * @return bool|object
	 */
	public function _getServerProperties() {

		if(!$this->node || !$this->server) {
			return false;
		}

		try {

			return Unirest\Request::get(
				"http://".$this->node->ip.":".$this->node->gsd_listen."/gameservers/".$this->server->gsd_id."/file/server.properties",
				array(
					"X-Access-Token" => $this->server->gsd_secret
				)
			);

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			return false;

		}

	}

}
