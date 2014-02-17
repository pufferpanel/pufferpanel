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
	
	if($_POST['command'] && $_POST['command'] == 'stats'){
	
		$maxSpace = $core->framework->server->getData('disk_space') * 1024 * 1024;
		
		/*
		* Get the Server Node Info
		*/
		$query = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :nodeid");
		$query->execute(array(
		':nodeid' => $core->framework->server->getData('node')
		));
		
		$node = $query->fetch();
		
			/*
			 * Run Command
			 */
			$con = ssh2_connect($node['sftp_ip'], 22);
			if(!ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv']))))
				exit('<div class="alert alert-danger">Unable to connect to the node.</div>');
			
			
			$stream = ssh2_exec($con, 'sudo du -s '.$node['server_dir'].$core->framework->server->getData('ftp_user').'/server', true);
			$errorStream = ssh2_fetch_stream($stream, SSH2_STREAM_STDERR);
			
			stream_set_blocking($errorStream, true);
			stream_set_blocking($stream, true);
			
			$getCommandData = stream_get_contents($stream);
			if(empty($getCommandData))
				exit('<div class="alert alert-danger">Unable to execute command on the server.</div>');
			
			fclose($errorStream);
			fclose($stream);
		
		/*
		 * Do Math
		 */
		$getCommandData = explode("\t", $getCommandData);
				
		$returnSpacePercent = round((($getCommandData[0] * 1024) / $maxSpace), 2) * 100;
		if($returnSpacePercent < 1){ $returnSpacePercent = 1; }
		
		$spaceUsedH = $core->framework->files->formatSize($getCommandData[0] * 1024);
		$maxSpaceH = $core->framework->files->formatSize($maxSpace);
		
		echo '	<div class="progress">
		  			<div class="progress-bar" style="width:'.$returnSpacePercent.'%"></div>
				</div>
				<p class="text-muted">You are using '.$spaceUsedH.' of your maximum '.$core->framework->server->getData('disk_space').' MB of disk space.</p>';
				
	}else if($_POST['command'] && $_POST['command'] == 'players'){
		
		/*
		 * Query Dodads
		 */
		if($core->framework->gsd->online() !== true){
			exit('<div class="alert alert-danger">The server appears to be offline.</div>');
		}
			
			$onlinePlayers = null;
			$players = $core->framework->query->getPlayers();
			$i = 0;
			
			if($players !== false){
			
				foreach($players as $player){
			
					$onlinePlayers .= '<img data-toggle="tooltip" src="http://i.fishbans.com/helm/'.$player.'/32" title="'.$player.'" style="padding: 0 2px 6px 0;"/>';
					$i++;
			
				}
				
			}else{
			
				$onlinePlayers = '<p class="text-muted">No players are currently online.</p>';
			
			}
		
		$playerPercentage = round($i / $core->framework->gsd->retrieve('maxplayers'), 2) * 100;
		$playerPercentage = ($playerPercentage < 1) ? 1 : $playerPercentage;
		
		echo '
				<div class="progress">
				  	<div class="progress-bar" style="width:'.$playerPercentage.'%"></div>
				</div>
				'.$onlinePlayers;
				
	}else if($_POST['command'] && $_POST['command'] == 'info'){
	
		/*
		 * Query Dodads
		 */
		
		
		if($core->framework->gsd->online() === true){
		
			$version = $core->framework->query->getInfo('Software');
			
			$plugins = null;

			if($pluginList = is_array($core->framework->query->getInfo('Plugins'))){

				foreach($pluginList as $id => $name){

					$plugins .= $name.', ';

				}
				$plugins = rtrim($plugins, ", ");

			}

			$plugins = (is_null($plugins)) ? "No plugins detected." : $plugins;
		
		}else{
		
			$version = "Unable to query the server.";
			$plugins = "Unable to query the server.";
		
		}
		
			echo '
					<table class="table table-striped table-bordered table-hover">
						<thead>
							<tr>
								<th>Information</th>
								<th>Data</th>
							</tr>
						</thead>
						<tbody>
							<tr>
								<td><strong>Connection</strong></td>
								<td>'.$core->framework->server->getData('server_ip').':'.$core->framework->server->getData('server_port').'</td>
							</tr>
							<tr>
								<td><strong>Node</strong></td>
								<td>'.$core->framework->settings->nodeName($core->framework->server->getData('node')).'</td>
							</tr>
							<tr>
								<td><strong>Version</strong></td>
								<td>'.$version.'</td>
							</tr>
							<tr>
								<td><strong>Plugins</strong></td>
								<td>'.$plugins.'</td>
							</tr>
						</tbody>
					</table>';
	
	}

}else{

	exit('Invalid Authentication.');

}
?>