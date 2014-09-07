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

/**
 * PufferPanel Core User Class File
 */
class user extends Auth\auth {

	/**
	 * @param array $_data Implements a blank array for the functions to write to.
	 */
	private $_data;

	/**
	 * @param bool $_l Defaults to true and will be changed to false if there is an error.
	 */
	private $_l;

	/**
	 * Constructor Class responsible for filling in arrays with the data from a specified user.
	 *
	 * @param string $ip The IP address of a user who is requesting the function, or if called from the Admin CP it is the user id.
	 * @param string|null $session The value of the pp_auth_token cookie.
	 * @param string|null $hash The server hash of the requesting user which is used when they are viewing node pages.
	 * @return void
	 */
	public function __construct($ip, $session = null, $hash = null){

		$this->mysql = self::connect();

		/*
		 * Reset Values
		 */
		$this->_data = array();
		$this->_l = true;

		if(self::isLoggedIn($ip, $session, $hash) === true){

			$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sesip AND `session_id` = :sesid");
			$this->query->execute(array(':sesip' => $ip, ':sesid' => $session));

				$this->row = $this->query->fetch();
				foreach($this->row as $this->id => $this->val){

					$this->_data[$this->id] = $this->val;

				}

		}else if(is_null($session) && is_null($hash) && is_numeric($ip)){

			$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `id` = :id");
			$this->query->execute(array(':id' => $ip));

				$this->row = $this->query->fetch();
				foreach($this->row as $this->id => $this->val){

					$this->_data[$this->id] = $this->val;

				}

		}else
			$this->_l = false;

	}

	/**
	 * Re-runs the __construct() class with a defined ID for the admin control panel.
	 *
	 * @param string $id The ID of a user requested in the Admin CP.
	 * @return void
	 */
	public function rebuildData($id){

		$this->__construct($id);

	}

	/**
	 * Provides the corresponding value for the id provided from the MySQL Database.
	 *
	 * @param string|null $id The column value for the data you need (e.g. email).
	 * @return string|array|bool A string is returned on success, array if nothing was passed, and if the command fails 'false' is returned.
	 */
	public function getData($id = null){

		if(is_null($id))
			if($this->_l === true)
				return $this->_data;
			else
				return false;
		else
			if($this->_l === true && array_key_exists($id, $this->_data))
				return $this->_data[$id];
			else
				return false;

	}

}

?>
