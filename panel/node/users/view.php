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

$klein->respond('*', function($request, $response) use ($core, $twig, $pageStartTime) {
	if($core->user->hasPermission('users.view') !== true)
		$response->redirect('/index.php?error=no_permission', 302)->send();

	user = ORM::forTable('users')->selectMany('permissions', 'email')->wh	ere('uuid', $_GET['id'])->findOne();

	if($user === false)
		$response->redirect('/index.php?error', 302)->send();

	if(empty($user->permissions) || !is_array(json_decode($user->permissions, true)))
		$response->redirect('/index.php?error', 302)->send();

	$permissions = json_decode($user->permissions, true);
	if(!array_key_exists($core->server->getData('hash'), $permissions))
		$response->redirect('/index.php?error', 302)->send();

	/*
	* Display Page
	*/
	echo $twig->render('node/users/view.html', array(
		'server' => $core->server->getData(),
		'permissions' => $core->user->twigListPermissions($permissions[$core->server->getData('hash')]),
		'user' => array('email' => $user->email),
		'xsrf' => $core->auth->XSRF(),
		'allow_subusers' => $core->settings->get('allow_subusers'),
		'footer' => array(
			'seconds' => number_format((microtime(true) - $pageStartTime), 4)
		)
	));
});
