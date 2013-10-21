<?php 
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
				
				$nodeQuery = $mysql->prepare("SELECT * FROM `nodes` WHERE `node_name` = ? LIMIT 1");
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