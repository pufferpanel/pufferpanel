<?php

/*
 * PufferPanel API
 * WHMCS Module Provisioning
 */
 
class api_provision extends api_main {

	public function provision($data){
	
		$data = unserialize(base64_decode($data));
		$this->e = false;
		
		/*
		 * All Variables Accounted For
		 */
		if(count($data) == 11){
			
			$time = time(); 
			
			/*
			 * Account Data
			 */
			$email = $data['email'];
			$name = $data['firstname'].".".$data['lastname'];
			$name = (strlen($name) > 25) ? substr($name, 0, 25) : $name;
			$name = strtolower($name);
			$whmcsuser = $data['userid'];
			$raw_password = parent::_keygen('16');
			$password = parent::_encrypt($raw_password);
			
				/*
				 * No Such User; Add to Database; Otherwise, skip this step
				 */
				if(parent::_userExists($email) === false){
				
					mysql_query("INSERT INTO `users` VALUES(NULL, '".mysql_real_escape_string($data['userid'])."', '".$name."', '".$email."', '".$password."', '".$time."', 'owner', NULL, NULL, NULL, 0, 1, 1)")or $this->e = true;
									
				}
				
				/*
				 * No Server in Database; Add to Database; Otherwise Skip Step
				 */
				if(parent::_serverExists($data['serviceid']) === false){
				
					$ram = str_replace('MB', '', $data['ram']);
					$node = parent::_selectNode($ram);
					if($node != 'error_no_node'){
				
						$disk = (str_replace('GB', '', $data['disk']))*1024;
						$bacupFiles = $data['backup_limit'];
						$backupSpace = (str_replace('GB', '', $data['backup_disk_limit']))*1024;
						$userid = parent::_getUserId($email);
						$ip = parent::_selectIP($node);
						$port = parent::_selectPort($node, $ip);
						$server = array("name" => parent::_generateServerName($data['domain']), "node" => $node, "ip" => $ip, "port" => $port);
						$sftpHost = parent::_selectNodeSFTP($node); #This gives us the nodes SFTP IP
						$sftpPassword = parent::_keygen('16');
						$backupFiles = $data['backup_limit'];
						$backupSpace = (str_replace('GB', '', $data['backup_disk_limit']))*1024;
					
						mysql_query("INSERT INTO `servers` VALUES(NULL, '".mysql_real_escape_string($data['serviceid'])."', '".parent::_keygen('42')."', '".$server['node']."', '".$server['name']."', 1, '".$userid."', '".$name."', '".$ram."', '".$disk."', '/srv/servers/".$server['name']."/server/', 'shared', '".mysql_real_escape_string($data['identifier'])."', '".$time."', '".$server['ip']."', '".$port."', NULL, NULL, NULL, '".$sftpHost."', '".$server['name']."', '".$sftpPassword."', '".$backupFiles."', '".$backupSpace."')")or $this->e = true;
						
						
							if($this->e === false){
							
								/*
								 * Server Added - Update WHMCS Data Values
								 */
								$doProcess = parent::_process($node, $server['name'], $sftpPassword, $backupSpace);
								if($doProcess === true){
								
									parent::_WHMCSUpdateValues($ram, $node, $ip, $port, 'provision');
									$customArray = array("1" => $ip, "2" => $port, "3" => $server['name'], "4" => $raw_password, "5" => $sftpHost, "6" => $sftpPassword, "7" => $node);
									$customData = base64_encode(serialize($customArray));
									return json_encode(array('success' => true, 'custom' => $customData));
								
								}else{
								
									#parent::_WHMCSUpdateValues($ram, $node, $ip, $port, 'terminate');
									return json_encode(array('success' => false, 'error' => 'Unable to create account on physical node. (Error Returned: '.$doProcess.')'));
								
								}
							
							}
						
					}else{
					
						/*
						 * No Node was able to be provisioned. Email Admins.
						 */
						return json_encode(array('success' => false, 'error' => 'Unable to provision a node for the server. Setup failed!'));
						
					}
				
				}else{
				
					return json_encode(array('success' => false, 'error' => 'Server already exists in server database.'));
				
				}
		
		}else{
		
			return json_encode(array('success' => false, 'error' => 'Incorrect number of data variables passed to the API.'));
		
		}
	
	}

}

?>