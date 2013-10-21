<?php
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash')) === true){

	if($_POST['command'] && $_POST['command'] == 'stats'){
	
		$maxSpace = ($core->framework->server->getData('disk_space') * 1024) * 1024;
		$spaceUsed = $core->framework->files->getFolderSize($core->framework->server->getData('path'));
		    
		$returnSpacePercent = round(($spaceUsed / $maxSpace), 2) * 100;
		if($returnSpacePercent < 1){ $returnSpacePercent = 1; }
		
		$spaceUsedH = $core->framework->files->formatSize($spaceUsed);
		$maxSpaceH = $core->framework->files->formatSize($maxSpace);
		
		echo '	<div class="progress">
		  			<div class="progress-bar" style="width:'.$returnSpacePercent.'%;background-color: #2069b4"></div>
				</div>
				<p class="center nomargin">You are using '.$spaceUsedH.' of your maximum '.$maxSpaceH.' of disk space.</p>';
				
	}else if($_POST['command'] && $_POST['command'] == 'players'){
	
		if($rcon->s->isOnline($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
				
			include('../../../core/framework/rcon/query.class.php');
			include('../../../core/framework/rcon/rcon.class.php');
			include('../../../core/framework/rcon/query.status.php');
			
			$rcon->status = new MinecraftStatus($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'));
			$rcon->query = new MinecraftQuery();
			$rcon->command = new MinecraftRcon();
			
			$rcon->query->Connect($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'), 1);
			$playersOnline = $rcon->query->GetPlayers();
			$pl = $rcon->query->GetInfo();
				
				$pOnlineL = '';
				if(!empty($playersOnline)){
					foreach($playersOnline as $id => $player){
				
						$pOnlineL .= '<img src="http://i.fishbans.com/player/'.$player.'/32" alt="'.$player.'" title="'.$player.'" style="padding:5px 10px;"/>';
				
					}
				}else{
				
					$pOnlineL = '<p class="center nomargin">No players are currently online.</p>';
				
				}
			
			$pPercent = round($pl['Players']/$pl['MaxPlayers'], 2) * 100;
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
	
		if($rcon->s->isOnline($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port')) === true){
				
			include('../../../core/framework/rcon/query.class.php');
			include('../../../core/framework/rcon/rcon.class.php');
			include('../../../core/framework/rcon/query.status.php');
			
			$rcon->status = new MinecraftStatus($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'));
			$rcon->query = new MinecraftQuery();
			$rcon->command = new MinecraftRcon();
			
			$rcon->query->Connect($core->framework->server->getData('server_ip'), $core->framework->server->getData('server_port'), 1);
			$pl = $rcon->query->GetInfo();
			
			$pluginList = '';
			if(is_array($pl['Plugins'])){			
				foreach($pl['Plugins'] as $id => $plugin){
					if(strpos(strtolower($plugin), 'mod') === false){
						$pluginList .= $plugin.', ';
					}
				}
			}
			
			$serverStatus = '<span style="color:green;">Online</span>';
			$sVersion = $pl['Version']; 
			$sSoftware = $pl['Software'];
				
		}else{
		
			$pluginList = 'Could not connect to server via RCON.';
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
							<td>'.$pluginList.'</td>
						</tr>
					</tbody>
				</table>';
	
	}

}else{

	exit('Invalid Authentication.');

}
?>