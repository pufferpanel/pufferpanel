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
require_once('../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) === true)
	Page\components::redirect('account.php?token='.@$_GET['token']);

if(isset($_GET['do']) && $_GET['do'] == 'register' && $_SERVER['REQUEST_METHOD'] === 'POST'){

	if(!isset($_POST['token']))
		Page\components::redirect('register.php?error=token');
	else
		list($encrypted, $iv) = explode('.', $_POST['token']);

	/* XSRF Check */
	if($core->auth->XSRF(@$_POST['xsrf']) !== true)
		Page\components::redirect('index.php?error=xsrf&token='.urlencode($_POST['token']));

	if(!preg_match('/^[\w-]{4,35}$/', $_POST['username']))
		Page\components::redirect('register.php?error=u_fail&token='.urlencode($_POST['token']));

	if(strlen($_POST['password']) < 8 || $_POST['password'] != $_POST['password_2'])
		Page\components::redirect('register.php?error=p_fail&token='.urlencode($_POST['token']));

	$query = $mysql->prepare("SELECT * FROM `users` WHERE `username` = :user OR `email` = :email");
	$query->execute(array(
		':user' => $_POST['username'],
		':email' => $core->auth->decrypt($encrypted, $iv)
	));

	if($query->rowCount() > 0)
		Page\components::redirect('register.php?error=a_fail&token='.$_POST['token']);

	$query = $mysql->prepare("SELECT * FROM `account_change` WHERE `type` = 'subuser' AND `key` = :key AND `verified` = 0");
	$query->execute(array(
		':key' => $_POST['token']
	));

	if($query->rowCount() != 1)
		Page\components::redirect('register.php?error=t_fail&token='.$_POST['token']);
	else
		$row = $query->fetch();

	$insert = $mysql->prepare("INSERT INTO `users` VALUES(NULL, NULL, :uuid, :user, :email, :pass, :permissions, :language, :time, NULL, NULL, 0, 0, 0, 0, NULL)");
	$insert->execute(array(
		':uuid' => $core->auth->gen_UUID(),
		':user' => $_POST['username'],
		':email' => $core->auth->decrypt($encrypted, $iv),
		':pass' => $core->auth->hash($_POST['password']),
		':permissions' => $row['content'],
		':language' => $core->settings->get('default_language'),
		':time' => time()
	));

	$userid = $mysql->lastInsertId();

	$query = $mysql->prepare("SELECT `subusers`, `hash` FROM `servers` WHERE `hash` = :hash");
	$query->execute(array(
		':hash' => key(json_decode($row['content'], true))
	));

	$row = $query->fetch();
	$row['subusers'] = json_decode($row['subusers'], true);
	unset($row['subusers'][$core->auth->decrypt($encrypted, $iv)]);
	$row['subusers'][$userid] = "verified";

	$query = $mysql->prepare("UPDATE `servers` SET `subusers` = :subusers WHERE `hash` = :hash");
	$query->execute(array(
		':subusers' => json_encode($row['subusers']),
		':hash' => $row['hash']
	));

	$query = $mysql->prepare("UPDATE `account_change` SET `verified` = 1 WHERE `key` = :key");
	$query->execute(array(
		':key' => $_POST['token']
	));

	Page\components::redirect('index.php?registered');

}else{

	echo $twig->render(
			'panel/register.html', array(
				'xsrf' => $core->auth->XSRF(),
				'footer' => array(
					'queries' => Database\databaseInit::getCount(),
					'seconds' => number_format((microtime(true) - $pageStartTime), 4)
				)
		));

}
?>
