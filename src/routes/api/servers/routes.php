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

$servers = $api->loadClass('Servers');

$klein->respond('GET', '/api/servers', function($request, $response) use ($servers) {
	$response->json($servers->getServers());
});

$klein->respond('GET', '/api/servers/[:hash]', function($request, $response) use ($servers) {

	$data = $servers->getServer($request->param('hash'));
	if(!$data) {

		$response->code(404);
		$response->json(array('message' => 'The requested server does not exist in the system.'));

	} else {
		$response->json($data);
	}

});

$klein->respond('POST', '/api/servers', function($request, $response) use ($servers) {

	if(!$request->isSecure()) {

		$response->code(403);
		$response->json(array('message' => 'Creating a new server must occur over a secure connection.'));

	}

	$json = json_decode(file_get_contents('php://input'), true);
	if(json_last_error() != "JSON_ERROR_NONE") {

		$response->code(409);
		$response->json(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

	}

	$addServer = $servers->addServer($json);
	if(is_numeric($addServer)) {

		$response->code(400);
		switch($addServer) {

			case 1:
				$response->json(array('message' => 'Missing a required parameter in your JSON.'));
				break;
			case 2:
				$response->json(array('message' => 'Unable to locate that node or user.'));
				break;
			case 3:
				$response->json(array('message' => 'Unable to connect to the server running GSD. Is it online?'));
				break;
			case 4:
				$response->json(array('message' => 'Invalid server name. Validated using ^[\w-]{4,35}$.'));
				break;
			case 5:
				$response->json(array('message' => 'Invalid server IP or Port provided, or they are not assigned to this node.'));
				break;
			case 6:
				$response->json(array('message' => 'The selected port is currently in use.'));
				break;
			case 7:
				$response->json(array('message' => 'Invalid user email provided.'));
				break;
			case 8:
				$response->json(array('message' => 'Non-numeric value provided for memory, disk, or cpu.'));
				break;
			default:
				$response->code(500);
				$response->json(array('message' => 'An unhandled error occured when trying to add the node.'));
				break;

		}

	} else {
		$response->json($addServer);
	}

});

$klein->respond('PUT', '/api/servers/[:hash]', function($request, $response) use ($servers) {

});