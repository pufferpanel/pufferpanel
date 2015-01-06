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
 * PufferPanel Core Error Logging Class
 *
 * @extends user
 */
class Log extends User {

	/**
	 * @param string $url
	 */
	private $url;

	/**
	 * Constructor class for logging
	 *
	 * @param int $uid The user ID.
	 * @return void
	 */
	public function __construct($uid) {

			$this->uid = ($uid !== false) ? $uid : null;

		}

	/**
	 * Logging Function
	 * action: (example) user.login, user.start_server, admin.remove_server, admin.node_add_port, system.unknown_error
	 * desc: Description of the Error
	 * uid: Sent from login commands where $user is not yet defined. Optional.
	 *
	 * @param int $priority The priority of the error on a sale of 0 to 4.
	 * @param int $viewable Set to 0 to only be viewable by administrators, 1 to be viewed by the user who caused the error.
	 * @param array $data An array of the data that caused the error. Should be in the form Array(action, desc, uid).
	 * @return void
	 */
	public function addLog($priority, $viewable, $data = array()) {

			$this->uid = (!array_key_exists(2, $data)) ? $this->uid : $data[2];
			$this->logger = ORM::forTable('actions_log')->create();

			$this->logger->priority = $priority;
			$this->logger->viewable = $viewable;
			$this->logger->user = $this->uid;
			$this->logger->time = time();
			$this->logger->ip = $_SERVER['REMOTE_ADDR'];
			$this->logger->url = $this->url;
			$this->logger->action = $data[0];
			$this->logger->desc = $data[1];

			$this->logger->save();

		}

	/**
	 * Get the current page URL where the log action is called.
	 *
	 * @return Log
	 */
	public function getUrl() {

			$this->url = (isset($_SERVER['HTTPS']) == 'on' ? 'https' : 'http').'://'.$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI'];
			return $this;

		}

}

?>