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

			$data = $users->listUsers($request->param('uuid'));
			if($data === false){

				$response->code(404);
				return json_encode(array('message' => 'The requested user does not exist in the system.'));

			}else
				return json_encode($data, JSON_PRETTY_PRINT);


		}else
			return json_encode($users->listUsers(), JSON_PRETTY_PRINT);

});

$klein->respond('GET', '/servers/[:hash]?', function ($request, $response) use ($api) {

	$servers = $api->loadClass('Servers');

		if($request->param('hash')){

			$data = $servers->listServers($request->param('hash'));
			if($data === false){

				$response->code(404);
				return json_encode(array('message' => 'The requested server does not exist in the system.'));

			}else
				return json_encode($data, JSON_PRETTY_PRINT);


		}else
			return json_encode($servers->listServers(), JSON_PRETTY_PRINT);

});

$klein->respond('GET', '/nodes/[:id]?', function ($request, $response) use ($api) {

	$nodes = $api->loadClass('Nodes');

		if($request->param('id')){

			$data = $nodes->listNodes($request->param('id'));
			if($data === false){

				$response->code(404);
				return json_encode(array('message' => 'The requested node does not exist in the system.'));

			}else
				return json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);


		}else
			return json_encode($nodes->listNodes(), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);

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
			)
		)
	), JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);

});

$klein->dispatch();

?>