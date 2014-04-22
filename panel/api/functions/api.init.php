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

trait globalInit {
	
	public function throwResponse($text, $success = false){
	
		exit(json_encode(
			array(
				'success' => $success,
				'info' => $text
			)
		));
	
	}
	
	public function getStoredData() {
	
		if(!isset($_POST['request']))
			$this->throwResponse("No data was sent in the request.");
		else
			return json_decode($_POST['request'], true);
	
	}

}

class apiInitializationClass extends getSettings {

	use globalInit;
	
	public function __construct() {
	
		$this->mysql = parent::getConnection();
	
	}
	
	public function init() {
		
		$this->data = $this->getStoredData();
		
		/*
		 * Throw Authentication Errors, otherwise allow script to continue running
		 */
		if(array_key_exists('key', $this->data['auth']) && !empty($this->data['auth']['key'])){
		
			if(getSettings::get('api_key') == $this->data['auth']['key']){
			
				if($_SERVER['REMOTE_ADDR'] != getSettings::get('api_request_ip'))
					self::throwResponse('This server is not permitted to access the API.', false);
			
			}else
				self::throwResponse('An invalid API key was provided for this request.', false);
		
		}else
			self::throwResponse('No API key was provided in the authentication method.', false);
	
	}

}

?>