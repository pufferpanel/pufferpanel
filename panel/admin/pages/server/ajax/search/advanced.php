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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), null, true) !== true){
	exit('<div class="alert alert-warning">Failed to Authenticate Account.</div>');
}

/*
 * Check Variables
 */
if(!isset($_POST['method'], $_POST['field_1'], $_POST['operator_1'], $_POST['term_1'], $_POST['mid_op'], $_POST['field_2'], $_POST['operator_2'], $_POST['term_2']))
	exit('<div class="alert alert-warning">Missing required variable.</div>');

if($_POST['method'] != 'advanced')
	exit('<div class="alert alert-warning">Invalid Search Method.</div>');
	
if(empty($_POST['field_1']) || empty($_POST['operator_1']) || empty($_POST['mid_op']) || empty($_POST['field_2']) || empty($_POST['operator_2']))
	exit('<div class="alert alert-warning">Required Variable Empty.</div>');
	
if(!in_array($_POST['field_1'], array('name', 'server_ip', 'server_port', 'owner_email', 'active')) || !in_array($_POST['field_2'], array('name', 'server_ip', 'owner_email', 'active')))
	exit('<div class="alert alert-warning">Required `field` contains unknown value.</div>');
	
if(!in_array($_POST['operator_1'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')) || !in_array($_POST['operator_2'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')))
	exit('<div class="alert alert-warning">Required `operator` contains unknown value.</div>');

if(strlen($_POST['term_1']) < 4 && $_POST['field_1'] != 'active' || strlen($_POST['term_2']) < 4 && $_POST['field_2'] != 'active')
	exit('<div class="alert alert-warning">Required `term` must be at least 4 characters.</div>');
	
if($_POST['field_1'] == 'active' && !in_array($_POST['term_1'], array('0', '1')) || $_POST['field_2'] == 'active' && !in_array($_POST['term_2'], array('0', '1')))
	exit('<div class="alert alert-warning">Required `term` for root_admin must be 1 or 0.</div>');
	
if(!in_array($_POST['mid_op'], array('and', 'or')))
	exit('<div class="alert alert-warning">Required `comparison_operator` must be AND/OR.</div>');
	
/*
 * Is Search Looking for Similar
 */
if($_POST['operator_1'] == 'starts_w'){
	$searchTerm = $_POST['term_1'].'%';
	$useOperator = 'LIKE';
}else if($_POST['operator_1'] == 'ends_w'){
	$searchTerm = '%'.$_POST['term_1'];
	$useOperator = 'LIKE';
}else if($_POST['operator_1'] == 'like'){
	$searchTerm = '%'.$_POST['term_1'].'%';
	$useOperator = 'LIKE';
}else if($_POST['operator_1'] == 'not_equal'){
	$searchTerm = $_POST['term_1'];
	$useOperator = '!=';
}else if($_POST['operator_1'] == 'equal'){
	$searchTerm = $_POST['term_1'];
	$useOperator = '=';
}

	/*
	 * Is Search Looking for Similar
	 */
	if($_POST['operator_2'] == 'starts_w'){
		$searchTerm_2 = $_POST['term_2'].'%';
		$useOperator_2 = 'LIKE';
	}else if($_POST['operator_2'] == 'ends_w'){
		$searchTerm_2 = '%'.$_POST['term_2'];
		$useOperator_2 = 'LIKE';
	}else if($_POST['operator_2'] == 'like'){
		$searchTerm_2 = '%'.$_POST['term_2'].'%';
		$useOperator_2 = 'LIKE';
	}else if($_POST['operator_2'] == 'not_equal'){
		$searchTerm_2 = $_POST['term_2'];
		$useOperator_2 = '!=';
	}else if($_POST['operator_2'] == 'equal'){
		$searchTerm_2 = $_POST['term_2'];
		$useOperator_2 = '=';
	}
	
/*
 * Comparison
 */
if($_POST['mid_op'] == 'and')
	$middleOperator = 'AND';
else
	$middleOperator = 'OR';

/*
 * Different Search Method for Owner Email
 */
if($_POST['field_1'] == 'owner_email' || $_POST['field_2'] == 'owner_email'){

	
	if($_POST['field_1'] == 'owner_email'){
		
		$findIDs = $mysql->prepare("SELECT `id` FROM `users` WHERE `email` ".$useOperator." :term");
		$findIDs->execute(array(
			':term' => $searchTerm
		));
		
		$otherField = $_POST['field_2'];
		$otherSearchOp = $useOperator_2;
		$otherSearch = $searchTerm_2;
		
	}else{
		
		$findIDs = $mysql->prepare("SELECT `id` FROM `users` WHERE `email` ".$useOperator_2." :term");
		$findIDs->execute(array(
			':term' => $searchTerm_2
		));
		
		$otherField = $_POST['field_1'];
		$otherSearchOp = $useOperator;
		$otherSearch = $searchTerm;
		
	}
	
	$find = $mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` IN (".implode(',', $findIDs->fetchAll(PDO::FETCH_COLUMN, 0)).") ".$middleOperator." `".$otherField."` ".$otherSearchOp." :term");
		
	$find->execute(array(
		':term' => $otherSearch
	));

}else{
	
	$find = $mysql->prepare("SELECT * FROM `servers` WHERE `".$_POST['field_1']."` ".$useOperator." :term ".$middleOperator." `".$_POST['field_2']."` ".$useOperator_2." :term_2");
	$find->execute(array(
		':term' => $searchTerm,
		':term_2' => $searchTerm_2
	));

}


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
			<td><a href="../node/view.php?id='.$row['node'].'">'.$core->framework->settings->nodeName($row['node']).'</a></td>
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