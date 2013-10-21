<?php

/*
 * Auth Framework Core
 */

class auth extends dbConn {
	
	public function __construct()
		{
		
			$this->mysql = parent::getConnection();
		
		}
	
	public function keygen($amount){
		
		$keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
		
		$randkey = "";
		
		for ($i=0; $i<$amount; $i++)
			$randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
		
		return $randkey;
			
	}
	
	public function isLoggedIn($ip, $session, $serverhash){
	
		$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sessip AND `session_id` = :session AND `session_expires` > :sesexp");
		$this->query->execute(array(
			':sessip' => $ip,
			':session' => $session,
			':sesexp' => time()
		));
		
			if($this->query->rowCount() == 1){
				
				$this->row = $this->query->fetch();
				
					/*
					 * Allow Admins Access to any Server
					 */
					if($this->row['root_admin'] != '1'){	
						
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
							
					}else{
					
						$this->updateUsers = $this->mysql->prepare("UPDATE `users` SET `session_expires` = :sesexp WHERE `session_ip` = :sesip AND `session_id` = :sesid");
						$this->updateUsers->execute(array(
							':sesexp' => time() + 1800,
							':sesip' => $ip,
							':sesid' => $session
							
						));
						
							return true;
					
					}
				
			}else{
			
				return false;
			
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