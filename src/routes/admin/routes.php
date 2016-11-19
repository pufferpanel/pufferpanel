<?php
/*
PufferPanel - A Game Server Management Panel
Copyright (c) 2015 Dane Everitt

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

$klein->respond('GET', '/admin/index', function($request, $response) use ($core) {

	$response->body($core->twig->render('admin/index.html'))->send();

});

$klein->respond('GET', '/admin/passgen/[i:count]?', function($request, $response) use ($core) {

	if(!$request->param('count')) {

		$response->body($core->auth->keygen(16))->send();

	} else {

		$response->body($core->auth->keygen($request->param('count')))->send();

	}

});

include 'account/routes.php';
include 'settings/routes.php';
include 'node/routes.php';
include 'server/routes.php';