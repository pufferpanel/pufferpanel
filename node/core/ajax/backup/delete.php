<?php 
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