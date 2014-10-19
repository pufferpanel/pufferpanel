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
require_once('../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php?login');
	exit();
}

if($core->user->hasPermission('users.view') !== true)
	Page\components::redirect('../index.php?error=no_permission');

$access = json_decode($core->server->getData('subusers'), true);
$users = array();

if(is_array($access) && !empty($access)){

	foreach($access as $id => $status) {

		if($status != "verified"){

			$query = $mysql->prepare("SELECT `content` FROM `account_change` WHERE `key` = :key AND `verified` = 0");
			$query->execute(array(
				':key' => $core->auth->encrypt($id, $status).".".$status
			));

			$row = $query->fetch();

			$_perms = json_decode($row['content'], true);
			$users = array_merge($users, array($id => array("status" => "pending", "permissions" => $_perms[$core->server->getData('hash')])));

		}else{

			$query = $mysql->prepare("SELECT `permissions`, `email`, `uuid` FROM `users` WHERE `id` = :id");
			$query->execute(array(
				':id' => $id
			));

			$row = $query->fetch();

			$_perms = json_decode($row['permissions'], true);
			$users = array_merge($users, array($row['email'] => array("status" => "verified", "id" => $row['uuid'], "permissions" => $_perms[$core->server->getData('hash')])));

		}

	}

}

/*
* Display Page
*/
echo $twig->render(
		'node/users/list.html', array(
			'users' => $users,
			'server' => $core->server->getData(),
			'footer' => array(
				'queries' => Database\databaseInit::getCount(),
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>