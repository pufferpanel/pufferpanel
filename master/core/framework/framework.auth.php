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
 * Auth Framework Core
 */

class auth extends dbConn {
	
	public function __construct()
		{
		
			$this->mysql = parent::getConnection();
		
		}
	
	public function encrypt($string, $algo = 'ripemd320'){
		
        $salt = crypt($string, '$6$rounds=5000$'.$this->getSalt().'$');
		return hash($algo, $salt);
	
	}
	
	public function keygen($amount){
		
		$keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
		
		$randkey = "";
		
		for ($i=0; $i<$amount; $i++)
			$randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
		
		return $randkey;
			
	}
	
	public function isLoggedIn($ip, $session, $acp = false){
	
		$this->expires = time();
		
			if($acp !== true){
						
				$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sesip AND `session_id` = :sesid AND `session_expires` > '".$this->expires."'");
				$this->query->execute(array(':sesip' => $ip, ':sesid' => $session));
				
					if($this->query->rowCount() == 1){
					
						$this->update = $this->mysql->prepare("UPDATE `users` SET `session_expires` = :sesexp WHERE `session_ip` = :sesip AND `session_id` = :sesid");
						$this->update->execute(array(':sesexp' => (time() + 1800), ':sesip' => $ip, ':sesid' => $session));
					
						return true;
						
					}else{
					
						return false;
					
					}
					
			}else{
			
				/*
				 * Admin CP Login
				 */
				
				$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sesip AND `session_id` = :sesid AND `session_expires` > '".$this->expires."' AND `root_admin` = 1");
				$this->query->execute(array(':sesip' => $ip, ':sesid' => $session));
				
					if($this->query->rowCount() == 1){
					
						$this->update = $this->mysql->prepare("UPDATE `users` SET `session_expires` = :sesexp WHERE `session_ip` = :sesip AND `session_id` = :sesid");
						$this->update->execute(array(':sesexp' => (time() + 1800), ':sesip' => $ip, ':sesid' => $session));
					
						return true;
				
						
					}else{
					
						return false;
					
					}
				
			
			}
	
	}
	
	public function getCookie($cookie){
	
		if(isset($_COOKIE[$cookie])){
		
			return $_COOKIE[$cookie];
		
		}else{
		
			return 'error';
		
		}
	
	}

}

?>