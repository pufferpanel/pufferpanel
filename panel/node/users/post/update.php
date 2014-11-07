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

require_once('../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Components\Page::redirect($core->settings->get('master_url').'index.php?login');
	exit();

}

	if($core->settings->get('allow_subusers') != 1)
		Components\Page::redirect('../list.php?error=not_enabled');

	if(!isset($_POST['uuid'], $_POST['permissions']))
		Components\Page::redirect('../list.php');

	if($core->auth->XSRF(@$_POST['xsrf']) !== true)
		Components\Page::redirect('../list.php?id='.$_POST['uuid'].'&error');

	if(empty($_POST['permissions']))
		Components\Page::redirect('../view.php?id='.$_POST['uuid'].'&error');

	$query = ORM::forTable('users')->select('permissions')->where('uuid', $_POST['uuid'])->findOne();

	if($query === false)
		Components\Page::redirect('../list.php?error');

	$permissions = @json_decode($query->permissions, true);
	if(!is_array($permissions) || !array_key_exists($core->server->getData('hash'), $permissions))
		Components\Page::redirect('../view.php?id='.$_POST['uuid'].'&error');

	$permissions[$core->server->getData('hash')] = $_POST['permissions'];
	$query->permissions = json_encode($permissions);
	$query->save();

	Components\Page::redirect('../view.php?id='.$_POST['uuid']);

?>