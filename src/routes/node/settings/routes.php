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

$klein->respond('GET', BASE_URL.'/node/settings', function($request, $response, $service) use ($core) {

	if(!$core->permissions->has('manage.view')) {

		$response->code(403);
		$response->body($core->twig->render('error/403.html'));
		return;

	}

	$response->body($core->twig->render('node/settings.html', array(
		'flash' => $service->flashes(),
		'xsrf' => $core->auth->XSRF(),
		'server' => array_merge($core->server->getData(), array('server_jar' => (str_replace(".jar", "", $core->server->getData('server_jar'))))),
		'node' => $core->server->nodeData()
	)));

});

$klein->respond('POST', BASE_URL.'/node/settings/password', function($request, $response) use ($core) {

	$response->body($core->auth->keygen(rand(6, 10))."-".$core->auth->keygen(rand(6, 14)));

});