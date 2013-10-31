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
 * PufferPanel API Functions File
*/

class api_terminate extends api_main {

	public function terminate($data)
	{
	
		$data = unserialize(base64_decode($data));
		$this->e = false;
		
		/*
		 * Run Terminate Process
		 */
		$this->getServerInfo = mysql_query("SELECT * FROM `servers` WHERE `whmcs_id` = '".mysql_real_escape_string($data['serviceid'])."'");
			
			if(mysql_num_rows($this->getServerInfo) == 1){
			
				$this->sInfo = mysql_fetch_assoc($this->getServerInfo);
				
				/*
				 * Send Shutdown Command
				 */
				parent::_serverShutdown($this->sInfo['node'], $this->sInfo['name']);
				sleep(1); #Rest Server for a Second
				
				/*
				 * Delete SFTP Account (Server)
				 */
				parent::_serverDelete($this->sInfo['node'], $this->sInfo['name'], 'sftp');
				 
				/*
				 * Delete All Files
				 */
				parent::_serverDelete($this->sInfo['node'], $this->sInfo['name'], 'files');
				
				/*
				 * Delete All Backups
				 */
				parent::_serverDelete($this->sInfo['node'], $this->sInfo['name'], 'backups');
				
				/*
				 * Update Ports
				 */
				parent::_WHMCSUpdateValues($this->sInfo['max_ram'], $this->sInfo['node'], $this->sInfo['server_ip'], $this->sInfo['server_port'], 'terminate');
				 
				/*
				 * Return Data
				 */
				return json_encode(array('success' => true));
			
			}else{
			
				return json_encode(array('success' => false, 'error' => 'Could not locate server in PufferPanel Database.'));
			
			}
	
	}
		
}
?>