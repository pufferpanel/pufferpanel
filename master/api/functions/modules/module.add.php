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

/* 	auth => key,
	function => add,
	data =>	(
		
	)
			
*/

trait addServer {

	use globalInit;
	
	public function validateRequest() {
	
		$data = $this->getStoredData();
		if(array_key_exists('function', $data) && $data['function'] == 'add'){
		
			/*
			 * Is all of the data here?
			 */
			$dataOptions = array(
				'server_name',
				'node',
				'modpack',
				'email',
				'server_ip',
				'server_port',
				'alloc_mem',
				'alloc_disk',
				'sftp_pass',
				'sftp_pass_2',
				'backup_disk',
				'backup_files',
				'cpu_limit'
			); 
			
			foreach($dataOptions as $dataOption) {
			
				if(!array_key_exists($dataOption, $data['data']) || $data['data'][$dataOption] == "")
					$this->throwResponse('Missing required data values in API call.', false);
			
			}
		
		}else{
		
			$this->throwResponse('Accessing API in an illegal manner.', false);
		
		}
	
	}
	
	public function runRequest() {
	
			
	
	}

}

class apiModuleAddServer extends GSD_Query {
	
	use addServer;
	
	public function __construct() {
	
		$this->validateRequest();
		$this->runRequest();
	
	}
	
	private function runSSH() {
	
		apiModuleAddServer_Extended::run();
	
	}

}

class apiModuleAddServer_Extended extends ssh {

	use addServer;
	
	public function run() {
	
	
	}

}

?>