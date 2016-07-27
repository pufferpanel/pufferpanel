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
 * PufferPanel Core Daemon Implementation Class
 */
class Daemon extends Server {

	use Components\Daemon;

	protected $server = false;

	protected $node = false;

	protected $query = false;

	protected $query_data;

	/**
	 * Builds server data using a specified ID, Hash, and Root Administrator Status.
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
	 * @param string $ip The IP address of the main daemon server to check aganist.
	 * @param int $port The port of the main daemon server to check aganist.
	 * @param string $hash The hash of the server to check for.
	 * @param string $secret The daemon secret of the server to check, or the god token for the node.
	 * @return bool Returns an true if server is on, false if off or invalid data was recieved.
	 */
	public function check_status() {

		$arguments = func_get_args();

		try {

			Unirest\Request::timeout(1);
			$request = Unirest\Request::get(
				"https://".$arguments[0].":".$arguments[1]."/server/%s/status",
				array(
					'X-Access-Token' => $arguments[3],
					'X-Access-Server' => $arguments[2]
				)
			);

		} catch(\Exception $e) {
			return false;
		}

		return (isset($request->body->status)) ? $request->body->status : false;

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
			$request = Unirest\Request::get(sprintf("https://%s:%s/server/%s/status?privkey=%s",
                                $this->node->fqdn,
                                $this->node->daemon_listen,
                                $this->server->hash,
                                $this->node->daemon_secret)
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

				$response = Unirest\Request::get(sprintf("https://%s:%s/server/%s/log?lines=%s&privkey=%s",
                                        $this->node->fqdn,
                                        $this->node->daemon_listen,
                                        $this->server->hash,
                                        $lines,
                                        $this->node->daemon_secret)
                                );

				if($response->code !== 200) {
					return isset($response->body->message) ? "An error occured using the daemon. Said '".$response->body->message."'" : "An unknown error occured with the daemon.";
				}

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

		$response = $this->_getServerProperties();

		if(!$response) {
			\Tracy\Debugger::log($this->node);
			return "Unable to connect to the Scales Daemon running on the node.";
		}

		if(!in_array($response->code, array(200, 500))) {

			switch($response->code) {

				case 403:
					return "Authentication error encountered.";
				default:
					return "[HTTP/{$response->code}] Invalid response was recieved. ({$response->raw_body})";

			}

		}

		if($response->code == 500 || !isset($response->body->contents) || empty($response->body->contents)) {

			/*
			* Create server.properties
			*/
			if(!file_exists(APP_DIR.'templates/server.properties.tpl') || empty(file_get_contents(APP_DIR.'templates/server.properties.tpl'))) {

				return "No Template Avaliable for server.properties";

			}

			try {

				Unirest\Request::put(
					"https://".$this->node->fqdn.":".$this->node->daemon_listen."/server/file/server.properties",
					array(
						"X-Access-Token" => $this->server->daemon_secret,
						"X-Access-Server" => $this->server->hash
					),
					array(
						"contents" => sprintf(file_get_contents(APP_DIR.'templates/server.properties.tpl'), $this->server->server_port, $this->server->server_ip)
					)
				);

			} catch(\Exception $e) {

				\Tracy\Debugger::log($e);
				return "An error occured when trying to write a server.properties file.";

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
                        sprintf("http://%s:%s/server/%s/start/%s?privkey=%s",
                                $this->server->nodeData('fqdn'),
                                $this->server->nodeData('daemon_listen'),
                                $this->server->getData('hash'),
                                $this->server->getData('daemon_secret')
                                )
			);

		} catch(\Exception $e) {

			//\Tracy\Debugger::log($e);
			return false;

		}

		return ($request->code != 200) ? false : true;

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
				"https://".$this->node->fqdn.":".$this->node->daemon_listen."/server/file/server.properties",
				array(
					"X-Access-Token" => $this->server->daemon_secret,
					"X-Access-Server" => $this->server->hash
				)
			);

		} catch(\Exception $e) {

			\Tracy\Debugger::log($e);
			return false;

		}

	}

}