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

require_once('../../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	exit('<div class="alert alert-warning">Failed to Authenticate Account.</div>');
}

/*
 * Check Variables
 */
if(!isset($_POST['method'], $_POST['field_1'], $_POST['operator_1'], $_POST['term_1'], $_POST['mid_op'], $_POST['field_2'], $_POST['operator_2'], $_POST['term_2']))
	exit('<div class="alert alert-warning">Missing required variable.</div>'.print_r($_POST));

if($_POST['method'] != 'advanced')
	exit('<div class="alert alert-warning">Invalid Search Method.</div>');

if(empty($_POST['field_1']) || empty($_POST['operator_1']) || empty($_POST['mid_op']) || empty($_POST['field_2']) || empty($_POST['operator_2']))
	exit('<div class="alert alert-warning">Required Variable Empty.</div>');

if(!in_array($_POST['field_1'], array('email', 'username', 'root_admin')) || !in_array($_POST['field_2'], array('email', 'username', 'root_admin')))
	exit('<div class="alert alert-warning">Required `field` contains unknown value.</div>');

if(!in_array($_POST['operator_1'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')) || !in_array($_POST['operator_2'], array('equal', 'not_equal', 'starts_w', 'ends_w', 'like')))
	exit('<div class="alert alert-warning">Required `operator` contains unknown value.</div>');

if(strlen($_POST['term_1']) < 4 && $_POST['field_1'] != 'root_admin' || strlen($_POST['term_2']) < 4 && $_POST['field_2'] != 'root_admin')
	exit('<div class="alert alert-warning">Required `term` must be at least 4 characters.</div>');

if($_POST['field_1'] == 'root_admin' && !in_array($_POST['term_1'], array('0', '1')) || $_POST['field_2'] == 'root_admin' && !in_array($_POST['term_2'], array('0', '1')))
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


$find = ORM::forTable('users')->rawQuery("SELECT * FROM `users` WHERE `".$_POST['field_1']."` ".$useOperator." '".$searchTerm."' ".$middleOperator." `".$_POST['field_2']."` ".$useOperator_2." '".$searchTerm_2."'")->findMany();

	$returnRows = '';
	foreach($find as &$row){

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
