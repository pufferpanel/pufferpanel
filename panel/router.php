<?php

/*
	PufferPanel - A Minecraft Server Management Panel
	Copyright (c) 2014 PufferPanel

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

session_start();

require_once('../src/core/core.php');

$klein = new \Klein();

$klein->respond(function($request, $response, $service, $app) use ($core) {
	$app->register('isLoggedIn', function() use ($core, $request) {
		return $core->auth->isLoggedIn($request->server()['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'));
	});
});

$klein->respond('!(logout|login|register|password)', function($request, $response, $service, $app) {
	if(!$app->loggedIn) {
		$response->redirect('/login')->send();
		throw new Exception("Not logged in");
	}
});

$klein->respond('GET', '/logout', function($request, $response, $service, $app) use ($core) {
	if($app->isLoggedIn) {

		/*
		 * Expire Session and Clear Database Details
		 */
		setcookie("pp_auth_token", null, time() - 86400, '/');
		setcookie("pp_server_node", null, time() - 86400, '/');
		setcookie("pp_server_hash", null, time() - 86400, '/');

		$logout = ORM::forTable('users')->where(array('session_id' => $_COOKIE['pp_auth_token'], 'session_ip' => $_SERVER['REMOTE_ADDR']))->findOne();
		$logout->session_id = null;
		$logout->session_ip = null;
		$logout->save();

		$core->log->getUrl()->addLog(0, 1, array('auth.user_logout', 'Account logged out from ' . $_SERVER['REMOTE_ADDR'] . '.'));

	}
	$response->redirect('/index');
});

$klein->with('/account', 'account.php');
$klein->with('/password', 'password.php');

$klein->onError(function ($klein, $err) {
	if($err->getMessage() !== "Not logged in" && !$response->isSent()) {
		//fatal error occurred somewhere, logging
		error_log($err);
	}
});

$klein->dispatch();