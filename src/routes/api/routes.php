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

$klein->respond(array('GET', 'POST', 'PUT', 'DELETE'), '/api/[*:trailing]/?', function($request, $response) use ($core, $api) {

	$response->header('Content-Type', 'application/json');

	if($core->settings->get('use_api') != 1) {

		http_response_code(404);
		exit(json_encode(array('message' => 'This API is not enabled.')));

	}

	if($core->settings->get('https') == 1) {

		if(!isset($_SERVER['HTTPS']) || $_SERVER['HTTPS'] == "") {

			http_response_code(403);
			exit(json_encode(array('message' => 'This API can only be accessed using a secure (HTTPS) connection.')));

		}

	}

});

$klein->with('/api/users', function() use ($klein, $api) {

	$users = $api->loadClass('Users');

	$klein->respond('GET', '/?', function($request, $response) use ($users) {
		$response->header('Content-Type', 'application/json');
		return json_encode($users->getUsers(), JSON_PRETTY_PRINT);
	});

	$klein->respond('GET', '/[:uuid]', function($request, $response) use ($users) {

		$response->header('Content-Type', 'application/json');

		$data = $users->getUser($request->param('uuid'));
		if(!$data) {

			$response->code(404);
			return json_encode(array('message' => 'The requested user does not exist in the system.'));

		} else {
			return json_encode($data, JSON_PRETTY_PRINT);
		}

	});

	$klein->respond('PUT', '/[:uuid]', function($request, $response) use ($users) {

		$response->header('Content-Type', 'application/json');

		$data = json_decode(file_get_contents('php://input'), true);
		if(json_last_error() != "JSON_ERROR_NONE") {

			$response->code(409);
			return json_encode(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

		}

		if(isset($data['password']) && !$request->isSecure()) {

			$response->code(403);
			return json_encode(array('message' => 'Updating a password must occur over a secure connection.'));

		}

		if(!$users->updateUser($request->param('uuid'), $data)) {

			$response->code(400);
			return json_encode(array('message' => 'An error occured when trying to process this data.'));

		} else {
			$response->code(204);
		}

	});

	$klein->respond('POST', '/?', function($request, $response) use ($users) {

		$response->header('Content-Type', 'application/json');

		$json = json_decode(file_get_contents('php://input'), true);
		if(json_last_error() != "JSON_ERROR_NONE") {

			$response->code(409);
			return json_encode(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

		}

		if(isset($json['password']) && !$request->isSecure()) {

			$response->code(403);
			return json_encode(array('message' => 'Using a password that you define must occur over a secure connection.'));

		}

		$addUser = $users->createUser($json);
		if(is_numeric($addUser)) {

			$response->code(400);
			switch($addUser) {

				case 1:
					return json_encode(array('message' => 'Missing a required parameter in your JSON.'));
					break;
				case 2:
					return json_encode(array('message' => 'Invalid username was provided. Matching using '));
					break;
				case 3:
					return json_encode(array('message' => 'Invalid email was provided.'));
					break;
				case 4:
					return json_encode(array('message' => 'A user with that email or username already exists in the system.'));
					break;
				default:
					$response->code(500);
					return json_encode(array('message' => 'An unhandled error occured when trying to add the user.'));
					break;

			}

		} else {
			return json_encode($addUser, JSON_PRETTY_PRINT);
		}

	});

});

$klein->with('/api/servers', function() use ($klein, $api) {

	$servers = $api->loadClass('Servers');

	$klein->respond('GET', '/?', function($request, $response) use ($servers) {
		$response->header('Content-Type', 'application/json');
		return json_encode($servers->getServers(), JSON_PRETTY_PRINT);
	});

	$klein->respond('GET', '/[:hash]', function($request, $response) use ($servers) {

		$response->header('Content-Type', 'application/json');

		$data = $servers->getServer($request->param('hash'));
		if(!$data) {

			$response->code(404);
			return json_encode(array('message' => 'The requested server does not exist in the system.'));

		} else {
			return json_encode($data, JSON_PRETTY_PRINT);
		}

	});

	$klein->respond('POST', '/?', function($request, $response) use ($servers) {

		$response->header('Content-Type', 'application/json');

		if(!$request->isSecure()) {

			$response->code(403);
			return json_encode(array('message' => 'Creating a new server must occur over a secure connection.'));

		}

		$json = json_decode(file_get_contents('php://input'), true);
		if(json_last_error() != "JSON_ERROR_NONE") {

			$response->code(409);
			return json_encode(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

		}

		$addServer = $servers->addServer($json);
		if(is_numeric($addServer)) {

			$response->code(400);
			switch($addServer) {

				case 1:
					return json_encode(array('message' => 'Missing a required parameter in your JSON.'));
					break;
				case 2:
					return json_encode(array('message' => 'Unable to locate that node or user.'));
					break;
				case 3:
					return json_encode(array('message' => 'Unable to connect to the server running GSD. Is it online?'));
					break;
				case 4:
					return json_encode(array('message' => 'Invalid server name. Validated using ^[\w-]{4,35}$.'));
					break;
				case 5:
					return json_encode(array('message' => 'Invalid server IP or Port provided, or they are not assigned to this node.'));
					break;
				case 6:
					return json_encode(array('message' => 'The selected port is currently in use.'));
					break;
				case 7:
					return json_encode(array('message' => 'Invalid user email provided.'));
					break;
				case 8:
					return json_encode(array('message' => 'Non-numeric value provided for memory, disk, or cpu.'));
					break;
				default:
					$response->code(500);
					return json_encode(array('message' => 'An unhandled error occured when trying to add the node.'));
					break;

			}

		} else {
			return json_encode($addServer, JSON_PRETTY_PRINT);
		}

	});

	$klein->respond('PUT', '/[:hash]', function($request, $response) use ($servers) {

	});

});

$klein->with('/api/nodes', function() use ($klein, $api) {

	$nodes = $api->loadClass('Nodes');

	$klein->respond('GET', '/?', function($request, $response) use ($nodes) {
		$response->header('Content-Type', 'application/json');
		return json_encode($nodes->getNodes(), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);
	});

	$klein->respond('GET', '/[i:id]', function($request, $response) use ($nodes) {

		$response->header('Content-Type', 'application/json');

		$data = $nodes->getNode($request->param('id'));
		if(!$data) {

			$response->code(404);
			return json_encode(array('message' => 'The requested node does not exist in the system.'));

		} else {
			return json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);
		}

	});

	$klein->respond('POST', '/?', function ($request, $response) use ($nodes) {

		$response->header('Content-Type', 'application/json');

		$json = json_decode(file_get_contents('php://input'), true);
		if(json_last_error() != "JSON_ERROR_NONE") {

			$response->code(409);
			return json_encode(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

		}

		$addNode = $nodes->addNode($json);
		if(is_numeric($addNode)) {

			$response->code(400);
			switch($addNode) {

				case 1:
					return json_encode(array('message' => 'Missing a required parameter in your JSON.'));
					break;
				case 2:
					return json_encode(array('message' => 'Invalid node name was provided. Matching using [\w.-]{1,15}'));
					break;
				case 3:
					return json_encode(array('message' => 'Invalid node IP provided.'));
					break;
				case 4:
					return json_encode(array('message' => 'Missing or invalid IP and Port information.'));
					break;
				default:
					$response->code(500);
					return json_encode(array('message' => 'An unhandled error occured when trying to add the node.'));
					break;

			}

		} else {
			return json_encode($addNode, JSON_PRETTY_PRINT);
		}

	});

});

$klein->respond(array('GET', 'POST', 'PUT', 'DELETE'), '/api/?', function($request, $response) {

	$response->header('Content-Type', 'application/json');
	$response->code(200);
	return json_encode(array(
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
	), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);

});