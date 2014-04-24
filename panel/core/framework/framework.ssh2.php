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

class ssh extends \Auth\auth {

	//use Database\database;
	
	private $ssh2_connection;
	private $connectFailed;
	
	public function __construct($settingValue) {
	
		$this->mysql = self::connect();
		$this->useSSHKeys = $settingValue;
		
	}
	
	public function generateSSH2Connection($id, $usekey = false, $return = false){
	
		/*
		 *Connect to Node or Connect to User?
		 */
		if($usekey === true){
							
			$query = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
			$query->execute(array(
				':id' => $id
			));
			
				if($query->rowCount() == 1)
					$this->node = $query->fetch();
				else
					exit('SQL Error encountered in SSH2 trying to select node (2).');
		
		}else{
							
			$query = $this->mysql->prepare("SELECT * FROM `servers` WHERE `id` = :id");
			$query->execute(array(
				':id' => $id
			));
			
				if($query->rowCount() == 1)
					$this->server = $query->fetch();
				else
					exit('SQL Error encountered in SSH2 trying to select server.');
					
			$query = $this->mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
			$query->execute(array(
				':id' => $this->server['node']
			));
			
				if($query->rowCount() == 1)
					$this->node = $query->fetch();
				else
					exit('SQL Error encountered in SSH2 trying to select node.');
			
		
		}
	
		/*
		 * Add check to ensure IP can be conencted to, otherwise this function runs like a dead turtle.
		 */
		$this->connectFailed = false;	
		if(!$fp = @fsockopen($this->node['sftp_ip'], 22, $errString, $errCode, 2))
			$this->connectFailed = true;
		
		/*
		 * Connect to Node
		 */
		if($this->connectFailed === false){
		
			$this->ssh2_connection = ssh2_connect($this->node['sftp_ip'], 22);	
		
				if($usekey === true && $this->useSSHKeys == 1){
				
					if(!ssh2_auth_pubkey_file(
						$this->ssh2_connection,
						$this->node['username'],
						$this->node['ssh_pub'],
						$this->node['ssh_priv'],
						$this->decrypt($this->node['ssh_secret'], $this->node['ssh_secret_iv'])
					))
						$this->connectFailed = true;
					else
						null;
						
				}else {
				
					if($usekey === true){
					
						if(!ssh2_auth_password(
							$this->ssh2_connection,
							$this->node['username'],
							$this->decrypt($this->node['ssh_secret'], $this->node['ssh_secret_iv'])
						))
							$this->connectFailed = true;
							
					}else{
					
						if(!ssh2_auth_password(
							$this->ssh2_connection,
							$this->server['ftp_user'],
							$this->decrypt($this->server['ftp_pass'], $this->server['encryption_iv'])
						))
							$this->connectFailed = true;
					
					}
					
				}
						
						
			if($return === false)
				return $this;
			else
				return $this->ssh2_connection;
						
		}else{
		
			return false;
		
		}
	
	}

	public function executeSSH2Command($command, $callback = false, $tty = true){
	
		if($this->connectFailed === true){
			error_log('[PufferPanel] Error in "framework.auth.php" --> unable to connect to server using ssh2_auth. Authentication failed.');
			return false;
		}else{
	
			$this->stream = ssh2_exec($this->ssh2_connection, $command, $tty);
			$this->errorStream = ssh2_fetch_stream($this->stream, SSH2_STREAM_STDERR);
			
			/*
			 * Handle Errors
			 */
			stream_set_blocking($this->errorStream, true);
			stream_set_blocking($this->stream, true);
			
			$this->isError = stream_get_contents($this->errorStream);
			
			/*
			 * Do we want the data?
			 */
			if($callback === true)
				$this->streamData = stream_get_contents($this->stream);
			
			/*
			 * Close the Streams
			 */
			fclose($this->errorStream);
			fclose($this->stream);
			
			/*
			 * Return Data (true = completed, false = error)
			 */
			if(!empty($this->isError))
				return false;
			else
				if($callback === false)
					return true;
				else
					return $this->streamData;
				
		}
	
	}

}

?>