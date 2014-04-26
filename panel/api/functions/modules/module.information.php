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
 
namespace Modules\Info;

trait validate {

	public static function validateRequest($req) {
	
		/*
		 * Handle Validation
		 */
		if($req == "add") {
		
			/*
			 * Handle Validation for Adding Server
			 */
		
		}
	
	}

}
 
class apiModuleGetInformation {
 
 	use \Database\database, \API\functions, \Functions\general;
 	private $node = array();
 	
	public function __construct() {
	
		$this->mysql = self::connect();
		$this->_listNodes();
		
		self::throwResponse($this->_listPorts(), true);
	
	}
	
	/*
	 * List Nodes
	 */
	private function _listNodes() {
	
		/*
		 * List All Nodes
		 */
		$this->query = $this->mysql->prepare("SELECT * FROM `nodes`");
		$this->query->execute(array());

			while($this->data = $this->query->fetch())
				$this->node[$this->data['id']] = $this->data;
	
	}
	
	/*
	 * List Avaliable IPs
	 */
	private function _listIPs($node = null) {
	
		$this->return = null;
		
		/*
		 * List all IPs on all Nodes
		 */
		if(is_null($node)) {
		
			foreach($this->node as $id => $data)
				$this->return .= $data['ips'];
		
			/*
			 * Return the Data
			 */
			return $this->return;
		
		}else{
		
			/*
			 * Validate Node
			 */
			if(!array_key_exists($node, $this->node))
				self::throwResponse("Unable to call listIP() due to an invalid node being passed.", false);
				
			return $this->node[$node]['ips'];
		
		}
	
	}
	
	/*
	 * List Avaliable Ports
	 */
	private function _listPorts($node = null, $ip = null) {
	
		$this->ports = array();

		/*
		 * List all Ports & IPs
		 */
		if(is_null($node)){
		
			foreach($this->node as $id => $data){
			
				/*
				 * Rebuild Array for Return
				 */
				$this->ports = array_merge($this->ports, array("node" => $data['id'], "ports" => json_decode($data['ports'], true)));
			
			}
			
			/*
			 * Return as JSON
			 */
			return $this->ports;
		
		}else{
		
			$this->ports = array("node" => $data[$node]['id'], "ports" => json_decode($data['ports'], true));
			
			return $this->ports;
		
		}
	
	}

}
 
 ?>