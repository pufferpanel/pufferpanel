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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === true){
	
	if($_POST['command'] && $_POST['command'] == 'stats'){
		
		$maxSpace = $core->server->getData('disk_space') * 1024 * 1024;
			
			/*
			 * Run Command
			 */
			$getCommandData = $core->ssh->generateSSH2Connection(array(
				'ip' => $core->server->nodeData('sftp_ip'),
				'user' => $core->server->nodeData('username')
			), array(
				'pub' => $core->server->nodeData('ssh_pub'),
				'priv' => $core->server->nodeData('ssh_priv'),
				'secret' => $core->server->nodeData('ssh_secret'),
				'secret_iv' => $core->server->nodeData('ssh_secret_iv')
			))->executeSSH2Command('sudo du -s '.$core->server->nodeData('server_dir').$core->server->getData('ftp_user').'/server', true);
						
			if($getCommandData === false)
				exit('<div class="alert alert-danger">Unable to connect to the node.</div>');
			else if(empty($getCommandData))
				exit('<div class="alert alert-danger">Unable to execute command on the server.</div>');
		
			
		/*
		 * Do Math
		 */
		$getCommandData = explode("\t", $getCommandData);
		$returnSpacePercent = round((($getCommandData[0] * 1024) / $maxSpace), 2) * 100;
		if($returnSpacePercent < 1){ $returnSpacePercent = 1; }
		
		$spaceUsedH = $core->files->formatSize($getCommandData[0] * 1024);
		$maxSpaceH = $core->files->formatSize($maxSpace);
		
		echo '	<div class="progress">
		  			<div class="progress-bar" style="width:'.$returnSpacePercent.'%"></div>
				</div>
				<p class="text-muted">You are using '.$spaceUsedH.' of your maximum '.$core->server->getData('disk_space').' MB of disk space.</p>';
				
	}else if($_POST['command'] && $_POST['command'] == 'players'){
		
		/*
		 * Query Dodads
		 */
		if($core->gsd->online() !== true){
			exit('<div class="alert alert-danger">The server appears to be offline.</div>');
		}
		
		$cpu = round(($core->gsd->retrieve_process('cpu') / $core->server->getData('cpu_limit')) * 100, 2);
		$cpu = ($cpu > "100") ? "100" : $cpu;	
		echo '	<h5>CPU Usage</h5>
				<div class="progress">
				  	<div class="progress-bar" id="cpu_bar" style="width:'.$cpu.'%;max-width:100%;">'.round($cpu, 2).'%</div>
				</div>';
			
		echo '	<h5>Memory Usage</h5>
				<div class="progress">
				  	<div class="progress-bar" id="memory_bar" style="width:'.(($core->files->format($core->gsd->retrieve_process('memory')) / $core->server->getData('max_ram')) * 100).'%;max-width:100%;">'.$core->files->format($core->gsd->retrieve_process('memory')).'MB / '.$core->server->getData('max_ram').'MB</div>
				</div>';
			
			$onlinePlayers = null;
			$players = $core->gsd->retrieve('players');
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
								<td>'.$core->server->getData('server_ip').':'.$core->server->getData('server_port').'</td>
							</tr>
							<tr>
								<td><strong>Node</strong></td>
								<td>'.$core->settings->nodeName($core->server->getData('node')).'</td>
							</tr>
							<tr>
								<td><strong>Memory Allocated</strong></td>
								<td>'.$core->server->getData('max_ram').' MB</td>
							</tr>
							<tr>
								<td><strong>Disk Allocated</strong></td>
								<td>'.$core->server->getData('disk_space').' MB</td>
							</tr>
						</tbody>
					</table>';
	
	}

}else{

	exit('Invalid Authentication.');

}
?>