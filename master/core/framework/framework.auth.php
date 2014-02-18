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
 * Auth Framework Core
 */

class auth extends dbConn {
	
	private $ssh2_connection;
	private $connectFailed;
	
	public function __construct()
		{
		
			$this->mysql = parent::getConnection();
		
		}
	
	public function hash($raw){
	
		return password_hash($raw, PASSWORD_BCRYPT);
	
	}
	
	private function password_compare($raw, $hashed){
	
		if(password_verify($raw, $hashed))
			return true;
		else
			return false;
	
	}
	
	public function generate_iv(){
	
		return base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
		
	}
	
	public function encrypt($raw, $iv, $method = 'AES-256-CBC'){
	
		return openssl_encrypt($raw, $method, file_get_contents(HASH), false, base64_decode($iv));
	
	}
	
	public function decrypt($encrypted, $iv, $method = 'AES-256-CBC'){
	
		return openssl_decrypt($encrypted, $method, file_get_contents(HASH), 0, base64_decode($iv));
	
	}
	
	public function verifyPassword($email, $raw){
	
		$this->get = $this->mysql->prepare("SELECT `password` FROM `users` WHERE `email` = :email");
		$this->get->execute(array(
			':email' => $email
		));
	
			if($this->get->rowCount() == 1){
			
				$this->row = $this->get->fetch();
				return $this->password_compare($raw, $this->row['password']);
				
			}else
				return false;
	
	}
	
	public function generateSSH2Connection($vars, $pubkey, $usekey = false){
	
		/*
		 * Connect to Node
		 */
		$this->ssh2_connection = ssh2_connect($vars['ip'], 22);
		$this->connectFailed = false;		
		
			if($usekey === true)
				if(!ssh2_auth_pubkey_file($this->ssh2_connection, $vars['user'], $pubkey['pub'], $pubkey['priv'], $this->decrypt($pubkey['secret'], $pubkey['secret_iv'])))
					$this->connectFailed = true;
			else
				if(!ssh2_auth_password($this->ssh2_connection, $vars['user'], $this->decrypt($vars['pass'], $vars['iv'])))
					$this->connectFailed = true;
		
		return $this;
	
	}

	public function executeSSH2Command($command, $callback = false, $tty = true){
	
		if($this->connectFailed === true){
			return false;
		}else{
	
			$this->stream = ssh2_exec($this->ssh2_connection, $command, $tty);
			$this->errorStream = ssh2_fetch_stream($this->stream, SSH2_STREAM_STDERR);
			
			/*
			 * Handle Errors
			 */
			stream_set_blocking($this->errorStream, true);
			stream_set_blocking($this->stream, true);
			
			$this->isError = stream_get_contents($errorStream);
			
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
	
	public function keygen($amount){
		
		$keyset  = "abcdefghijklmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ0123456789";
		
		$randkey = null;
		
		for ($i=0; $i<$amount; $i++)
			$randkey .= substr($keyset, rand(0, strlen($keyset)-1), 1);
		
		return $randkey;
					
	}

	public function isLoggedIn($ip, $session, $serverhash = null, $acp = false){

		$this->query = $this->mysql->prepare("SELECT * FROM `users` WHERE `session_ip` = :sessip AND `session_id` = :session AND `session_expires` > :sesexp");
		$this->query->execute(array(
			':sessip' => $ip,
			':session' => $session,
			':sesexp' => time()
		));

			if($this->query->rowCount() == 1){

				$this->row = $this->query->fetch();

					/*
					 * Accessing Admin Directory
					 */
					if($this->row['root_admin'] != 1 && $acp === true)
						return false;
					else{

						/*
						 * Allow Admins Access to any Server
						 */
						if($this->row['root_admin'] != '1'){	
	
							/*
							 * Validate User is Owner of Server
							 */
							if(!is_null($serverhash)){
							
								$this->_validateServer = $this->mysql->prepare("SELECT * FROM `servers` WHERE `hash` = :shash AND `owner_id` = :ownerid AND `active` = 1");
								$this->_validateServer->execute(array(
									':shash' => $serverhash,
									':ownerid' => $this->row['id']
								));
		
									if($this->_validateServer->rowCount() == 1){
		
										$this->updateUsers = $this->mysql->prepare("UPDATE `users` SET `session_expires` = :sesexp WHERE `session_ip` = :sesip AND `session_id` = :sesid");
										$this->updateUsers->execute(array(
											':sesexp' => time() + 1800,
											':sesip' => $ip,
											':sesid' => $session
		
										));
		
											return true;
		
									}else{
		
										return false;
		
									}
							
							/*
							 * Just Check if they are Logged In
							 */	
							}else{
							
								return true;
								
							}
	
						}else{
	
							$this->updateUsers = $this->mysql->prepare("UPDATE `users` SET `session_expires` = :sesexp WHERE `session_ip` = :sesip AND `session_id` = :sesid");
							$this->updateUsers->execute(array(
								':sesexp' => time() + 1800,
								':sesip' => $ip,
								':sesid' => $session
	
							));
	
								return true;
	
						}
						
					}

			}else{

				return false;

			}

	}
	
	public function getCookie($cookie){
	
		if(isset($_COOKIE[$cookie])){
		
			return $_COOKIE[$cookie];
		
		}else{
		
			return null;
		
		}
	
	}

}

?>