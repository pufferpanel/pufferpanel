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
			$getCommandData = $core->framework->auth->generateSSH2Connection(array(
				'ip' => $node['sftp_ip'],
				'user' => $node['username']
			), array(
				'pub' => $node['ssh_pub'],
				'priv' => $node['ssh_priv'],
				'secret' => $node['ssh_secret'],
				'secret_iv' => $node['ssh_secret_iv']
			), true)->executeSSH2Command('sudo du -s '.$node['server_dir'].$core->framework->server->getData('ftp_user').'/server', true);
			
			if($getCommandData === false)
				exit('<div class="alert alert-danger">Unable to connect to the node.</div>');
			elseif(empty($getCommandData))
				exit('<div class="alert alert-danger">Unable to execute command on the server.</div>');
			
			
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
			
		echo '	<h5>CPU Usage</h5>
				<div class="progress">
				  	<div class="progress-bar" id="cpu_bar" style="width:'.round($core->framework->gsd->retrieve_process('cpu'), 0, PHP_ROUND_HALF_UP).'%;max-width:100%;">'.round($core->framework->gsd->retrieve_process('cpu'), 0, PHP_ROUND_HALF_UP).'%</div>
				</div>';
			
		echo '	<h5>Memory Usage</h5>
				<div class="progress">
				  	<div class="progress-bar" id="memory_bar" style="width:'.(($core->framework->files->format($core->framework->gsd->retrieve_process('memory')) / $core->framework->server->getData('max_ram')) * 100).'%;max-width:100%;">'.$core->framework->files->format($core->framework->gsd->retrieve_process('memory')).'MB / '.$core->framework->server->getData('max_ram').'MB</div>
				</div>';
			
			$onlinePlayers = null;
			$players = $core->framework->gsd->retrieve('players');
			$i = 0;
			
			if(count($players) > 0){
			
				foreach($players as $player){
			
					$onlinePlayers .= '<img data-toggle="tooltip" src="http://i.fishbans.com/helm/'.$player.'/32" title="'.$player.'" style="padding: 0 2px 6px 0;"/>';
					$i++;
			
				}
				
			}else{
			
				$onlinePlayers = '<p class="text-muted">No players are currently online.</p>';
			
			}
		
		echo '	<h5>Players Online</h5>
				<span id="player_list">'.$onlinePlayers.'</span>';
				
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