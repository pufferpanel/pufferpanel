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

$klein->respond('*', function($request, $response) use ($core) {
	if(!$core->user->hasPermission('manage.rename.jar')) {
		$response->redirect('../index?error=no_permission', 302)->send();
	}

	if(!isset($_POST['jarfile']) || empty($_POST['jarfile'])) {
		$response->redirect('../settings', 302)->send();
	}

	if(!preg_match('/^([\w\d_.-]+)$/', $_POST['jarfile'])) {
		$response->redirect('../settings', 302)->send();
	}

	/*
	 * Update It
	 */
	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->server_jar = $_POST['jarfile'];
	$server->save();

	/*
	 * Update GSD Setting
	 */
	$data = http_build_query(array(
		"variables" => array(
			"-jar" => $_POST['jarfile'] . '.jar',
			"-Xmx" => $core->server->getData('max_ram') . 'M'
		)
	));
	$context_options = array(
		'http' => array(
			'method' => 'PUT',
			'header' =>
			"Content-Type: application/x-www-form-urlencoded \r\n" .
			"X-Access-Token: " . $core->server->nodeData('gsd_secret'),
			'content' => $data
		)
	);

	$context = stream_context_create($context_options);
	file_get_contents('http://' . $core->server->nodeData('ip') . ':' . $core->server->nodeData('gsd_listen') . '/gameservers/' . $core->server->getData('gsd_id'), false, $context);

	$response->redirect('../settings', 302)->send();
});
