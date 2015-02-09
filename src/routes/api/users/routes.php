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

$users = $api->loadClass('Users');

$klein->respond('GET', '/api/users', function($request, $response) use ($users) {
	$response->json($users->getUsers());
});

$klein->respond('GET', '/api/users/[:uuid]', function($request, $response) use ($users) {

	$data = $users->getUser($request->param('uuid'));
	if(!$data) {

		$response->code(404);
		$response->json(array('message' => 'The requested user does not exist in the system.'));

	} else {
		$response->json($data);
	}

});

$klein->respond('PUT', '/api/users/[:uuid]', function($request, $response) use ($users) {

	$data = json_decode(file_get_contents('php://input'), true);
	if(json_last_error() != "JSON_ERROR_NONE") {

		$response->code(409);
		$response->json(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

	}

	if(isset($data['password']) && !$request->isSecure()) {

		$response->code(403);
		$response->json(array('message' => 'Updating a password must occur over a secure connection.'));

	}

	if(!$users->updateUser($request->param('uuid'), $data)) {

		$response->code(400);
		$response->json(array('message' => 'An error occured when trying to process this data.'));

	} else {
		$response->code(204);
	}

});

$klein->respond('POST', '/api/users', function($request, $response) use ($users) {

	$json = json_decode(file_get_contents('php://input'), true);
	if(json_last_error() != "JSON_ERROR_NONE") {

		$response->code(409);
		$response->json(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

	}

	if(isset($json['password']) && !$request->isSecure()) {

		$response->code(403);
		$response->json(array('message' => 'Using a password that you define must occur over a secure connection.'));

	}

	$addUser = $users->createUser($json);
	if(is_numeric($addUser)) {

		$response->code(400);
		switch($addUser) {

			case 1:
				$response->json(array('message' => 'Missing a required parameter in your JSON.'));
				break;
			case 2:
				$response->json(array('message' => 'Invalid username was provided. Matching using '));
				break;
			case 3:
				$response->json(array('message' => 'Invalid email was provided.'));
				break;
			case 4:
				$response->json(array('message' => 'A user with that email or username already exists in the system.'));
				break;
			default:
				$response->code(500);
				$response->json(array('message' => 'An unhandled error occured when trying to add the user.'));
				break;

		}

	} else {
		$response->json($addUser);
	}

});