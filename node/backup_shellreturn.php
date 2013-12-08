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
 * Include Dependencies
 */
require_once('core/framework/framework.core.php');


/*
 * Backup Process Done
 * Function Cannot run after Auth check
 */
if(isset($_GET['do']) && $_GET['do'] == 'backup_done' && isset($_GET['server']) && $_SERVER['REMOTE_ADDR'] == $_SERVER['SERVER_ADDR']){

	$selectServerData = $mysql->prepare("SELECT * FROM `servers` WHERE `hash` = ?");
	$selectServerData->execute(array($_GET['server']));
	
		if($selectServerData->rowCount() == 1){
		
			/*
			 * Send save-on & save-all and update backup status
			 */
			$serverData = $selectServerData->fetch();
			
			$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node` = ? LIMIT 1");
			$nodeSQLConnect->execute(array($serverData['node']));
			$node = $nodeSQLConnect->fetch();
		
			/*
			 * Send Command to Backup Stuff
			 */
			if($core->framework->rcon->online($serverData['server_ip'], $serverData['server_port']) === true){
				
                /*
                 * Connections
                 */
                $con = ssh2_connect($node['node_ip'], 22);
                ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])));
                
                $s = ssh2_exec($con, 'cd /srv/scripts; ./send_command.sh '.$core->framework->server->getData('ftp_user').' "save-on"');
                stream_set_blocking($s, true);
            
			
				$selectBackup = $mysql->prepare("SELECT * FROM `backups` WHERE `backup_token` = ?");
				$selectBackup->execute(array($_GET['token']));
				
				$row = $selectBackup->fetch();
					
				$fileSHA1 = sha1_file($node['backup_dir'].$serverData['name'].'/'.$row['file_name'].'.tar.gz');
				$fileMD5 = md5_file($node['backup_dir'].$serverData['name'].'/'.$row['file_name'].'.tar.gz');
					
					/*
					 * Update MySQL Stuff
					 */
					$updateBackups = $mysql->prepare("UPDATE `backups` SET `complete` = '1', `timeend` = :time, `md5` = :md5, `sha1` = :sha1 WHERE `server` = :server AND `complete` = 0 AND `backup_token` = :token");	
					$updateBackups->execute(array(
						':time' => time(),
						':md5' => md5_file($node['backup_dir'].$serverData['ftp_user'].'/'.$row['file_name'].'.tar.gz'), 
						':sha1' => sha1_file($node['backup_dir'].$serverData['ftp_user'].'/'.$row['file_name'].'.tar.gz'),
						':server' => $serverData['hash'],
						':token' => $_GET['token']
					));		
			
		}else{
		
			echo 'Not Found in MYSQL';
		
		}
		
	exit();

}