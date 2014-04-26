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
 
namespace Modules;
 
trait validate {
	
	use \Modules\functions;
	
	/*
	 * Handle Incoming Queries
	 */
	public static function run() {
	
		if(!array_key_exists(self::getStoredData()['function'], array(
			'add', 'delete', 'info'
		))) 
			self::throwResponse("Accessing API in an illegal manner.", false);
		else
			self::validateData();
	
	}
	
	/* 
	 * Middle Man for Handing Functions
	 */
	private function validateData() {
	
		switch(self::getStoredData()['function']) {
		
			case 'add':
				self::validateAddServerRequest();
				break;
			case 'delete':
				self::validateDeleteServerRequest();
				break;
			case 'info':
				self::validateInformationRequest();
				break;
			default:
				self::throwResponse('Accessing API in an illegal manner.', false);
		
		}
	
	}
	
	/*
	 * Validate Request Data for Adding a Server
	 */
	private function validateAddServerRequest() {
	
		/*
		 * Is all of the data here?
		 */
		$dataOptions = array(
			'server_name', 'node', 'modpack', 'email', 'server_ip', 'server_port', 'alloc_mem', 'alloc_disk', 'sftp_pass', 'sftp_pass_2', 'cpu_limit'
		);
					
		foreach($dataOptions as $dataOption)
			if(!array_key_exists($dataOption, $data['data']))
				self::throwResponse('Missing required data values in API call.', false);
	
		/*
		 * Run Function
		 */
		$run = new apiModuleAddServer();
	
	}
	
}