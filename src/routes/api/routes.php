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

require SRC_DIR.'core/api/initalize.php';
$api = new API\Initalize();

$klein->respond('/api/?[*]?', function($request, $response) {

	$response->header('Content-Type', 'application/json');
	if(Settings::config()->use_api != 1) {

		$response->code(404);
		$response->json(array('message' => 'This API is not enabled.'));
		return;

	}

	if(Settings::config()->https == 1 && !$request->isSecure()) {

		$response->code(403);
		$response->json(array('message' => 'This API can only be accessed using a secure (HTTPS) connection.'));
		return;

	}

});

include 'users/routes.php';
include 'servers/routes.php';
include 'nodes/routes.php';

$klein->respond('/api/?[*]?', function($request, $response) {

	if(!$response->isSent()) {

		$response->code(404);
		$response->body(json_encode(array(
			'message' => 'You have reached an invalid API endpoint.',
			'endpoints' => array(
				'/users' => array(
					'GET' => array(
						'/' => 'List all users on the system.',
						'/:uuid' => 'List information about a specific user including servers they own.'
					),
					'POST' => array(
						'/' => 'Create a user given the correct JSON structure.'
					),
					'DELETE' => array(),
					'PUT' => array(
						'/:uuid' => 'Update information for a user.'
					)
				),
				'/servers' => array(
					'GET' => array(
						'/' => 'List all servers on the system.',
						'/:hash' => 'List detailed information about a specific server.'
					),
					'POST' => array(),
					'DELETE' => array(),
					'PUT' => array()
				),
				'/nodes' => array(
					'GET' => array(
						'/' => 'List all nodes on the system.',
						'/:id' => 'List detailed information about a specific node.'
					),
					'POST' => array(
						'/' => 'Add a new node to the system.'
					),
					'DELETE' => array(),
					'PUT' => array()
				)
			)
		), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES))->send();

	}

});