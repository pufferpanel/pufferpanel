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
 * PufferPanel API
 * WHMCS Module Provisioning
 * Bad, pls no use
 */
 
class api_main {
 
	public function __construct(){
	
	
	}
	
	public function _encrypt($password){
		
		$salt = 'c2$VA9*^xj#0sA';
		return md5(md5($password).$salt);
	
	}
	
	public function _keygen($amount){
		
		$keyset  = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
		
		$randkey = "";
		
		for ($i=0; $i<$amount; $i++)
			$randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
		
		return $randkey;
			
	}
	
	public function _userExists($email){
	
		$query = mysql_query("SELECT id FROM `users` WHERE `email` = '".mysql_real_escape_string($email)."'");
		
			if(mysql_num_rows($query) > 0){
			
				return true;
				
			}else{
			
				return false;
				
			}
	
	}
	
	public function _serverExists($serverid){
	
		$query = mysql_query("SELECT id FROM `servers` WHERE `whmcs_id` = '".mysql_real_escape_string($serverid)."'");
		
			if(mysql_num_rows($query) > 0){
			
				return true;
				
			}else{
			
				return false;
				
			}
	
	}
	
	public function _getUserId($email, $id = 'id'){
	
		$query = mysql_query("SELECT id, whmcs_id FROM `users` WHERE `email` = '".mysql_real_escape_string($email)."' LIMIT 1");
		$row = mysql_fetch_assoc($query);
		
			if($id == 'id'){
				
				return $row['id'];
				
			}else{
			
				return $row['whmcs_id'];
			
			}
	
	}
	
	public function _generateServerName($domain){
	
		/*
		 * Set the Name
		 */
		if(strpos($domain, ".") === false){
		
			/*
			 * Invalid Domain. Shorten Name
			 */
			$this->_serverName = (strlen($domain) > 20) ? substr($domain, 0, 20) : $domain;
		
		}else{
		
			$domain = explode(".", $domain);
			$domain = $domain[0];
			
			$this->_serverName = (strlen($domain) > 20) ? substr($domain, 0, 20) : $domain;
		
		}
		
		/*
		 * Verify Name doesn't Exist Already
		 */
		$query = mysql_query("SELECT * FROM `servers` WHERE `name` = '".mysql_real_escape_string($this->_serverName)."'");
		
			if(mysql_num_rows($query) > 0){
			
				/*
				 * That Name Exisits, Make it more unique.
				 */
				$identifier = md5(time().$this->_serverName.mt_rand());
				$this->_serverName = substr($this->_serverName, 0, 10).'_'.substr($identifier, 0, 9);
			
			}
		
		/*
		 * Illegal Name Array
		 */
		$this->_illegalNames = array("root", "bin", "bash", "ssh", "ftp", "apache", "mysql", "nobody", "centos", "sftp", "debian", "backup", "www-data", "daemon", "danee", "dane", "lockwood", "klockwood", "puffrfish", "kylel", "kyle", "proxy", "mail", "http", "https", "news");
		
		/*
		 * Check for Illegal Server Name
		 */
		if(!in_array(strtolower($this->_serverName), $this->_illegalNames)){
		
			$this->_serverName = $this->_serverName;
			
		}else{
		
			$identifier = md5(time().$this->_serverName.mt_rand());
			$this->_serverName = substr($this->_serverName, 0, 10).'_'.substr($identifier, 0, 9);
		
		}
		
		/*
		 * Return Server Name
		 */
		return $this->_serverName;
	
	}
	
	public function _selectNodeSFTP($node){
	
		/*
		 * Get Values
		 */
		$this->getData = mysql_query("SELECT * FROM `whmcs_data` WHERE `node` = '".$node."' LIMIT 1");
		$this->row = mysql_fetch_assoc($this->getData);
		
		return $this->row['sftp_ip'];
	
	}
	
	public function _selectNode($ram){
	
		/*
		 * No node Defined; Select Random Node
		 */
		$query = mysql_query("SELECT * FROM `whmcs_data` WHERE `disabled` = 0 AND `full` = 0 AND `curr_alloc_ram` + ".$ram." <= `max_alloc_ram`");
		
			if(mysql_num_rows($query) > 0){
			
				/*
				 * Build Array with Data
				 */
				$this->_array = array();
				while($row = mysql_fetch_assoc($query)){
				
					/*
					 * Did this in Query, Check Again Regardless
					 */
					if(($ram + $row['curr_alloc_ram']) <= $row['max_alloc_ram']){
					
						$this->_array = array_merge($this->_array, array(array('node' => $row['node'], 'space_remaining' => ($row['max_alloc_ram'] - $row['curr_alloc_ram']))));
						
					}
				
				}
				
				
				/*
				 * Select Node with lowest space remaining
				 */
				$this->_min = $this->_array[0]['space_remaining'];
				$this->_key = 0;
				foreach($this->_array as $key => $val){
				
					if($val['space_remaining'] < $this->_min){
						$this->_min = $val['space_remaining'];
						$this->_key = $key;
					}
				
				}
				
				$this->node = $this->_array[$this->_key]['node'];
				
				/*
				 * Return the Node
				 */
				return $this->node;
			
			}else{
			
				return 'error_no_node';
			
			}
	
	}
	
