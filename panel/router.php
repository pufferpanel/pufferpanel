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

require '../src/core/core.php';

$klein = new \Klein\Klein();

$klein->respond(function($request, $response, $service, $app) use ($core) {
	$app->register('isLoggedIn', function() use ($core, $request) {
		return $core->auth->isLoggedIn($request->ip(), $request->cookies()['pp_auth_token'], $request->cookies()['pp_server_hash']);
	});
});

$klein->respond('!@^(/auth/|/langauge/|/api/)', function($request, $response, $service, $app) use ($klein) {

	if(!$app->isLoggedIn) {

		if(!strpos($request->pathname(), "/ajax/")) {

			$service->flash('<div class="alert alert-danger">You must be logged in to access that page.</div>');
			$response->redirect('/auth/login')->send();

		} else {

			$response->code(403);
			$response->body('Not Authenticated.')->send();

		}

		$klein->skipRemaining();

	}

});

$klein->respond('@^/auth/', function($request, $response, $service, $app) use ($klein) {

	if($app->isLoggedIn && $request->pathname() != "/auth/logout") {

		$response->redirect('/index')->send();
		$klein->skipRemaining();

	}

});

$klein->respond(array('GET', 'POST'), '/node/[*]', function($request, $response, $service, $app, $klein) use($core) {

	if(is_null($request->cookies()['pp_server_hash']) || !$app->isLoggedIn()) {

		if(!$request->method('post')) {

			$service->flash('<div class="alert alert-danger">You seem to have accessed that page in an invalid manner.</div>');
			$response->redirect('/index')->send();

		} else {

			$response->body('<div class="alert alert-danger">You seem to have accessed that page in an invalid manner.</div>')->send();

		}

		$klein->skipRemaining();

	}

});

include SRC_DIR.'routes/admin/routes.php';
include SRC_DIR.'routes/ajax/routes.php';
include SRC_DIR.'routes/auth/routes.php';
include SRC_DIR.'routes/panel/routes.php';
include SRC_DIR.'routes/node/routes.php';
include SRC_DIR.'routes/api/routes.php';

$klein->respond('*', function($request, $response, $service) use ($core) {

	if(!$response->isSent()) {

		$response->code(404);
		$response->body($core->twig->render('errors/404.html'))->send();

	}

});

$klein->dispatch();