<?php

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