	public function _selectIP($node){
	
		/*
		 * Select the node and JSON IP Structure
		 */
		$query = mysql_query("SELECT * FROM `whmcs_data` WHERE `node` = '".$node."' LIMIT 1");
		$row = mysql_fetch_assoc($query);
		
		/*
		 * Make JSON into Array and Select Random IP
		 */
		$ips = json_decode($row['ips'], true);
		
			/*
			 * Select IP with least free space in order to maximize.
			 */
			$useIP = array_keys($ips['options']);
			$useIP = $useIP[0];
			$currUsedPorts = $ips['max_ports_per_ip'];
			foreach($ips['options'] as $ip => $ports){
			
				if($ports['ports_free'] > 0){
				
					if($ports['ports_free'] < $currUsedPorts && $ports['ports_used'] < $ips['max_ports_per_ip']){
					
						$currUsedPorts = $ports['ports_free'];
						$useIP = $ip;
					
					}
				
				}
			
			}
			
		/*
		 * Return the IP
		 */
		return $useIP;
	
	}
	
	public function _selectPort($node, $ip){
	
		/*
		 * Select the Node and Ports JSON
		 */
		$query = mysql_query("SELECT * FROM `whmcs_data` WHERE `node` = '".$node."' LIMIT 1");
		$row = mysql_fetch_assoc($query);
		
		/*
		 * Filter out the Ports for Specific IP
		 */
		$ports = json_decode($row['ports'], true);
		
			$this->_stop = false;
			foreach($ports[$ip] as $port => $free){
			
				/*
				 * Get First Avaliable Port
				 */
				if($this->_stop === false){
				
					if($free == 0){
					
						$this->_usePort = $port;
						$this->_stop = true;
					
					}
				
				}
			
			}
			
		/*
		 * Return the Port
		 */
		return $this->_usePort;
	
	}
	
