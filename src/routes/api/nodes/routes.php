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

$nodes = $api->loadClass('Nodes');

$klein->respond('GET', '/api/nodes', function($request, $response) use ($nodes) {
	$response->json($nodes->getNodes());
});

$klein->respond('GET', '/api/nodes/[i:id]', function($request, $response) use ($nodes) {

	$response->header('Content-Type', 'application/json');

	$data = $nodes->getNode($request->param('id'));
	if(!$data) {

		$response->code(404);
		$response->json(array('message' => 'The requested node does not exist in the system.'));

	} else {
		$response->json($data);
	}

});

$klein->respond('POST', '/?', function ($request, $response) use ($nodes) {

	$response->header('Content-Type', 'application/json');

	$json = json_decode(file_get_contents('php://input'), true);
	if(json_last_error() != "JSON_ERROR_NONE") {

		$response->code(409);
		$response->json(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

	}

	$addNode = $nodes->addNode($json);
	if(is_numeric($addNode)) {

		$response->code(400);
		switch($addNode) {

			case 1:
				$response->json(array('message' => 'Missing a required parameter in your JSON.'));
				break;
			case 2:
				$response->json(array('message' => 'Invalid node name was provided. Matching using [\w.-]{1,15}'));
				break;
			case 3:
				$response->json(array('message' => 'Invalid node IP provided.'));
				break;
			case 4:
				$response->json(array('message' => 'Missing or invalid IP and Port information.'));
				break;
			default:
				$response->code(500);
				$response->json(array('message' => 'An unhandled error occured when trying to add the node.'));
				break;

		}

	} else {
		$response->json($addNode);
	}

});