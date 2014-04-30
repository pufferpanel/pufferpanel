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

class apiInitializationClass {

	use \Database\database, \Modules\validate;
	private $data;
	private $keyData;
	
	public function __construct() {
	
		$this->mysql = self::connect();
		$this->data = self::getStoredData();
	
	}
	
	public function init() {
		
		/*
		 * Throw Authentication Errors, otherwise allow script to continue running
		 */
		if(array_key_exists('key', $this->data['auth']) && !empty($this->data['auth']['key'])){
		
			$this->_validateKey();
			self::run();
		
		}else
			self::throwResponse('No API key was provided in the authentication method.', false);
			
	
	}
	
	private function _validateKey() {
	
		$this->_validate = $this->mysql->prepare("SELECT * FROM `api` WHERE `key` = :key");
		$this->_validate->execute(array(
			':key' => $this->data['auth']['key']
		));
		
			if($this->_validate->rowCount() == 0)
				self::throwResponse("Invalid API key provided.", false);
			else
				$this->keyData = $this->_validate->fetch();
		
		$this->_validateRequestIP();
		$this->_validateRequestPermissions();
		
	}
	
	private function _validateRequestIP() {
	
		$this->allowedIPs = json_decode($this->keyData['request_ips'], true);

		/*
		 * Is Key Locked to Specific IPs
		 */
		if($this->allowedIPs[0] != "*") {
		
			$this->hitMatch = false;
			
			/*
			 * Run Through List
			 */
			foreach($this->allowedIPs as $ip)
				if($_SERVER["REMOTE_ADDR"] == $ip)
					$this->hitMatch = true;
					
			/*
			 * No Match
			 */
			if($this->hitMatch === false)
				self::throwResponse("This IP Address is not permitted to access the API in this manner.", false);
			 
		}
	
	}
	
	
	private function _validateRequestPermissions() {
	
		/*
		 * @TODO: Permissions Setup
		 */
	
	}

}

?>