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

class ssh extends auth {

	private $ssh2_connection;
	private $connectFailed;
	
	public function generateSSH2Connection($vars, $pubkey = null, $return = false){
	
		/*
		 * Connect to Node
		 */
		$this->ssh2_connection = ssh2_connect($vars['ip'], 22);
		$this->connectFailed = false;		
		
			if(!empty($pubkey))
				if(!ssh2_auth_pubkey_file($this->ssh2_connection, $vars['user'], $pubkey['pub'], $pubkey['priv'], $this->decrypt($pubkey['secret'], $pubkey['secret_iv'])))
					$this->connectFailed = true;
				else
					null;
			else
				if(!ssh2_auth_password($this->ssh2_connection, $vars['user'], $this->decrypt($vars['pass'], $vars['iv'])))
					$this->connectFailed = true;
		
		if($return === false)
			return $this;
		else
			return $this->ssh2_connection;
	
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