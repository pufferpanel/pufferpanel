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

class GSD_Query extends dbConn {

	public function __construct($serverid){
	
		if($serverid === false)
			$this->_queryData = false;
		else {
		
			$this->mysql = parent::getConnection();
			$this->gsid = (int)$serverid;
			$this->_queryData = array();
			
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
					
		}
	
	}
	
	public function online($override = false) {
	
		if($this->_queryData === false && $override === false)
			return false;
		else {
			
			/*
			 * Allow Override of Specific Server (used for pinging)
			 */
			$this->_queryData['gsd_id'] = ($override !== false) ? (int)$override : $this->_queryData['gsd_id'];
			
			$this->context = stream_context_create(array(
				"http" => array(
					"method" => "GET",
					"timeout" => 3
				)
			));
			$this->gatherData = @file_get_contents("http://".$this->_queryData['ftp_host'].":8003/gameservers/".$this->_queryData['gsd_id'] , 0, $this->context);
		
			$this->raw = json_decode($this->gatherData, true);
			
				/*
				 * Valid Data was Returned
				 */
				if(json_last_error() == JSON_ERROR_NONE){
				
					if($this->raw['status'] == 0)
						return false;
					else{
						$this->_jsonData = $this->raw['query'];
						$this->_jsonProcess = $this->raw['process'];
						return true;
					}
				
				}else
					return false;
		
		}
		
	
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