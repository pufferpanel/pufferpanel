<?php

/*
 * PufferPanel API Functions File
 * Suspend a Server
 * Variables Reqired:
 *	- server_id
*/

class api_suspend extends api_main {

	public function suspend($data)
	{
	
		$data = unserialize(base64_decode($data));
		$this->e = false;
		
		/*
		 * Run Suspend Process
		 */
		$this->getServerInfo = mysql_query("SELECT * FROM `servers` WHERE `whmcs_id` = '".mysql_real_escape_string($data['serviceid'])."'");
			
			if(mysql_num_rows($this->getServerInfo) == 1){
			
				$this->sInfo = mysql_fetch_assoc($this->getServerInfo);
				$this->newPassword = parent::_keygen('18');
				
				/*
				 * Suspend Server from Panel
				 */
				mysql_query("INSERT INTO `whmcs_suspend_data` VALUES(NULL, '".$this->sInfo['id']."', '".mysql_real_escape_string($data['serviceid'])."', '".$this->sInfo['ftp_pass']."', 0)");
				
				mysql_query("UPDATE `servers` SET `active` = 0, `ftp_pass` = '".$this->newPassword."' WHERE `whmcs_id` = '".mysql_real_escape_string($data['serviceid'])."' AND `hash` = '".$this->sInfo['hash']."'");
			
				/*
				 * Send Shutdown Command
				 */
				parent::_serverShutdown($this->sInfo['node'], $this->sInfo['name']);
				sleep(1); #Rest Server for a Second
				 
				/*
				 * Suspend SFTP Account (Server)
				 */
				parent::_serverSuspend($this->sInfo['node'], $this->sInfo['name'], $this->newPassword);
				 
				/*
				 * Update Ports
				 */
				parent::_WHMCSUpdateValues($this->sInfo['max_ram'], $this->sInfo['node'], $this->sInfo['server_ip'], $this->sInfo['server_port'], 'suspend');
				 
				/*
				 * Return Data
				 */
				return json_encode(array('success' => true));
			
			}else{
			
				return json_encode(array('success' => false, 'error' => 'Could not locate server in PufferPanel Database.'));
			
			}
	
	}
	
	public function unsuspend($data)
	{
	
		$data = unserialize(base64_decode($data));
		$this->e = false;
	
		/*
		 * Run UnSuspend Process
		 */
		$this->getServerInfo = mysql_query("SELECT * FROM `servers` WHERE `whmcs_id` = '".mysql_real_escape_string($data['serviceid'])."'");
			
			if(mysql_num_rows($this->getServerInfo) == 1){
			
				$this->sInfo = mysql_fetch_assoc($this->getServerInfo);
				
				/*
				 * Get Old FTP Password
				 */
				$this->oldFTPQuery = mysql_query("SELECT * FROM `whmcs_suspend_data` WHERE `server_id` = '".$this->sInfo['id']."' AND `whmcs_server_id` = '".mysql_real_escape_string($data['serviceid'])."' AND `unsuspended` = '0'");
				
				$this->oldFTP = mysql_fetch_assoc($this->oldFTPQuery);
				$this->newPassword = $this->oldFTP['old_password'];
				
				/*
				 * Unsuspend Server from Panel
				 */
				mysql_query("UPDATE `servers` SET `active` = 1, `ftp_pass` = '".$this->newPassword."' WHERE `whmcs_id` = '".mysql_real_escape_string($data['serviceid'])."' AND `hash` = '".$this->sInfo['hash']."'");
							 
				/*
				 * Unsuspend SFTP Account (Server)
				 */
				parent::_serverSuspend($this->sInfo['node'], $this->sInfo['name'], $this->oldFTP['old_password']);
				
				/*
				 * Run Query
				 */
				mysql_query("UPDATE `whmcs_suspend_data` SET `unsuspended` = '1' WHERE `server_id` = '".$this->sInfo['id']."' AND `whmcs_server_id` = '".mysql_real_escape_string($data['serviceid'])."' AND `unsuspended` = '0'");
				 
				/*
				 * Update Ports
				 */
				parent::_WHMCSUpdateValues($this->sInfo['max_ram'], $this->sInfo['node'], $this->sInfo['server_ip'], $this->sInfo['server_port'], 'unsuspend');
				 
				/*
				 * Return Data
				 */
				return json_encode(array('success' => true));

			
			}else{
			
				return json_encode(array('success' => false, 'error' => 'Could not locate server/old FTP info in PufferPanel Database.'));
			
			}
	
	}

}

?>