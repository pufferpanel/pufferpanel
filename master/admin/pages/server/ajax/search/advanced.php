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

/*
 * Check Variables
 */
if(!isset($_POST['method'], $_POST['field_1'], $_POST['operator_1'], $_POST['term_1'], $_POST['mid_op'], $_POST['field_2'], $_POST['operator_2'], $_POST['term_2']))
	exit('<div class="error-box round">Missing required variable.</div>');

if($_POST['method'] != 'advanced')
	exit('<div class="error-box round">Invalid Search Method.</div>');
	
if(empty($_POST['field_1']) || empty($_POST['operator_1']) || empty($_POST['mid_op']) || empty($_POST['field_2']) || empty($_POST['operator_2']))
	exit('<div class="error-box round">Required Variable Empty.</div>');
	
if(!in_array($_POST['field_1'], array('node', 'name', 'server_ip', 'owner_email', 'active')) || !in_array($_POST['field_2'], array('node', 'name', 'server_ip', 'owner_email', 'active')))
	exit('<div class="error-box round">Required `field` contains unknown value.</div>');
	
if(!in_array($_POST['operator_1'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')) || !in_array($_POST['operator_2'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')))
	exit('<div class="error-box round">Required `operator` contains unknown value.</div>');

if(strlen($_POST['term_1']) < 4 && $_POST['field_1'] != 'active' || strlen($_POST['term_2']) < 4 && $_POST['field_2'] != 'active')
	exit('<div class="error-box round">Required `term` must be at least 4 characters.</div>');
	
if($_POST['field_1'] == 'active' && !in_array($_POST['term_1'], array('0', '1')) || $_POST['field_2'] == 'active' && !in_array($_POST['term_2'], array('0', '1')))
	exit('<div class="error-box round">Required `term` for root_admin must be 1 or 0.</div>');
	
if(!in_array($_POST['mid_op'], array('and', 'or')))
	exit('<div class="error-box round">Required `comparison_operator` must be AND/OR.</div>');
	
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
		
		$isActive = ($row['active'] == 1) ? '<i class="fa fa-check-circle-o"></i>' : '<i class="fa fa-times-circle-o"></i>';
		
		$find = $mysql->prepare("SELECT `email` FROM `users` WHERE `id`  = :id");
		$find->execute(array(
			':id' => $row['owner_id']
		));
		$user = $find->fetch();
		
		$returnRows .= '
		<tr>
			<td><a href="../../../servers.php?goto='.$row['hash'].'"><i class="fa fa-terminal"></i></a></td>
			<td><a href="../account/view.php?id='.$row['owner_id'].'">'.$user['email'].'</a> (uID #'.$row['owner_id'].')</td>
			<td><a href="view.php?id='.$row['id'].'">'.$row['name'].'</a></td>
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
			<th style="width:20%">Server Name</th>
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