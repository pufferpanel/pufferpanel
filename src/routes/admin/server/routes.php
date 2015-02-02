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
use \ORM, \Tracy\Debugger, \Unirest\Request;

$klein->respond('GET', '/admin/server/add', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render('admin/server/add.html'))->send();

});

$klein->respond('GET', '/admin/server', function($request, $response, $service) use ($core) {

	$servers = ORM::forTable('servers')->select('servers.*')->select('nodes.node', 'node_name')->select('users.email', 'user_email')
		->join('users', array('servers.owner_id', '=', 'users.id'))
		->join('nodes', array('servers.node', '=', 'nodes.id'))
		->orderByDesc('active')
		->findArray();

	$response->body($core->twig->render(
		'admin/server/find.html',
		array(
			'flash' => $service->flashes(),
			'servers' => $servers
		)
	))->send();

});

$klein->respond(array('GET', 'POST'), '/admin/server/view/[i:id]/[*]?', function($request, $response, $service) use($core) {

	if(!$core->server->rebuildData($request->param('id'))) {

		if($request->method('post')) {

			$response->body('A server by that ID does not exist in the system.')->send();

		} else {

			$service->flash('<div class="alert alert-danger">A server by that ID does not exist in the system.</div>');
			$response->redirect('/admin/server')->send();

		}

		return;

	}

	if(!$core->user->rebuildData($core->server->getData('owner_id'))) {
		throw new \Exception("This error should never occur. Attempting to access a server with an unknown user id.");
	}

});

$klein->respond('GET', '/admin/server/view/[i:id]', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render(
		'admin/server/view.html',
		array(
			'flash' => $service->flashes(),
			'node' => $core->server->nodeData(),
			'decoded' => array('ips' => json_decode($core->server->nodeData('ips'), true), 'ports' => json_decode($core->server->nodeData('ports'), true)),
			'server' => $core->server->getData(),
			'user' => $core->user->getData()
		)
	))->send();

});

$klein->respond('POST', '/admin/server/view/[i:id]/delete', function() {

});

