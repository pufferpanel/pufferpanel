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
require_once('../../../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}
	
$find = $mysql->prepare("SELECT * FROM `servers`");
$find->execute();

	$returnRows = '';
	while($row = $find->fetch()){
		
		$isActive = ($row['active'] == '1') ? '<span class="label label-success">Enabled</span>' : '<span class="label label-danger">Disabled</span>';
		
		$findUser = $mysql->prepare("SELECT `email` FROM `users` WHERE `id`  = :id");
		$findUser->execute(array(
			':id' => $row['owner_id']
		));
		$user = $findUser->fetch();
		
		$row['name'] = (strlen($row['name']) > 20) ? substr($row['name'], 0, 17).'...' : $row['name'];
		$user['email'] = (strlen($user['email']) > 25) ? substr($user['email'], 0, 22).'...' : $user['email'];
		
		$returnRows .= '
		<tr>
			<td><a href="../../../servers.php?goto='.$row['hash'].'"><i class="fa fa-tachometer"></i></a></td>
			<td><a href="view.php?id='.$row['id'].'">'.$row['name'].'</a></td>
			<td><a href="../node/view.php?id='.$row['node'].'">'.$core->settings->nodeName($row['node']).'</a></td>
			<td>'.$row['server_ip'].':'.$row['server_port'].'</td>
			<td><a href="../account/view.php?id='.$row['owner_id'].'">'.$user['email'].'</a></td>
			<td style="text-align:center;">'.$isActive.'</td>
		</tr>
		';
	
	}

echo '
<table class="table table-striped table-bordered table-hover">
	<thead>
		<tr>
			<th style="width:2%"></th>
			<th>Server Name</th>
			<th>Node</th>
			<th>Connection Address</th>
			<th>Owner Information</th>
			<th></th>
		</tr>
	</thead>
	<tbody>
		'.$returnRows.'
	</tbody>
</table>';

?>
