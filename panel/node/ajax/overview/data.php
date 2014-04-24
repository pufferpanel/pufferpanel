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
		
	$maxSpace = $core->server->getData('disk_space') * 1024 * 1024;
		
		/*
		 * Run Command
		 */
		$getCommandData = $core->ssh->generateSSH2Connection($core->server->nodeData('id'), true)->executeSSH2Command('sudo du -s '.$core->server->nodeData('server_dir').$core->server->getData('ftp_user').'/server', true);
					
		if($getCommandData === false)
			exit('<div class="alert alert-danger">'.$_l->tpl('node.overview.ajax.no_conn').'</div>');
		else if(empty($getCommandData))
			exit('<div class="alert alert-danger">'.$_l->tpl('node.overview.ajax.no_exec').'</div>');
	
		
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
			<p class="text-muted">'.sprintf($_l->tpl('node.overview.space_used'), $spaceUsedH, $core->server->getData('disk_space')).'</p>';
			
}else{

	exit('Invalid Authentication.');

}
?>