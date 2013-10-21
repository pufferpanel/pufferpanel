<?php
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	exit('<div class="error-box round">Failed to Authenticate Account.</div>');
}

/*
 * Check Variables
 */
if(!isset($_POST['method'], $_POST['field'], $_POST['operator'], $_POST['term']))
	exit('<div class="error-box round">Missing required variable.</div>');

if($_POST['method'] != 'simple')
	exit('<div class="error-box round">Invalid Search Method.</div>');
	
if(empty($_POST['field']) || empty($_POST['operator']))
	exit('<div class="error-box round">Required Variable Empty.</div>');
	
if(!in_array($_POST['field'], array('node', 'server_name', 'server_ip', 'owner_email', 'active')))
	exit('<div class="error-box round">Required `field` contains unknown value.</div>');
	
if(!in_array($_POST['operator'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')))
	exit('<div class="error-box round">Required `operator` contains unknown value.</div>');

if(strlen($_POST['term']) < 4 && $_POST['field'] != 'active')
	exit('<div class="error-box round">Required `term` must be at least 4 characters.</div>');
	
if($_POST['field'] == 'active' && !in_array($_POST['term'], array('0', '1')))
	exit('<div class="error-box round">Required `term` for active must be 1 or 0.</div>');
	
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

/*
 * Different Search Method for Owner Email
 */
if($_POST['field'] == 'owner_email'){

	$findIDs = $mysql->prepare("SELECT `id` FROM `users` WHERE `email` ".$useOperator." :term");
	$findIDs->execute(array(
		':term' => $searchTerm
	));
	
	$find = $mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` IN (".implode(',', $findIDs->fetchAll(PDO::FETCH_COLUMN, 0)).")");
	$find->execute();

}else{

	$find = $mysql->prepare("SELECT * FROM `servers` WHERE `".$_POST['field']."` ".$useOperator." :term");
	$find->execute(array(
		':term' => $searchTerm
	));

}

	$returnRows = '';
	while($row = $find->fetch()){
		
		$isActive = ($row['active'] == 1) ? '<i class="icon-ok-circle"></i>' : '<i class="icon-remove-circle"></i>';
		
		$returnRows .= '
		<tr>
			<td><a href="../account/view.php?id='.$row['owner_id'].'">'.$row['owner_id'].'</a></td>
			<td style="text-align:center;">'.$row['whmcs_id'].'</td>
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
			<th style="width:7%">Owner ID</th>
			<th style="width:8%;text-align:center;">WHMCS ID</th>
			<th style="width:20%">Server Name</th>
			<th style="width:10%">Node</th>
			<th style="width:30%">Connection Address</th>
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