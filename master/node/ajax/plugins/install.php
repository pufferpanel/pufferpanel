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
		
		$context = stream_context_create(array(
			"http" => array(
				"method" => "GET",
				"timeout" => 5
			)
		));
		$data = @file_get_contents('http://api.bukget.org/3/plugins/bukkit/'.$slug, 0, $context);
		$data = json_decode($data, true);

		$filename = $data['versions'][$pluginID]['filename'];
		$downloadPath = $data['versions'][$pluginID]['download'];
			
			/*
			 * Connect and Run Function
			 */
			$callbackData = $core->framework->ssh->generateSSH2Connection(array(
				'ip' => $core->framework->server->nodeData('sftp_ip'),
				'user' => $core->framework->server->nodeData('username')
			), array(
				'pub' => $core->framework->server->nodeData('ssh_pub'),
				'priv' => $core->framework->server->nodeData('ssh_priv'),
				'secret' => $core->framework->server->nodeData('ssh_secret'),
				'secret_iv' => $core->framework->server->nodeData('ssh_secret_iv')
			))->executeSSH2Command('cd /srv/scripts; sudo ./install_plugin.sh '.$core->framework->server->getData('ftp_user').' "'.$downloadPath.'" "'.$core->framework->server->getData('path').'plugins" "'.$data['plugin_name'].'" "'.$filename.'"', true);
			
			if(!empty($callbackData))
				echo $callbackData;
			
            $core->framework->log->getUrl()->addLog(0, 1, array('user.install_plugin', 'A plugin was installed.'));
			
	}

}
?>