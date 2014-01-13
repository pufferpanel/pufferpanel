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
		$spaceUsed = $core->framework->files->getFolderSize($core->framework->server->getData('path'));
		    
		$returnSpacePercent = round(($spaceUsed / $maxSpace), 2) * 100;
		if($returnSpacePercent < 1){ $returnSpacePercent = 1; }
		
		$spaceUsedH = $core->framework->files->formatSize($spaceUsed);
		$maxSpaceH = $core->framework->files->formatSize($maxSpace);
		
		echo '	<div class="progress">
		  			<div class="progress-bar" style="width:'.$returnSpacePercent.'%"></div>
				</div>
				<p class="text-muted">You are using '.$spaceUsedH.' of your maximum '.$core->framework->server->getData('disk_space').' MB of disk space.</p>';
				
	}else if($_POST['command'] && $_POST['command'] == 'players'){
					
		if($core->framework->rcon->online($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
			
			$core->framework->rcon->getStatus($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'));
			
									
				$pOnlineL = '';
				if(count($core->framework->rcon->data('players')) > 0){
					foreach($core->framework->rcon->data('players') as $id => $player){
				
						$pOnlineL .= '<img data-toggle="tooltip" src="http://i.fishbans.com/helm/'.$player.'/32" title="'.$player.'" style="padding: 0 2px 6px 0;"/>';
				
					}
				}else{
				
					$pOnlineL = '<p class="text-muted">No players are currently online.</p>';
				
				}
			
			$pPercent = round(count($core->framework->rcon->data('players'))/$core->framework->rcon->data('maxplayers'), 2) * 100;
			($pPercent < 1) ? $pPercent = '1' : $pPercent = $pPercent;
			
			echo '
					<div class="progress">
					  	<div class="progress-bar" style="width:'.$pPercent.'%"></div>
					</div>
					'.$pOnlineL;
				
		}else{
		
			echo '
					<div class="alert alert-info">Unable to connect and query the server. Is it online?</div>
				';
		
		}
	
	}else if($_POST['command'] && $_POST['command'] == 'info'){
	
		if($core->framework->rcon->online($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
			
			$core->framework->rcon->getStatus($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'));
			$sVersion = $core->framework->rcon->data('version');
							
		}else{
		
			$sPlugins = 'Unable to query the server.';
			$sVersion = 'Unable to query the server.'; 
			$sSoftware = 'Unable to query the server.';
		
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
							<td><strong>IP</strong></td>
							<td>'.$core->framework->server->getData('server_ip').'</td>
						</tr>
						<tr>
							<td><strong>Port</strong></td>
							<td>'.$core->framework->server->getData('server_port').'</td>
						</tr>
						<tr>
							<td><strong>Node</strong></td>
							<td>'.$core->framework->settings->nodeName($core->framework->server->getData('node')).'</td>
						</tr>
						<tr>
							<td><strong>Version</strong></td>
							<td>'.$sVersion.'</td>
						</tr>
					</tbody>
				</table>';
	
	}

}else{

	exit('Invalid Authentication.');

}
?>