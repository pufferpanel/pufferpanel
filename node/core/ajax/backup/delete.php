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

	if(isset($_POST['id']) && !empty($_POST['id']) && is_numeric($_POST['id'])){

		$query = $mysql->prepare("SELECT * FROM `backups` WHERE `id` = :id AND `server` = :hash LIMIT 1");
		$query->execute(array(
			':id' => $_POST['id'],
			':hash' => $core->framework->server->getData('hash')
		));
		
			if($query->rowCount() == 1){
			
				$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node_name` = ? LIMIT 1");
				$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
				
				$row = $query->fetch();
				$node = $nodeSQLConnect->fetch();
					
				if(file_exists($node['backup_dir'].$core->framework->server->getData('name').'/server/'.$row['file_name'].'.tar.gz')){
				
					$con = ssh2_connect($node['node_ip'], 22);
					ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])));
				
					$deleteFile = escapeshellarg($node['backup_dir'].$core->framework->server->getData('name').'/server/'.$row['file_name'].'.tar.gz');
				
					$s = ssh2_exec($con, 'echo "'.openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])).'" | sudo -S su - '.$node['username'].' -c "cd '.$row['backup_dir'].$core->framework->server->getData('name').'; rm -rf '.$deleteFile.'"');
					stream_set_blocking($s, true);
					
					$deleteBackup = $mysql->prepare("DELETE FROM `backups` WHERE `id` = :id AND `server` = :hash");
					$deleteBackup->execute(array(
						':id' => $_POST['id'],
						':hash' => $core->framework->server->getData('hash')
					));
					
					echo 'File has been deleted.';
				
				}else{
				
					die('error');
				
				}
			
			}else{
			
				die('error');
			
			}

	}else{
	
		die('error');
	
	}
	
}
?>