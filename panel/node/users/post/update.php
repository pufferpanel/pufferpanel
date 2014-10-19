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
require_once('../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php?login');
	exit();

}

	if(!isset($_POST['uuid'], $_POST['permissions']))
		Page\components::redirect('../list.php');

	if($core->auth->XSRF(@$_POST['xsrf']) !== true)
		Page\components::redirect('../list.php?id='.$_POST['uuid'].'&error');

	if(empty($_POST['permissions']))
		Page\components::redirect('../view.php?id='.$_POST['uuid'].'&error');

	$query = $mysql->prepare("SELECT `permissions` FROM `users` WHERE `uuid` = :uid");
	$query->execute(array(
		':uid' => $_POST['uuid']
	));

		if($query->rowCount() != 1)
			Page\components::redirect('../list.php?error');
		else
			$row = $query->fetch();

		$permissions = @json_decode($row['permissions'], true);
		if(!is_array($permissions) || !array_key_exists($core->server->getData('hash'), $permissions))
			Page\components::redirect('../view.php?id='.$_POST['uuid'].'&error');

		$permissions[$core->server->getData('hash')] = $_POST['permissions'];

	$query = $mysql->prepare("UPDATE `users` SET `permissions` = :permissions WHERE `uuid` = :uid");
	$query->execute(array(
		':permissions' => json_encode($permissions),
		':uid' => $_POST['uuid']
	));

	Page\components::redirect('../view.php?id='.$_POST['uuid']);

?>