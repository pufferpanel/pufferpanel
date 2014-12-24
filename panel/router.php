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
		return $core->auth->isLoggedIn($request->ip(), $request->cookies()['pp_auth_token']);
	});
});

$klein->respond('!@^/(logout|index|register|password|api)', function($request, $response, $service, $app) use ($klein) {
	if(!$app->isLoggedIn) {
		$response->redirect('/index')->send();
		$klein->skipRemaining();
	}
});

include(SRC_DIR.'routes/assets/routes.php');
include(SRC_DIR.'routes/admin/routes.php');
include(SRC_DIR.'routes/ajax/routes.php');
include(SRC_DIR.'routes/base/routes.php');
include(SRC_DIR.'routes/api/routes.php');

// @TODO 404 Handler

$klein->dispatch();