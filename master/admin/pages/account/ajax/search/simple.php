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
	
if(!in_array($_POST['field'], array('email', 'username', 'root_admin')))
	exit('<div class="error-box round">Required `field` contains unknown value.</div>');
	
if(!in_array($_POST['operator'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')))
	exit('<div class="error-box round">Required `operator` contains unknown value.</div>');

if(strlen($_POST['term']) < 4 && $_POST['field'] != 'root_admin')
	exit('<div class="error-box round">Required `term` must be at least 4 characters.</div>');
	
if($_POST['field'] == 'root_admin' && !in_array($_POST['term'], array('0', '1')))
	exit('<div class="error-box round">Required `term` for root_admin must be 1 or 0.</div>');
	
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
		
		$isRoot = ($row['root_admin'] == 1) ? '<i class="fa fa-ok-circle"></i>' : '<i class="fa fa-remove-circle"></i>';
		
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