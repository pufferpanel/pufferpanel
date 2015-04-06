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
use \ORM;

$klein->respond('GET', '/node/index', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render('node/index.html', array(
		'server' => array_merge($core->server->getData(), array(
			'daemon_secret' => ($core->permissions->get('daemon_secret')) ? $core->permissions->get('daemon_secret') : $core->server->getData('daemon_secret'),
			'node' => $core->server->nodeData('node'),
			'console_inner' => $core->gsd->serverLog()
		)),
		'node' => $core->server->nodeData(),
		'flash' => $service->flashes()
	)))->send();

});

include 'ajax/routes.php';
include 'files/routes.php';
include 'settings/routes.php';
include 'users/routes.php';