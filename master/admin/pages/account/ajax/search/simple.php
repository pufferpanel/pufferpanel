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
	exit('<div class="alert alert-warning">Failed to Authenticate Account.</div>');
}

/*
 * Check Variables
 */
if(!isset($_POST['method'], $_POST['field'], $_POST['operator'], $_POST['term']))
	exit('<div class="alert alert-warning">Missing required variable.</div>');

if($_POST['method'] != 'simple')
	exit('<div class="alert alert-warning">Invalid Search Method.</div>');
	
if(empty($_POST['field']) || empty($_POST['operator']))
	exit('<div class="alert alert-warning">Required Variable Empty.</div>');
	
if(!in_array($_POST['field'], array('email', 'username', 'root_admin')))
	exit('<div class="alert alert-warning">Required `field` contains unknown value.</div>');
	
if(!in_array($_POST['operator'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')))
	exit('<div class="alert alert-warning">Required `operator` contains unknown value.</div>');

if(strlen($_POST['term']) < 4 && $_POST['field'] != 'root_admin')
	exit('<div class="alert alert-warning">Required `term` must be at least 4 characters.</div>');
	
if($_POST['field'] == 'root_admin' && !in_array($_POST['term'], array('0', '1')))
	exit('<div class="alert alert-warning">Required `term` for root_admin must be 1 or 0.</div>');
	
/*
 * Is Search Looking for Similar
 */
if($_POST['operator'] == 'starts_w'){
	$searchTerm = $_POST['term'].'%';
	$useOperator = 'LIKE';
}else if($_POST['operator'] == 'ends_w'){
	$searchTerm = '%'.$_POST['term'];
	$useOperator = 'LIKE';
}else if($_POST['operator'] == 'like'){
	$searchTerm = '%'.$_POST['term'].'%';
	$useOperator = 'LIKE';
}else if($_POST['operator'] == 'not_equal'){
	$searchTerm = $_POST['term'];
	$useOperator = '!=';
}else if($_POST['operator'] == 'equal'){
	$searchTerm = $_POST['term'];
	$useOperator = '=';
}
	
	
$find = $mysql->prepare("SELECT * FROM `users` WHERE `".$_POST['field']."` ".$useOperator." :term");
$find->execute(array(
	':term' => $searchTerm
));

	$returnRows = '';
	while($row = $find->fetch()){
		
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