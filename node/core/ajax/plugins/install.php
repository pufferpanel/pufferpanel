<?php
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === true){

	if(isset($_POST['plugin'])){
	
		list($slug, $pluginID) = explode('|', $_POST['plugin']);
		
		$pluginName = str_replace(array(' ', '+', '%20'), '', $slug);
		$data = file_get_contents('http://api.bukget.org/3/plugins/bukkit/'.$slug);
		$data = json_decode($data, true);

		$filename = $data['versions'][$pluginID]['filename'];
		$downloadPath = $data['versions'][$pluginID]['download'];

			$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node_name` = ? LIMIT 1");
			$nodeSQLConnect->execute(array($core->framework->server->getData('node')));
			$node = $nodeSQLConnect->fetch();
			
			/*
			 * Connect and Run Function
			 */
			$con = ssh2_connect($node['node_ip'], 22);
			ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])));
			
			$stream = ssh2_exec($con, 'cd /srv/scripts && ./install_plugin.sh '.$core->framework->server->getData('ftp_user').' "'.$downloadPath.'" "'.$core->framework->server->getData('path').'plugins" "'.$data['plugin_name'].'" "'.$filename.'"');
			$errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
			
			stream_set_blocking($errorStream, true);
			stream_set_blocking($stream, true);
			
			echo "Output: " . stream_get_contents($stream);
			
			fclose($errorStream);
			fclose($stream);
			
	}

}
?>