	public function _WHMCSUpdateValues($ram, $node, $ip, $port, $mode){
	
		/*
		 * Get Values
		 */
		$this->getData = mysql_query("SELECT * FROM `whmcs_data` WHERE `node` = '".$node."' LIMIT 1");
		$this->row = mysql_fetch_assoc($this->getData);
		 
		/*
		 * Processes
		 */ 
		if($mode == 'provision'){
		
			/*
			 * JSON Editing
			 */
			$this->_oldIPInfo = json_decode($this->row['ips'], true);
			$this->_oldIPInfo['options'][$ip]['ports_free'] = ($this->_oldIPInfo['options'][$ip]['ports_free'] - 1);
			$this->_oldIPInfo['options'][$ip]['ports_used'] = ($this->_oldIPInfo['options'][$ip]['ports_used'] + 1);
			
			$this->_oldPortInfo = json_decode($this->row['ports'], true);
			$this->_oldPortInfo[$ip][$port] = 1;
		
			/*
			 * Remove from Stocks
			 */
			$this->_newRamAlloc = ($this->row['curr_alloc_ram'] + $ram);
			$this->_newIPInfo = json_encode($this->_oldIPInfo);
			$this->_newPortInfo = json_encode($this->_oldPortInfo);
			
			/*
			 * Verify that server still has space after this
			 */
			if($this->_newRamAlloc < $this->row['max_alloc_ram']){
			
				$this->_isFull = 0;
			
			}else{
			
				$this->_isFull = 1;
			
			}
			
			/*
			 * Update Database
			 */
			mysql_query("UPDATE `whmcs_data` SET `curr_alloc_ram` = '".$this->_newRamAlloc."', `ips` = '".$this->_newIPInfo."', `ports` = '".$this->_newPortInfo."', `curr_servers` = curr_servers + 1, `full` = '".$this->_isFull."' WHERE `node` = '".$node."'");
			
		
		}else if($mode == 'suspend'){
		
			/*
			 * JSON Editing
			 */
			$this->_oldIPInfo = json_decode($this->row['ips'], true);
			$this->_oldIPInfo['options'][$ip]['ports_suspended'] = ($this->_oldIPInfo['options'][$ip]['ports_suspended'] + 1);
			
			/*
			 * Remove from Stocks
			 */
			$this->_newIPInfo = json_encode($this->_oldIPInfo);
			
			/*
			 * Update Database
			 */
			mysql_query("UPDATE `whmcs_data` SET `ips` = '".$this->_newIPInfo."' WHERE `node` = '".$node."'");
			
		}else if($mode == 'unsuspend'){
		
			/*
			 * JSON Editing
			 */
			$this->_oldIPInfo = json_decode($this->row['ips'], true);
			$this->_oldIPInfo['options'][$ip]['ports_suspended'] = ($this->_oldIPInfo['options'][$ip]['ports_suspended'] - 1);
			
			/*
			 * Remove from Stocks
			 */
			$this->_newIPInfo = json_encode($this->_oldIPInfo);
			
			/*
			 * Update Database
			 */
			mysql_query("UPDATE `whmcs_data` SET `ips` = '".$this->_newIPInfo."' WHERE `node` = '".$node."'");
		
		}else if($mode == 'terminate'){
		
			/*
			 * JSON Editing
			 */
			$this->_oldIPInfo = json_decode($this->row['ips'], true);
			$this->_oldIPInfo['options'][$ip]['ports_free'] = ($this->_oldIPInfo['options'][$ip]['ports_free'] + 1);
			$this->_oldIPInfo['options'][$ip]['ports_used'] = ($this->_oldIPInfo['options'][$ip]['ports_used'] - 1);
			
			$this->_oldPortInfo = json_decode($this->row['ports'], true);
			$this->_oldPortInfo[$ip][$port] = 0;
		
			/*
			 * Remove from Stocks
			 */
			$this->_newRamAlloc = ($this->row['curr_alloc_ram'] - $ram);
			$this->_newIPInfo = json_encode($this->_oldIPInfo);
			$this->_newPortInfo = json_encode($this->_oldPortInfo);
			
			/*
			 * Verify that server still has space after this
			 */
			if($this->_newRamAlloc < $this->row['max_alloc_ram']){
			
				$this->_isFull = 0;
			
			}else{
			
				$this->_isFull = 1;
			
			}
			
			/*
			 * Update Database
			 */
			mysql_query("UPDATE `whmcs_data` SET `curr_alloc_ram` = '".$this->_newRamAlloc."', `ips` = '".$this->_newIPInfo."', `ports` = '".$this->_newPortInfo."', `cur_servers` = cur_servers - 1, `full` = '".$this->_isFull."' WHERE `node` = '".$node."'");
		
		}
	
	}
	
	public function _process($node, $serverName, $usrSFTPPassword, $maxDisk){
	
		/*
		 * MySQL Connection
		 */
		$this->_nodeSQLConnect = mysql_query("SELECT * FROM `nodes` WHERE `node` = '".$node."' LIMIT 1");
		
			/*
			 * Verify
			 */
			if(mysql_num_rows($this->_nodeSQLConnect) == 1){
		
				$this->_nDataRow = mysql_fetch_assoc($this->_nodeSQLConnect);
		
				/*
				 * Connect to SFTP and Execute Commands
				 */
				if (!function_exists("ssh2_connect")) return 'No such function "ssh2_connect".';
				
				if(!($con = ssh2_connect($this->_nDataRow['ip'], 22))){
				    
				    return 'Cnnection failed: host.';
				    
				} else {
				    
				    /*
				     * NOTICE: These variables are specific to each server.
				     */
				    if(!ssh2_auth_password($con, $this->_nDataRow['user'], $this->_nDataRow['password'])) {
				        
				        return 'Connection failed: password.';
				    
				    }else{
				    
				    	/*
				    	 * Create User Stuff
				    	 */
				    	$softLimit = $maxDisk*1000;
				    	$hardLimit = ($maxDisk+1024)*1000;
				    									    									    									    	
				    	$s = ssh2_exec($con, 'echo "'.$this->_nDataRow['password'].'" | sudo -S su - root -c "cd /srv/scripts; ./create_user.sh '.$serverName.' '.$usrSFTPPassword.' '.$softLimit.' '.$hardLimit.'"');
				    	
				    	return true;
		
				    }
				    
				}
				
			}else{
			
				return 'MySQL Unable to locate Node.';
			
			}
	
	}
	
