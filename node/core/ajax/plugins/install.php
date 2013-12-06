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

	if(isset($_POST['plugin'])){
	
		list($slug, $pluginID) = explode('|', $_POST['plugin']);
		
		$pluginName = str_replace(array(' ', '+', '%20'), '', $slug);
		$data = file_get_contents('http://api.bukget.org/3/plugins/bukkit/'.$slug);
		$data = json_decode($data, true);

		$filename = $data['versions'][$pluginID]['filename'];
		$downloadPath = $data['versions'][$pluginID]['download'];

			$nodeSQLConnect = $mysql->prepare("SELECT * FROM `nodes` WHERE `node` = ? LIMIT 1");
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
			
            $core->framework->log->getUrl()->addLog(0, 1, array('user.install_plugin', 'A plugin was installed.'));
        
			echo "Output: " . stream_get_contents($stream);
			
			fclose($errorStream);
			fclose($stream);
			
	}

}
?>