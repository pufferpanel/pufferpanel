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
	
	if($_POST['data'] && $_POST['data'] == 'memory')
		exit($core->framework->files->format($core->framework->gsd->retrieve_process('memory')));
	
	if($_POST['data'] && $_POST['data'] == 'cpu'){
		$cpu = round(($core->framework->gsd->retrieve_process('cpu') / $core->framework->server->getData('cpu_limit')) * 100, 2);
		$cpu = ($cpu > "100") ? "100" : $cpu;
		
		exit($cpu);
		
	}
	
	if($_POST['data'] && $_POST['data'] == 'players'){
	
		$onlinePlayers = null;
		if(count(($players = $core->framework->gsd->retrieve('players'))) > 0)
			foreach($players as $player)
				$onlinePlayers .= '<img data-toggle="tooltip" src="http://i.fishbans.com/helm/'.$player.'/32" title="'.$player.'" style="padding: 0 2px 6px 0;"/>';
		else
			$onlinePlayers = '<p class="text-muted">No players are currently online.</p>';
			
		exit($onlinePlayers);
		
	}
	
}
?>