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
	if(!$core->user->hasPermission('users.view')) {
		$response->redirect('/index.php?error=no_permission', 302)->send();
	}

	$access = json_decode($core->server->getData('subusers'), true);
	$users = array();

	if($core->settings->get('allow_subusers') == 1) {

		if(is_array($access) && !empty($access)) {

			foreach($access as $id => $status) {

				if($status != "verified") {

					$user = ORM::forTable('account_change')->select('content')->where(array('key' => $core->auth->encrypt($id, $status).".".$status, 'verified' => 0))->findOne();
					$_perms = json_decode($user->content, true);
					$users = array_merge($users, array($id => array("status" => "pending", "revoke" => urlencode($core->auth->encrypt($id, $status).".".$status), "permissions" => $_perms[$core->server->getData('hash')])));

				} else {

					$user = ORM::forTable('users')->selectMany('permissions', 'email', 'uuid')->where('id', $id)->findOne();
					$user = json_decode($user->permissions, true);
					$users = array_merge($users, array($user->email => array("status" => "verified", "id" => $user->uuid, "permissions" => $_perms[$core->server->getData('hash')])));

				}

			}

		}

	}

	/*
	* Display Page
	*/
	echo $twig->render('node/users/list.html', array(
		'users' => $users,
		'server' => $core->server->getData(),
		'allow_subusers' => $core->settings->get('allow_subusers'),
		'footer' => array(
			'seconds' => number_format((microtime(true) - $pageStartTime), 4)
		)
	));
});