$klein->respond('POST', '/admin/server/view/[i:id]/connection', function($request, $response, $service) use($core) {

	$server_port = $request->param('server_port_'.str_replace('.', '_', $request->param('server_ip')));

	$ports = json_decode($core->server->nodeData('ports'), true);
	$ips = json_decode($core->server->nodeData('ips'), true);

	if(!array_key_exists($request->param('server_ip'), $ports)) {

		$service->flash('<div class="alert alert-danger">The selected IP does not exist on the node.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	if(!array_key_exists($server_port, $ports[$request->param('server_ip')])) {

		$service->flash('<div class="alert alert-danger">The selected port does not exist on the node for this IP.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	if($ports[$request->param('server_ip')][$server_port] == 0 && $server_port != $core->server->getData('server_port')) {

		$service->flash('<div class="alert alert-danger">The selected port is currently in use for this IP.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	try {

		$unirest = Request::put(
			"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id'),
			array(
				'X-Access-Token' => $core->server->nodeData('gsd_secret')
			),
			array(
				"game_cfg" => json_encode(array(
					"gameport" => (int) $server_port,
					"gamehost" => $request->param('server_ip')
				))
			)
		);

		if($unirest->code > 204){

			Debugger::log($e);
			$service->flash('<div class="alert alert-danger">GSD returned an error when trying to process your request. GSD said: '.$unirest->raw_body.' [HTTP/1.1 '.$unirest->code.']</div>');
			$response->redirect('/admin/server/view/'.$request->param('id'))->send();
			return;

		}

		$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
		$server->server_ip = $request->param('server_ip');
		$server->server_port = $server_port;
		$server->save();

		/*
		* Update Old
		*/
		$ports[$core->server->getData('server_ip')][$core->server->getData('server_port')] = 1;
		$ips[$core->server->getData('server_ip')]['ports_free']++;

		/*
		* Update Old
		*/
		$ports[$request->param('server_ip')][$server_port] = 0;
		$ips[$request->param('server_ip')]['ports_free']--;

		$node = ORM::forTable('nodes')->findOne($core->server->getData('node'));
		$node->ports = json_encode($ports);
		$node->ips = json_encode($ips);
		$node->save();

	} catch(\Exception $e) {

		Debugger::log($e);
		$service->flash('<div class="alert alert-danger">An error occured while trying to connect to the remote node. Please check that GSD is running and try again.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

});

$klein->respond('POST', '/admin/server/view/[i:id]/ftp', function($request, $response, $service) use($core) {

	if($request->param('ftp_pass') != $request->param('ftp_pass_2')) {

		$service->flash('<div class="alert alert-danger">The FTP passwords did not match.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id').'?tab=ftp_sett')->send();
		return;

	}

	if(!$core->auth->validatePasswordRequirements($request->param('ftp_pass'))) {

		$service->flash('<div class="alert alert-danger">The assigned password does not meet the requirements for passwords on this server.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id').'?tab=ftp_sett')->send();
		return;

	}

	$iv = $core->auth->generate_iv();
	$pass = $core->auth->encrypt($request->param('ftp_pass'), $iv);

	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->ftp_pass = $pass;
	$server->encryption_iv = $iv;
	$server->save();

	if($request->param('email_user')){

		$core->email->buildEmail('admin_new_ftppass', array(
			'PASS' => $request->param('ftp_pass'),
			'SERVER' => $core->server->getData('name')
		))->dispatch($core->user->getData('email'), Settings::config()->company_name.' - Your FTP Password was Reset');

	}

});

$klein->respond('POST', '/admin/server/view/[i:id]/reset-token', function($request, $response) use($core) {

	$secret = $core->auth->generateUniqueUUID('servers', 'gsd_secret');

	try {

		$unirest = Request::put(
			'http://'.$core->server->nodeData('ip').':'.$core->server->nodeData('gsd_listen').'/gameservers/'.$core->server->getData('gsd_id'),
			array(
				"X-Access-Token" => $core->server->nodeData('gsd_secret')
			),
			array(
				"keys" => json_encode(array(
					$core->server->getData('gsd_secret') => array(),
					$secret => array("s:ftp", "s:get", "s:power", "s:files", "s:files:get", "s:files:put", "s:query", "s:console", "s:console:send")
				))
			)
		);

		if($unirest->code != 200) {

			$response->body("Error trying to update token. GSD said: {$unirest->raw_body} [HTTP/1.1 {$unirest->code}]")->send();
			return;

		}

	} catch(\Exception $e) {

		Debugger::log($e);
		$response->body("The server management daemon is not responding, we were unable to update the GSD token.")->send();
		return;

	}

	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->gsd_secret = $secret;
	$server->save();

	$response->body($secret)->send();

});

$klein->respond('POST', '/admin/server/view/[i:id]/settings', function($request, $response, $service) use($core) {

	if(!is_numeric($request->param('alloc_mem')) || !is_numeric($request->param('cpu_limit'))) {

		$service->flash('<div class="alert alert-danger">Allocated memory and CPU limits must be an integer.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	if(!preg_match('/^[\w-]{4,35}$/', $request->param('server_name'))) {

		$service->flash('<div class="alert alert-danger">The server name did not meet server requirements. Server names must be between 4 and 35 characters and not contain any special characters.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->name = $request->param('server_name');
	$server->max_ram = $request->param('alloc_mem');
	//$server->disk_space = $request->param('alloc_disk');
	$server->cpu_limit = $request->param('cpu_limit');
	$server->save();

	/*
	* Build the Data
	*/
	try {

		$unirest = Request::put(
			"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id'),
			array(
				"X-Access-Token" => $core->server->nodeData('gsd_secret')
			),
			array(
				"variables" => json_encode(array(
					"-jar" => str_replace(".jar", "", $core->server->getData('server_jar')) . '.jar',
					"-Xmx" => $request->param('alloc_mem')."M"
				)),
				"build" => json_encode(array(
					"cpu" => (int) $request->param('cpu_limit')
				))
			)
		);

		$service->flash('<div class="alert alert-success">Server settings have been updated.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	} catch(\Exception $e) {

		Debugger::log($e);
		$service->flash('<div class="alert alert-danger">An error occured while trying to connect to the remote node. Please check that GSD is running and try again.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

});