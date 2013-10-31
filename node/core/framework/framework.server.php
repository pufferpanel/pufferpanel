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