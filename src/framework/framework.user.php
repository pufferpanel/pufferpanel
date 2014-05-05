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
 
/*
 * PufferPanel User Class File
 */

class user extends Auth\auth {
	
	public function __construct($ip, $session, $serverhash = null){

		//Re-Initalize the MySQL PDO Class
		$this->mysql = self::connect();
		
		if(self::isLoggedIn($ip, $session, $serverhash) === true){
		
			$this->_l = true;
			
				$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sesip AND `session_id` = :sesid AND `session_expires` >  :time");
				$this->query->execute(array(':sesip' => $ip, ':sesid' => $session, ':time' => time()));
			
					$this->row = $this->query->fetch();
					foreach($this->row as $this->id => $this->val){
				
						$this->_data[$this->id] = $this->val;
				
					}
		
		}else{
		
			$this->_l = false;
			
		}
	
	}
	
	public function getData($id){
	
		if($this->_l === true){
		
			return $this->_data[$id];
		
		}else{
		
			return false;
		
		}
	
	}

}

?>