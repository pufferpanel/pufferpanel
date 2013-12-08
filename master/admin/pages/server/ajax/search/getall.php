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
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}
	
$find = $mysql->prepare("SELECT * FROM `servers`");
$find->execute();

	$returnRows = '';
	while($row = $find->fetch()){
		
		$isActive = ($row['active'] == 1) ? '<i class="fa fa-check-circle-o"></i>' : '<i class="fa fa-times-circle-o"></i>';
		
		$findUser = $mysql->prepare("SELECT `email` FROM `users` WHERE `id`  = :id");
		$findUser->execute(array(
			':id' => $row['owner_id']
		));
		$user = $findUser->fetch();
		
		$returnRows .= '
		<tr>
			<td><a href="../../../servers.php?goto='.$row['hash'].'"><i class="fa fa-terminal"></i></a></td>
			<td><a href="../account/view.php?id='.$row['owner_id'].'">'.$user['email'].'</a> (uID #'.$row['owner_id'].')</td>
			<td><a href="view.php?id='.$row['id'].'">'.$row['name'].'</a> ('.$row['ftp_user'].')</td>
			<td><a href="../node/view.php?do=redirect&node='.$row['node'].'">'.$row['node'].'</a></td>
			<td>'.$row['server_ip'].':'.$row['server_port'].'</td>
			<td>'.$row['max_ram'].' MB</td>
			<td>'.$row['disk_space'].' MB</td>
			<td style="text-align:center;">'.$isActive.'</td>
		</tr>
		';
	
	}

echo '
<table>
	<thead>
		<tr>
			<th style="width:5%"></th>
			<th style="width:20%">Owner Information</th>
			<th style="width:20%">Server Name (FTP User)</th>
			<th style="width:10%">Node</th>
			<th style="width:20%">Connection Address</th>
			<th style="width:10%">RAM</th>
			<th style="width:10%">Disk Space</th>
			<th style="width:5%;text-align:center;">Active</th>
		</tr>
	</thead>
	<tbody>
		'.$returnRows.'
	</tbody>
</table>';

?>