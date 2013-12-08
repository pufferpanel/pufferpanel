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
	
$find = $mysql->prepare("SELECT * FROM `users`");
$find->execute();

	$returnRows = '';
	while($row = $find->fetch()){
		
		$isRoot = ($row['root_admin'] == 1) ? '<i class="fa fa-check-circle-o"></i>' : '<i class="fa fa-times-circle-o"></i>';
		
		$returnRows .= '
		<tr>
			<td><a href="view.php?id='.$row['id'].'">'.$row['username'].'</a></td>
			<td>'.$row['email'].'</td>
			<td>'.date('r', $row['register_time']).'</td>
			<td style="text-align:center;">'.$isRoot.'</td>
		</tr>
		';
	
	}

echo '
<table>
	<thead>
		<tr>
			<th style="width:20%">Username</th>
			<th style="width:30%">Email</th>
			<th style="width:35%">Registered</th>
			<th style="width:5%;text-align:center;">Admin</th>
		</tr>
	</thead>
	<tbody>
		'.$returnRows.'
	</tbody>
</table>';

?>