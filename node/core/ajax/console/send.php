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
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === true){

	if($core->framework->rcon->online($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
	
		$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = ? LIMIT 1");
		$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
		
		$row = $nodeSQLConnect->fetch();
		/*
		 * Connect and Run Function
		 */		
		if(!($con = ssh2_connect($row['sftp_ip'], 22))){
		    
            $core->framework->log->getUrl()->addLog(4, 1, array('system.sftp_connect_fail', 'The server `'.$core->framework->server->getData('name').'` was unable to have a command sent due to a connection error.'));
		    die("Unable to establish connection to Server.");
		    
		} else {
		    
		    
		    if(!ssh2_auth_password($con, $row['username'], openssl_decrypt($row['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($row['encryption_iv'])))) {
		        
                $core->framework->log->getUrl()->addLog(4, 1, array('system.sftp_auth_fail', 'The server `'.$core->framework->server->getData('name').'` was unable to have a command sent due to an authentication error.'));
		        die("Unable to Authenticate with Server.");
		    
		    }else{

				$stream = ssh2_exec($con, 'cd /srv/scripts; sudo ./send_command.sh '.$core->framework->server->getData('ftp_user').' "'.escapeshellcmd($_POST['command']).'"', true);
				stream_set_blocking($stream, true);
				fclose($stream);
		
		    }
		    
		}
		
	}else{
	
        $core->framework->log->getUrl()->addLog(1, 1, array('system.server_off', 'The server `'.$core->framework->server->getData('name').'` was unable to recieve a command because it is not on.'));
		exit('Could not establish connection to server. Is it turned on?');
	
	}

}else{

	exit('Invalid Authentication.');

}
?>