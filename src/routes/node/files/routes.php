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
use \ORM;

$klein->respond(array('GET', 'POST'), '/node/files/[*]', function($request, $response, $service, $app, $klein) use($core) {

	if(!$core->permissions->has('files.view')) {

		$response->code(403);
		$response->body($core->twig->render('error/403.html'))->send();
		$klein->skipRemaining();

	}

});

$klein->respond('GET', '/node/files', function($request, $response, $service) use($core) {

	$response->body($core->twig->render('node/files/index.html', array(
		'server' => $core->server->getData(),
		'node' => $core->server->nodeData(),
		'flash' => $service->flashes()
	)))->send();

});