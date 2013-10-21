<?php

/*
 * PufferPanel Server Class File
 $core->framework->server = new server($core->framework->auth->getCookie('pp_server_hash'), $core->framework->user->getData('id'), $core->framework->user->getData('root_admin'));
 */

class server extends user {

	public function __construct($hash, $userid, $isroot){

		$this->mysql = parent::getConnection();
		
		$this->_data = array();
		$this->_s = true;
		
			if($isroot == '1'){
				$this->query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = :hash AND `active` = 1");
				$this->query->execute(array(
					':hash' => $hash
				));
			}else{
				$this->query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = :hash AND `owner_id` = :ownerid AND `active` = 1");
				$this->query->execute(array(
					':hash' => $hash,
					':ownerid' => $userid
				));
			}
		
			if($this->query->rowCount() == 1){
			
				$this->row = $this->query->fetch();
		
					foreach($this->row as $this->id => $this->val){
			
						$this->_data = array_merge($this->_data, array($this->id => $this->val));
			
					}
					
			}else{
			
				$this->_s = false;
				
			}
	
	}
	
	public function getData($id){
	
		if($this->_s === true){
		
			return $this->_data[$id];
		
		}else{
		
			return false;
		
		}

	}

}

?>