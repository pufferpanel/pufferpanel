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
class Query {

	use Components\GSD;

	protected $server = false;

	/**
	 * Builds server data using a specified ID, Hash, and Root Administrator Status.
	 *
	 * @param int $serverid The ID of the server that will be implementing this class.
	 * @return array Returns an array on success or false on failure.
	 */
	public function __construct($serverid) {

		if($serverid === false) {
			$this->_queryData = false;
		} else {

			$this->gsid = (int)$serverid;

			/*
			 * Load Information into Script
			 */
			$this->server = ORM::forTable('servers')->findOne($this->gsid);

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
	 * @param int $id The GSD ID of the server to check.
	 * @param string $secret The GSD secret of the server to check, or the god token for the node.
	 * @return bool Returns an true if server is on, false if off or invalid data was recieved.
	 */
	public function check_status($ip, $port, $id, $secret){

		$this->context = stream_context_create(array(
			"http" => array(
				"method" => "GET",
				"timeout" => 1,
				"header" => "X-Access-Token: ".$secret
			)
		));
		$this->gatherData = @file_get_contents("http://".$ip.":".$port."/gameservers/".$id , 0, $this->context);

		$this->raw = json_decode($this->gatherData, true);

		if(!$this->gatherData) {
			return false;
		} else if(json_last_error() == JSON_ERROR_NONE) {
			return $this->raw['status'] != 0;
		} else {
			return false;
		}

	}

	/**
	 * Gets the status of any specified server given an IP address.
	 *
	 * @param int $override If an ID is provided this function will check that specific server, otherwise it is assumed to be the loaded server.
	 * @return bool Returns an true if server is on, false if off or invalid data was recieved.
	 */
	public function online($override = false) {

		if($this->server === false) {
			return false;
		} else {

			$this->server->gsd_id = ($override !== false) ? (int)$override : $this->server->gsd_id;

			$this->context = stream_context_create(array(
				"http" => array(
					"method" => "GET",
					"timeout" => 1,
					"header" => "X-Access-Token: ".$this->node->gsd_secret
				)
			));
			$this->gatherData = file_get_contents("http://".$this->node->ip.":".$this->node->gsd_listen."/gameservers/".$this->server->gsd_id , 0, $this->context);

			$this->raw = json_decode($this->gatherData, true);

			/*
				 * Valid Data was Returned
				 */
			if(!$this->gatherData) {
				return false;
			} else if(json_last_error() == JSON_ERROR_NONE) {

				if($this->raw['status'] == 0) {
					return false;
				} else {

					$this->_jsonData = (array_key_exists('query', $this->raw)) ? $this->raw['query'] : null;
					$this->_serverPID = $this->raw['pid'];
					$this->_jsonProcess = (array_key_exists('process', $this->raw)) ? $this->raw['process'] : array();
					return true;
					
				}
			} else {
				return false;
			}
		}

	}

	/**
	 * Returns the process ID of the last server query.
	 *
	 * @return int
	 */
	public function pid() {

		return $this->online() ? $this->_serverPID : null;

	}

	/**
	 * Returns process information for the last server query.
	 *
	 * @param string $element A specific part of the JSON to return, if not provided the entire array is returned.
	 * @return mixed
	 */
	public function retrieve_process($element = null) {

		if($this->online()) {

			if(is_null($element)) {
				return $this->_jsonProcess;
			} else {
				return (array_key_exists($element, $this->_jsonProcess)) ? $this->_jsonProcess[$element] : null;
			}
			
		} else {
			return null;
		}

	}

	/**
	* Returns data about the last server query.
	*
	* @param string $element A specific part of the JSON to return, if not provided the entire array is returned.
	* @return mixed
	*/
	public function retrieve($element = null) {

		if($this->online() === true) {
			
			if(is_null($element)) {
				return $this->_jsonData;
			} else {
				return (array_key_exists($element, $this->_jsonData)) ? $this->_jsonData[$element] : null;
			}

		} else {
			return null;
		}

	}

	/**
	 * Returns the last x lines from the server log for a server.
	 *
	 * @param int $lines The number of lines from the server log to return.
	 * @return array
	 */
	public function serverLog($lines = 750) {

		if($this->online()) {

			try {

				$response = Unirest::get(
					"http://".$this->node->ip.":".$this->node->gsd_listen."/gameservers/".$this->server->gsd_id."/log/".$lines,
					array(
						"X-Access-Token" => $this->server->gsd_secret
					)
				);

			} catch(\Exception $e) {

				\Tracy\Debugger::log($e);
				return array("Daemon error occured.");

			}

			return $response->body;

		} else {
			return array('Server is currently offline.');
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
			return "Unable to connect to the GSD Daemon running on the node.";
		}

		if(!in_array($this->response->code, array(200, 500))) {

			switch($this->response->code) {

				case 403:
					return "Authentication error encountered.";
				default:
					return "[HTTP/{$this->response->code}] Invalid response was recieved. ({$this->response->raw_body})";

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

				$this->put = Unirest::put(
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

			$this->get = Unirest::get(
				"http://".$this->node->ip.":".$this->node->gsd_listen."/gameservers/".$this->server->gsd_id."/on",
				array(
					"X-Access-Token" => $this->server->gsd_secret
				)
			);

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			return false;

		}

		return ($this->get->body != "ok") ? false : true;

	}

	/**
	 * Returns the contents of server.properties as a \Unirest object or false if the file doesn't exist.
	 *
	 * @return bool|object
	 */
	public function _getServerProperties() {

		try {

			return Unirest::get(
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