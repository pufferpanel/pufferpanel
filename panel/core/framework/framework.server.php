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
 $core->server = new server($core->auth->getCookie('pp_server_hash'), $core->user->getData('id'), $core->user->getData('root_admin'));
 */

class server extends user {

	use Page\components;
	
	public function __construct($hash, $userid, $isroot){
		
		if($userid !== false && !empty($hash)){
		
			$this->mysql = self::connect();
			
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
	        
	        /*
	         * Grab Node Information
	         */
	        if($this->_data['node'] !== false){
	            
	            $this->_ndata = array();
	            $this->_n = true;
	            
	            $this->query->node = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :node LIMIT 1");
	            $this->query->node->execute(array(
	                ':node' => $this->_data['node'] 
	            ));
	            
	            if($this->query->node->rowCount() == 1){
	                
	                $this->node = $this->query->node->fetch();
	        
	                    foreach($this->node as $this->nid => $this->nval){
	            
	                        $this->_ndata = array_merge($this->_ndata, array($this->nid => $this->nval));
	            
	                    }
	                    
	            }else{
	            
	                $this->_n = false;
	                
	            }
	            
	        }else{
	            
	            $this->_n = false;
	            
	        }
        
        }else{
        
        	$this->_s = false;
        	
        }
	
	}
	
	public function getData($id){
	
		if($this->_s === true && array_key_exists($id, $this->_data)){
		
			return $this->_data[$id];
		
		}else{
		
			return false;
		
		}

	}
    
    public function nodeData($id) {
        
        if($this->_n === true && array_key_exists($id, $this->_ndata)){
		
			return $this->_ndata[$id];
		
		}else{
		
			return null;
		
		}
        
    }
    
    public function nodeRedirect($hash) {
    	
		if($this->user->getData('root_admin') == '1'){
			$query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = ? AND `active` = '1'");
			$query->execute(array($hash));
		}else{
			$query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :ownerid AND `hash` = :hash AND `active` = '1'");
			$query->execute(array(
				':ownerid' => user::getData('id'),
				':hash' => $hash
			));
		}
		
			if($query->rowCount() == 1){
			
				$row = $query->fetch();
				
					setcookie('pp_server_hash', $row['hash'], 0, '/', $this->settings->get('cookie_website'));
				
					$this->redirect('node/index.php');
			
			}else{
			
				$this->redirect('servers.php?error=error');
			
			}
	
	}

}

?>