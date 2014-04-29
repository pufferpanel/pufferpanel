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

class query {

	use Database\database;
	
	public function __construct($serverid){
		
		if($serverid === false)
			$this->_queryData = false;
		else {
		
			$this->mysql = self::connect();
			$this->gsid = (int)$serverid;
			$this->_queryData = array();
			$this->_nodeData = array();
			
			/*
			 * Load Information into Script
			 */
			$this->executeQuery = $this->mysql->prepare("SELECT * FROM `servers` WHERE `id` = :sid");
			$this->executeQuery->execute(array(
				':sid' => $this->gsid
			));
			
				if($this->executeQuery->rowCount() == 1){
				
					 $this->row = $this->executeQuery->fetch();
					 
					 foreach($this->row as $this->id => $this->val)
					 	$this->_queryData = array_merge($this->_queryData, array($this->id => $this->val));
					 	
				}else
					$this->_queryData = false;
					
			/*
			 * Load Node Information into Script
			 */
			$this->executeNodeQuery = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :nid");
			$this->executeNodeQuery->execute(array(
				':nid' => $this->_queryData['node']
			));
			
				if($this->executeNodeQuery->rowCount() == 1){
				
					 $this->node = $this->executeNodeQuery->fetch();

					 foreach($this->node as $this->id => $this->val)
					 	$this->_nodeData = array_merge($this->_nodeData, array($this->id => $this->val));
					 	
				}else{
					$this->_nodeData = false;
					exit('here');	
				}
										
		}
		
	}
	
	/*
	 * Get status of any specificed server
	 */
	public function check_status($ip, $id, $secret){
	
		$this->context = stream_context_create(array(
			"http" => array(
				"method" => "GET",
				"timeout" => 3,
				"header" => "X-Access-Token: ".$secret
			)
		));
		$this->gatherData = @file_get_contents("http://".$ip.":8003/gameservers/".$id , 0, $this->context);

		$this->raw = json_decode($this->gatherData, true);
		
			/*
			 * Valid Data was Returned
			 */
			if(!$this->gatherData)
				return false;
			else
				if(json_last_error() == JSON_ERROR_NONE)
					if($this->raw['status'] == 0)
						return false;
					else
						return true;
				else
					return false;
	
	}
	
	/*
	 * Get status of currently selected server
	 */
	public function online() {
		
		if($this->_queryData === false)
			return false;
		else {
			
			/*
			 * Allow Override of Specific Server (used for pinging)
			 */
			$this->_queryData['gsd_id'] = ($override !== false) ? (int)$override : $this->_queryData['gsd_id'];
			
			$this->context = stream_context_create(array(
				"http" => array(
					"method" => "GET",
					"timeout" => 3,
					"header" => "X-Access-Token: ".$this->_nodeData['gsd_secret']
				)
			));
			$this->gatherData = @file_get_contents("http://".$this->_nodeData['sftp_ip'].":8003/gameservers/".$this->_queryData['gsd_id'] , 0, $this->context);
		
			$this->raw = json_decode($this->gatherData, true);
			
				/*
				 * Valid Data was Returned
				 */
				if(!$this->gatherData)
					return false;
				else{
				
					if(json_last_error() == JSON_ERROR_NONE){
					
						if($this->raw['status'] == 0)
							return false;
						else{
							$this->_jsonData = $this->raw['query'];
							$this->_serverPID = $this->raw['pid'];
							$this->_jsonProcess = $this->raw['process'];
							return true;
						}
					
					}else
						return false;
						
				}
		
		}
		
	
	}
	
	public function pid() {
		
		if($this->online() === true)
			return $this->_serverPID;
		else
			return null;
		
	}
	
	public function retrieve_process($element = null) {
		
		if($this->online() === true)
			if(is_null($element))
				return $this->_jsonProcess;
			else
				return (array_key_exists($element, $this->_jsonProcess)) ? $this->_jsonProcess[$element] : null;
		else
			return null;
		
	}
	
	public function retrieve($element = null) {
		
		if($this->online() === true)
			if(is_null($element))
				return $this->_jsonData;
			else
				return (array_key_exists($element, $this->_jsonData)) ? $this->_jsonData[$element] : null;
		else
			return null;
		
	}

}