<?php
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
	
if(!in_array($_POST['field_1'], array('email', 'username', 'root_admin')) || !in_array($_POST['field_2'], array('email', 'username', 'root_admin')))
	exit('<div class="error-box round">Required `field` contains unknown value.</div>');
	
if(!in_array($_POST['operator_1'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')) || !in_array($_POST['operator_2'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')))
	exit('<div class="error-box round">Required `operator` contains unknown value.</div>');

if(strlen($_POST['term_1']) < 4 && $_POST['field_1'] != 'root_admin' || strlen($_POST['term_2']) < 4 && $_POST['field_2'] != 'root_admin')
	exit('<div class="error-box round">Required `term` must be at least 4 characters.</div>');
	
if($_POST['field_1'] == 'root_admin' && !in_array($_POST['term_1'], array('0', '1')) || $_POST['field_2'] == 'root_admin' && !in_array($_POST['term_2'], array('0', '1')))
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


$find = $mysql->prepare("SELECT * FROM `users` WHERE `".$_POST['field_1']."` ".$useOperator." :term ".$middleOperator." `".$_POST['field_2']."` ".$useOperator_2." :term_2");
$find->execute(array(
	':term' => $searchTerm,
	':term_2' => $searchTerm_2
));

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