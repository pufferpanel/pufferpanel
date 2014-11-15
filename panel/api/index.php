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
use \ORM as ORM;
header('Content-Type: application/json');
require_once '../../src/core/core.php';
require_once '../../src/core/api/initalize.php';

$klein = new \Klein\Klein();
$base  = dirname($_SERVER['PHP_SELF']);
if(ltrim($base, '/'))
	$_SERVER['REQUEST_URI'] = substr($_SERVER['REQUEST_URI'], strlen($base));

if($core->settings->get('use_api') != 1) {

	http_response_code(404);
	echo json_encode(array('message' => 'This API is not enabled.'));
	exit();

}
// check all requests for a header
// $headers = getallheaders();
// if(array_key_exists('X-Access-Token', $headers)){
//
// 	$authenticate = ORM::forTable('api')->where('key', $headers['X-Access-Token'])->findOne();
//
// 	if(!$authenticate){
//
// 		http_response_code(401);
// 		exit();
//
// 	}
//
// }else{
//
// 	http_response_code(401);
// 	exit();
//
// }

if($core->settings->get('https') == 1) {

	if(!isset($_SERVER['HTTPS']) || $_SERVER['HTTPS'] == "") {

		http_response_code(403);
		echo json_encode(array('message' => 'This API can only be accessed using a secure (HTTPS) connection.'));
		exit();

	}

}

$api = new API\Initalize();

$klein->respond('GET', '/users/[:uuid]?', function ($request, $response) use ($api) {

	$users = $api->loadClass('Users');
	if($request->param('uuid')){

		$data = $users->getUser($request->param('uuid'));
		if($data === false){

			$response->code(404);
			return json_encode(array('message' => 'The requested user does not exist in the system.'));

		}else
			return json_encode($data, JSON_PRETTY_PRINT);


	}else
		return json_encode($users->getUsers(), JSON_PRETTY_PRINT);

});

$klein->respond('GET', '/servers/[:hash]?', function ($request, $response) use ($api) {

	$servers = $api->loadClass('Servers');

		if($request->param('hash')){

			$data = $servers->getServer($request->param('hash'));
			if($data === false){

				$response->code(404);
				return json_encode(array('message' => 'The requested server does not exist in the system.'));

			}else
				return json_encode($data, JSON_PRETTY_PRINT);


		}else
			return json_encode($servers->getServers(), JSON_PRETTY_PRINT);

});

$klein->respond('GET', '/nodes/[i:id]?', function ($request, $response) use ($api) {

	$nodes = $api->loadClass('Nodes');

		if($request->param('id')){

			$data = $nodes->getNode($request->param('id'));
			if($data === false){

				$response->code(404);
				return json_encode(array('message' => 'The requested node does not exist in the system.'));

			}else
				return json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);


		}else
			return json_encode($nodes->getNodes(), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);

});

$klein->respond('POST', '/nodes', function ($request, $response) use ($api) {

	$node = $api->loadClass('Nodes');

	if(!isset($_POST['data']) || empty($_POST['data'])){

		$response->code(400);
		return json_encode(array('message' => 'You did not pass the required POST parameters.'));

	}else{

		$decodedData = json_decode($_POST['data'], true);

		if(json_last_error() != "JSON_ERROR_NONE"){

			$response->code(409);
			return json_encode(array('message' => 'The JSON provided was invalid. ('.json_last_error().')'));

		}

		$runNodeAddition = $node->addNode($decodedData);
		if(is_numeric($decodedData)){

			$response->code(400);
			switch($decodedData){
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

		}else{

			$response->code(200);
			return json_encode($decodedData, JSON_PRETTY_PRINT);

		}

	}

});

$klein->onHttpError(function() {

	http_response_code(404);
	echo json_encode(array(
		'endpoints' => array(
			'/users' => array(
				'GET' => array(
					'/' => 'List all users on the system.',
					'/:uuid' => 'List information about a specific user including servers they own.'
				),
				'POST' => array(),
				'DELETE' => array(),
				'PUT' => array()
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

$klein->dispatch();

?>