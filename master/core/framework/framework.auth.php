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
		
		$keyset  = "abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ0123456789";
		
		$randkey = null;
		
		for ($i=0; $i<$amount; $i++)
			$randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
		
		return $randkey;
					
	}

	public function isLoggedIn($ip, $session, $serverhash = null, $acp = false){

		$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sessip AND `session_id` = :session AND `session_expires` > :sesexp");
		$this->query->execute(array(
			':sessip' => $ip,
			':session' => $session,
			':sesexp' => time()
		));

			if($this->query->rowCount() == 1){

				$this->row = $this->query->fetch();

					/*
					 * Accessing Admin Directory
					 */
					if($this->row['root_admin'] != 1 && $acp === true)
						return false;
					else{

						/*
						 * Allow Admins Access to any Server
						 */
						if($this->row['root_admin'] != '1'){	
	
							/*
							 * Validate User is Owner of Server
							 */
							if(!is_null($serverhash)){
							
								$this->_validateServer = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = :shash AND `owner_id` = :ownerid AND `active` = 1");
								$this->_validateServer->execute(array(
									':shash' => $serverhash,
									':ownerid' => $this->row['id']
								));
		
									if($this->_validateServer->rowCount() == 1){
		
										$this->updateUsers = $this->mysql->prepare("UPDATE `users` SET `session_expires` = :sesexp WHERE `session_ip` = :sesip AND `session_id` = :sesid");
										$this->updateUsers->execute(array(
											':sesexp' => time() + 1800,
											':sesip' => $ip,
											':sesid' => $session
		
										));
		
											return true;
		
									}else{
		
										return false;
		
									}
							
							/*
							 * Just Check if they are Logged In
							 */	
							}else{
							
								return true;
								
							}
	
						}else{
	
							$this->updateUsers = $this->mysql->prepare("UPDATE `users` SET `session_expires` = :sesexp WHERE `session_ip` = :sesip AND `session_id` = :sesid");
							$this->updateUsers->execute(array(
								':sesexp' => time() + 1800,
								':sesip' => $ip,
								':sesid' => $session
	
							));
	
								return true;
	
						}
						
					}

			}else{

				return false;

			}

	}
	
	public function getCookie($cookie){
	
		if(isset($_COOKIE[$cookie])){
		
			return $_COOKIE[$cookie];
		
		}else{
		
			return null;
		
		}
	
	}

}

?>