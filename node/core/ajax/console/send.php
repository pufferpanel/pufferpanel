<?php
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === true){

	require_once('../../../core/framework/rcon/query.class.php');
	require_once('../../../core/framework/rcon/rcon.class.php');
	require_once('../../../core/framework/rcon/query.status.php');
	
	$rcon->status = new MinecraftStatus($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'));
	$rcon->query = new MinecraftQuery();
	$rcon->command = new MinecraftRcon();

	if($rcon->s->isOnline($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
	
		$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node_name` = ? LIMIT 1");
		$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
		
		$row = $nodeSQLConnect->fetch();
		/*
		 * Connect and Run Function
		 */		
		if(!($con = ssh2_connect($row['node_ip'], 22))){
		    
		    die("Unable to establish connection to Server.");
		    
		} else {
		    
		    
		    if(!ssh2_auth_password($con, $row['username'], openssl_decrypt($row['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($row['encryption_iv'])))) {
		        
		        die("Unable to Authenticate with Server.");
		    
		    }else{

				$stream = ssh2_exec($con, 'cd /srv/scripts; ./send_command.sh '.$core->framework->server->getData('name').' \"'.$_POST['command'].'\"');
				stream_set_blocking($stream, true);
				fclose($stream);
		
		    }
		    
		}
		
	}else{
	
		exit('Could not establish connection to server. Is it turned on?');
	
	}

}else{

	exit('Invalid Authentication.');

}
?>