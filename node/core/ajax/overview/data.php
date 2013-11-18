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
		  			<div class="progress-bar" style="width:'.$returnSpacePercent.'%;background-color: #2069b4"></div>
				</div>
				<p class="center nomargin">You are using '.$spaceUsedH.' of your maximum '.$core->framework->server->getData('disk_space').' MB of disk space.</p>';
				
	}else if($_POST['command'] && $_POST['command'] == 'players'){
					
		if($core->framework->rcon->online($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
			
			$core->framework->rcon->getStatus($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'));
			
									
				$pOnlineL = '';
				if(count($core->framework->rcon->data('players')) > 0){
					foreach($core->framework->rcon->data('players') as $id => $player){
				
						$pOnlineL .= '<img src="http://i.fishbans.com/player/'.$player.'/32" alt="'.$player.'" title="'.$player.'" style="padding:5px 10px;"/>';
				
					}
				}else{
				
					$pOnlineL = '<p class="center nomargin">No players are currently online.</p>';
				
				}
			
			$pPercent = round(count($core->framework->rcon->data('players'))/$core->framework->rcon->data('maxplayers'), 2) * 100;
			($pPercent < 1) ? $pPercent = '1' : $pPercent = $pPercent;
			
			echo '
					<div class="progress">
					  	<div class="progress-bar" style="width:'.$pPercent.'%;background-color: #2069b4"></div>
					</div>
					'.$pOnlineL;
				
		}else{
		
			echo '
					<div class="progress">
					  	<div class="progress-bar" style="width:1%;background-color: #2069b4"></div>
					</div>
					<p class="center nomargin">Could not connect to server via RCON.</p>
				';
		
		}
	
	}else if($_POST['command'] && $_POST['command'] == 'info'){
	
		if($core->framework->rcon->online($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
			
			$core->framework->rcon->getStatus($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'));
			$serverStatus = '<span style="color:green;">Online</span>';
			$sVersion = $core->framework->rcon->data('version');
			$sSoftware = $core->framework->rcon->data('software');
			$sPlugins = null;
			foreach(explode(',', $core->framework->rcon->data('plugins')) as $id => $plugin)
				{
				
					$pData = explode(' ', $plugin);
					$sPlugins .= str_replace(end($pData), '', $plugin).' ('.end($pData).')';
				
				}
							
		}else{
		
			$sPlugins = 'Could not connect to server via RCON.';
			$serverStatus = '<span style="color:red;">Offline</span>';
			$sVersion = 'Could not connect to server via RCON.'; 
			$sSoftware = 'Could not connect to server via RCON.';
		
		}
		
		echo '
				<table>
					<thead>
						<tr>
							<th style="width:1%"></th>
							<th style="width:20%">Information</th>
							<th style="width:79%">Data</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>&nbsp;</td>
							<td><strong>IP</strong></td>
							<td>'.$core->framework->server->getData('server_ip').'</td>
						</tr>
						<tr>
							<td>&nbsp;</td>
							<td><strong>Port</strong></td>
							<td>'.$core->framework->server->getData('server_port').'</td>
						</tr>
						<tr>
							<td>&nbsp;</td>
							<td><strong>Node</strong></td>
							<td>'.$core->framework->server->getData('node').'</td>
						</tr>
						<tr>
							<td>&nbsp;</td>
							<td><strong>Version</strong></td>
							<td>'.$sVersion.'</td>
						</tr>
						<tr>
							<td>&nbsp;</td>
							<td><strong>Software</strong></td>
							<td>'.$sSoftware.'</td>
						</tr>
						<tr>
							<td>&nbsp;</td>
							<td><strong>Plugins</strong></td>
							<td>'.$sPlugins.'</td>
						</tr>
					</tbody>
				</table>';
	
	}

}else{

	exit('Invalid Authentication.');

}
?>