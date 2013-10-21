<?php

/*
 * Auth Framework Core
 */

class auth extends dbConn {
	
	public function __construct()
		{
		
			$this->mysql = parent::getConnection();
		
		}
	
	public function encrypt($string, $algo = 'ripemd320'){
		
		$salt = crypt($string, '$6$rounds=5000$c2aX7d0JJSVA9*^xj#0sA$');
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
						$this->update->execute(array(':sesexp' => time() + 1800, ':sesip' => $ip, ':sesid' => $session));
					
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
						$this->update->execute(array(':sesexp' => time() + 1800, ':sesip' => $ip, ':sesid' => $session));
					
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