	public function _serverShutdown($node, $account){
	
		/*
		 * MySQL Connection
		 */
		$this->_nodeSQLConnect = mysql_query("SELECT * FROM `nodes` WHERE `node` = '".$node."' LIMIT 1");
		
			/*
			 * Verify
			 */
			if(mysql_num_rows($this->_nodeSQLConnect) == 1){
		
				$this->_nDataRow = mysql_fetch_assoc($this->_nodeSQLConnect);
		
				/*
				 * Connect to SFTP and Execute Commands
				 */
				if (!function_exists("ssh2_connect")) return 'No such function "ssh2_connect".';
				
				if(!($con = ssh2_connect($this->_nDataRow['ip'], 22))){
				    
				    return 'Cnnection failed: host.';
				    
				} else {
				    
				    /*
				     * NOTICE: These variables are specific to each server.
				     */
				    if(!ssh2_auth_password($con, $this->_nDataRow['user'], $this->_nDataRow['password'])) {
				        
				        return 'Connection failed: password.';
				    
				    }else{
				    				    									    									    									    	
				    	ssh2_exec($con, 'echo "'.$this->_nDataRow['password'].'" | sudo -S su - root -c "cd /srv/scripts; ./send_command.sh '.$account.' \"stop\""');
				    	
				    	return true;
		
				    }
				    
				}
				
			}else{
			
				return 'MySQL Unable to locate Node.';
			
			}
	
	}
	
	public function _serverSuspend($node, $account, $newPassword){
	
		/*
		 * MySQL Connection
		 */
		$this->_nodeSQLConnect = mysql_query("SELECT * FROM `nodes` WHERE `node` = '".$node."' LIMIT 1");
		
			/*
			 * Verify
			 */
			if(mysql_num_rows($this->_nodeSQLConnect) == 1){
		
				$this->_nDataRow = mysql_fetch_assoc($this->_nodeSQLConnect);
		
				/*
				 * Connect to SFTP and Execute Commands
				 */
				if (!function_exists("ssh2_connect")) return 'No such function "ssh2_connect".';
				
				if(!($con = ssh2_connect($this->_nDataRow['ip'], 22))){
				    
				    return 'Cnnection failed: host.';
				    
				} else {
				    
				    /*
				     * NOTICE: These variables are specific to each server.
				     */
				    if(!ssh2_auth_password($con, $this->_nDataRow['user'], $this->_nDataRow['password'])) {
				        
				        return 'Connection failed: password.';
				    
				    }else{
				    				    									    									    									    	
				    	ssh2_exec($con, 'echo "'.$this->_nDataRow['password'].'" | sudo -S su - root -c "cd /srv/scripts; ./suspend_server.sh '.$account.' '.$newPassword.'"');
				    	
				    	return true;
		
				    }
				    
				}
				
			}else{
			
				return 'MySQL Unable to locate Node.';
			
			}
	
	}
	
	public function _serverDelete($node, $account, $action){
	
		/*
		 * MySQL Connection
		 */
		$this->_nodeSQLConnect = mysql_query("SELECT * FROM `nodes` WHERE `node` = '".$node."' LIMIT 1");
		
			/*
			 * Verify
			 */
			if(mysql_num_rows($this->_nodeSQLConnect) == 1){
		
				$this->_nDataRow = mysql_fetch_assoc($this->_nodeSQLConnect);
		
				/*
				 * Connect to SFTP and Execute Commands
				 */
				if (!function_exists("ssh2_connect")) return 'No such function "ssh2_connect".';
				
				if(!($con = ssh2_connect($this->_nDataRow['ip'], 22))){
				    
				    return 'Cnnection failed: host.';
				    
				} else {
				    
				    /*
				     * NOTICE: These variables are specific to each server.
				     */
				    if(!ssh2_auth_password($con, $this->_nDataRow['user'], $this->_nDataRow['password'])) {
				        
				        return 'Connection failed: password.';
				    
				    }else{
				    				    									    									    									    	
				    	if($action == 'sftp'){
				    		
				    		ssh2_exec($con, 'echo "'.$this->_nDataRow['password'].'" | sudo -S su - root -c "cd /srv/scripts; ./terminate_sftp.sh '.$account.'"');
				    		return true;
				    		
				    	}else if($action == 'files'){
				    	
				    		ssh2_exec($con, 'echo "'.$this->_nDataRow['password'].'" | sudo -S su - root -c "cd /srv/scripts; ./terminate_files.sh '.$this->_nDataRow['data_dir'].' '.$account.'"');
				    		return true;
				    	
				    	}else if($action == 'backups'){
				    	
				    		ssh2_exec($con, 'echo "'.$this->_nDataRow['password'].'" | sudo -S su - root -c "cd /srv/scripts; ./terminate_backups.sh '.$this->_nDataRow['backup_dir'].' '.$account.'"');
				    		return true;
				    	
				    	}
		
				    }
				    
				}
				
			}else{
			
				return 'MySQL Unable to locate Node.';
			
			}
	
	}
 
}
 
?>