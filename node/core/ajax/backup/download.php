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

	if(isset($_GET['id']) && !empty($_GET['id']) && is_numeric($_GET['id'])){
	
		$query = $mysql->prepare("SELECT * FROM `backups` WHERE `id` = :id AND `server` = :hash LIMIT 1");
		$query->execute(array(
			':id' => $_GET['id'],
			':hash' => $core->framework->server->getData('hash')
		));
		
			if($query->rowCount() == 1){
			
				$row = $query->fetch();
				
				$nodeQuery = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = ? LIMIT 1");
				$nodeQuery->execute(array($core->framework->server->getData('node')));
				$node = $nodeQuery->fetch();
				
				header("Pragma: private");
				header("Expires: 0");
				header("Cache-Control: must-revalidate, post-check=0, pre-check=0");
				header("Content-Type: application/force-download");
				header("Content-Description: File Transfer");
				header('Content-Disposition: response; filename="'.$row['file_name'].'.tar.gz"');
				header("Content-Transfer-Encoding: binary");
				header('Accept-Ranges: bytes');
				header("Content-Length: ".filesize($node['backup_dir'].$core->framework->server->getData('name').'/'.$row['file_name'].'.tar.gz'));
					
				$core->framework->files->download($node['backup_dir'].$core->framework->server->getData('name').'/'.$row['file_name'].'.tar.gz');
				
			}else{
			
				die("error");
			
			}
	
	}else{
	
		die('error');
	
	}

}
?>