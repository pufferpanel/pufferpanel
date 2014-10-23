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

//id means pending, uid means not pending
if(isset($_GET['id']) && !empty($_GET['id'])){

	$_GET['id'] = urldecode($_GET['id']);

	$query = $mysql->prepare("SELECT * FROM `account_change` WHERE `key` = :key AND `verified` = 0");
	$query->execute(array(
		':key' => $_GET['id']
	));

	if($query->rowCount() != 1)
		Page\components::redirect('../list.php?error');

	$row = $query->fetch();

	// verify that this user is assigned to this server
	if(!array_key_exists($core->server->getData('hash'), json_decode($row['content'], true)))
		Page\components::redirect('../list.php?error=c1');

	// remove verification codes
	$mysql->exec("DELETE FROM `account_change` WHERE `id` = ".$row['id']);

	// update server in database
	list($encrypted, $iv) = explode('.', $_GET['id']);
	$_perms = json_decode($row['subusers'], true);
	unset($_perms[$core->auth->decrypt($encrypted, $iv)]);
	$_perms = json_encode($_perms);
	$mysql->exec("UPDATE `servers` SET `subusers` = '".$_perms."' WHERE `hash` = '".$core->server->getData('hash')."'");

	Page\components::redirect('../list.php?revoked');

}elseif(isset($_GET['uid']) && !empty($_GET['uid'])){

	$_GET['uid'] = urldecode($_GET['uid']);

	$query = $mysql->prepare("SELECT * FROM `users` WHERE `uuid` = :uuid");
	$query->execute(array(
		':uuid' => $_GET['uid']
	));

	if($query->rowCount() != 1)
		Page\components::redirect('../list.php?error');

	$row = $query->fetch();

	// verify that this user is assigned to this server
	if(!array_key_exists($core->server->getData('hash'), json_decode($row['permissions'], true)))
		Page\components::redirect('../list.php?error');

	// update server in database
	$_perms = json_decode($core->server->getData('subusers'), true);
	unset($_perms[$row['id']]);
	$_perms = json_encode($_perms);
	$mysql->exec("UPDATE `servers` SET `subusers` = '".$_perms."' WHERE `hash` = '".$core->server->getData('hash')."'");

	Page\components::redirect('../list.php?revoked');

}
?>