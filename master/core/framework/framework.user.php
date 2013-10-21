<?php

/*
 * PufferPanel User Class File
 */

class user extends auth {

	public function __construct($ip, $session){
	
		//Re-Initalize the MySQL PDO Class
		$this->mysql = parent::getConnection();
		
		if(auth::isLoggedIn($ip, $session) === true){
		
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