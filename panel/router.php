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

$klein = new \Klein\Klein();

$klein->respond('GET', '/assets/[**:trail]', function($request, $response) {
	$path = 'assets/' . $request->param('trail');

	if(file_exists($path)) {

		//This is a workaround to klein sending files not working
		//It is advised the web server handles the /assets/ route instead
		$ext = pathinfo($path)['extension'];
		switch ($ext) {
			case 'js': $mimetype = 'text/javascript';
				break;

			case 'css': $mimetype = 'text/css';
				break;

			default: $mimetype = 'text/plain';
				break;
		}

		$filename = basename($path);
		$response->header('Content-type', $mimetype);
		$response->header('Content-length', filesize($path));
		$response->header('Content-Disposition', 'attachment; filename="' . $filename . '"');
		$response->body(readFile($path));
		$response->send();

	} else {

		$response->code('404');

	}
});

$klein->respond(function($request, $response, $service, $app) use ($core, $twig, $pageStartTime) {
	$app->register('isLoggedIn', function() use ($core, $request) {
		return $core->auth->isLoggedIn($request->server()['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'));
	});

	$app->register('core', function() use ($core) {
		return $core;
	});

	$app->register('twig', function() use ($twig) {
		return $twig;
	});

	$app->register('pageStartTime', function() use ($pageStartTime) {
		return $pageStartTime;
	});
});

$klein->respond('/!(logout|index|register|password|api)', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	if(!$app->isLoggedIn) {
		$response->redirect('/index')->send();
	}
});

$klein->respond('/ajax/[**:trail]', function($request, $response, $service, $app) {
	$path = 'ajax/' . $request->param('trail');

	if(file_exists($path)) {

		$core = $app->core;
		$twig = $app->twig;
		include($path);

	} else {

		$response->code(404);

	}
});

$klein->respond('/admin/[**:trail]', function($request, $response, $service, $app) {
	$path = 'admin/' . $request->param('trail');

	if(file_exists($path)) {

		$core = $app->core;
		$twig = $app->twig;
		include($path);

	} else {

		$response->code(404);

	}
});

$klein->with('/account', 'account.php');
$klein->with('/password', 'password.php');
$klein->with('', 'root.php');

$klein->respond('/api?[**]', function($request, $response, $service, $app) {
	$core = $app->core;
	include('api/index.php');
});

$klein->respond('GET', '/', function($request, $response, $service, $app) {
	if($response->isSent()) {
		return;
	}

	if($app->isLoggedIn) {
		$response->redirect('/servers', 302)->send();
	} else {
		$response->redirect('/index')->send();
	}
});

$klein->onError(function ($klein, $err) {
	//If this is not an expected exception (from not logged in)
	//then we need to rethrow that error because something broke
	if($err !== "Not logged in" && !$klein->response()->isSent()) {
		\Tracy\Debugger::log($err);
	}
});

$klein->dispatch();
