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
namespace PufferPanel\Core;
use \ORM as ORM;

require_once('../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}

$users = ORM::forTable('users')->findMany();

	$returnRows = '';
	foreach($users as &$row){

		$isRoot = ($row['root_admin'] == 1) ? '<span class="label label-danger">Admin</span>' : '<span class="label label-success">User</span>';

		$returnRows .= '
		<tr>
			<td><a href="view.php?id='.$row['id'].'">'.$row['username'].'</a></td>
			<td>'.$row['email'].'</td>
			<td>'.date('r', $row['register_time']).'</td>
			<td style="text-align:center;">'.$isRoot.'</td>
		</tr>
		';

	}

echo '<table class="table table-striped table-bordered table-hover">
	<thead>
		<tr>
			<th>Username</th>
			<th>Email</th>
			<th>Registered</th>
			<th></th>
		</tr>
	</thead>
	<tbody>
		'.$returnRows.'
	</tbody>
</table>';

?>
