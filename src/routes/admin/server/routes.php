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
use \ORM, \Exception, \Tracy\Debugger, \Unirest\Request, \PufferPanel\Core\Components\Functions;

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

$klein->respond(array('GET', 'POST'), '/admin/server/view/[i:id]/[*]?', function($request, $response, $service, $app, $klein) use($core) {

	if(!$core->server->rebuildData($request->param('id'))) {

		if($request->method('post')) {

			$response->body('A server by that ID does not exist in the system.')->send();

		} else {

			$service->flash('<div class="alert alert-danger">A server by that ID does not exist in the system.</div>');
			$response->redirect('/admin/server')->send();

		}

		$klein->skipRemaining();
		return;

	}

	if(!$core->user->rebuildData($core->server->getData('owner_id'))) {

		throw new Exception("This error should never occur. Attempting to access a server with an unknown user id.");
		$klein->skipRemaining();

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

$klein->respond('POST', '/admin/server/view/[i:id]/delete', function($request, $response, $service) use ($core) {

	// Start Transaction so if the daemon errors we can rollback changes
	ORM::get_db()->beginTransaction();

	$node = ORM::forTable('nodes')->findOne($core->server->getData('node'));

		$ports = json_decode($node->ports, true);
		$ips = json_decode($node->ips, true);
		$ports[$core->server->getData('server_ip')][$core->server->getData('server_port')]++;
		$ips[$core->server->getData('server_ip')]['ports_free']++;
		$node->ips = json_encode($ips);
		$node->ports = json_encode($ports);

	$node->save();

	ORM::forTable('subusers')->where('server', $core->server->getData('id'))->deleteMany();
	ORM::forTable('permissions')->where('server', $core->server->getData('id'))->deleteMany();
	ORM::forTable('downloads')->where('server', $core->server->getData('id'))->deleteMany();
	ORM::forTable('servers')->where('id', $core->server->getData('id'))->deleteMany();

	try {

		$unirest = Request::delete(
			"https://".$core->server->nodeData('fqdn').":".$core->server->nodeData('daemon_listen')."/server",
			array(
				'X-Access-Token' => $core->server->nodeData('daemon_secret'),
				'X-Access-Server' => $core->server->getData('hash')
			)
		);

		if($unirest->code !== 204) {
			throw new Exception('<div class="alert alert-danger">Scales returned an error when trying to process your request. Scales said: '.$unirest->raw_body.' [HTTP/1.1 '.$unirest->code.']</div>');
		}

		ORM::get_db()->commit();
		$service->flash('<div class="alert alert-success">The requested server has been deleted from PufferPanel.</div>');
		$response->redirect('/admin/server')->send();

	} catch(Exception $e) {

		Debugger::log($e);
		ORM::get_db()->rollBack();

		$service->flash('<div class="alert alert-danger">An error was encountered with the daemon while trying to delete this server from the system.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id').'?tab=delete')->send();
		return;

	}

});

$klein->respond('POST', '/admin/server/view/[i:id]/connection', function($request, $response, $service) use($core) {

	$ports = json_decode($core->server->nodeData('ports'), true);
	$ips = json_decode($core->server->nodeData('ips'), true);

	if(!array_key_exists($request->param('server_ip'), $ports)) {

		$service->flash('<div class="alert alert-danger">The selected IP does not exist on the node.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	if(!array_key_exists($request->param('server_port'), $ports[$request->param('server_ip')])) {

		$service->flash('<div class="alert alert-danger">The selected port does not exist on the node for this IP.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	if($ports[$request->param('server_ip')][$request->param('server_port')] == 0 && $request->param('server_port') != $core->server->getData('server_port')) {

		$service->flash('<div class="alert alert-danger">The selected port is currently in use for this IP.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	ORM::get_db()->beginTransaction();

	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->server_ip = $request->param('server_ip');
	$server->server_port = $request->param('server_port');
	$server->save();

	/*
	* Update Old
	*/
	$ports[$core->server->getData('server_ip')][$core->server->getData('server_port')] = 1;
	$ips[$core->server->getData('server_ip')]['ports_free']++;

	/*
	* Update Old
	*/
	$ports[$request->param('server_ip')][$request->param('server_port')] = 0;
	$ips[$request->param('server_ip')]['ports_free']--;

	$node = ORM::forTable('nodes')->findOne($core->server->getData('node'));
	$node->ports = json_encode($ports);
	$node->ips = json_encode($ips);
	$node->save();

	try {

		$unirest = Request::put(
			"https://".$core->server->nodeData('fqdn').":".$core->server->nodeData('daemon_listen')."/server",
			array(
				'X-Access-Token' => $core->server->nodeData('daemon_secret'),
				'X-Access-Server' => $core->server->getData('hash')
			),
			array(
				"json" => json_encode(array(
					"gameport" => (int) $request->param('server_port'),
					"gamehost" => $request->param('server_ip')
				)),
				"object" => null,
				"overwrite" => false
			)
		);

		if($unirest->code !== 204) {

			$service->flash('<div class="alert alert-danger">Scales returned an error when trying to process your request. Scales said: '.$unirest->raw_body.' [HTTP/1.1 '.$unirest->code.']</div>');
			$response->redirect('/admin/server/view/'.$request->param('id'))->send();
			return;

		}

		ORM::get_db()->commit();
		$service->flash('<div class="alert alert-success">The connection information for this server has been updated.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();

	} catch(Exception $e) {

		Debugger::log($e);
		ORM::get_db()->rollBack();

		$service->flash('<div class="alert alert-danger">An error occured while trying to connect to the remote node. Please check that the daemon is running and try again.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

});

$klein->respond('POST', '/admin/server/view/[i:id]/sftp', function($request, $response, $service) use($core) {

	if($request->param('sftp_pass') != $request->param('sftp_pass_2')) {

		$service->flash('<div class="alert alert-danger">The SFTP passwords did not match.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id').'?tab=sftp_sett')->send();
		return;

	}

	if(!$core->auth->validatePasswordRequirements($request->param('sftp_pass'))) {

		$service->flash('<div class="alert alert-danger">The assigned password does not meet the requirements for passwords on this server.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id').'?tab=sftp_sett')->send();
		return;

	}

	try {

		$unirest = Request::post(
			'https://'.$core->server->nodeData('fqdn').':'.$core->server->nodeData('daemon_listen').'/server/reset-password',
			array(
				"X-Access-Token" => $core->server->nodeData('daemon_secret'),
				"X-Access-Server" => $core->server->getData('hash')
			),
			array(
				"password" => $request->param('sftp_pass')
			)
		);

		if($unirest->code !== 204) {

			$service->flash('<div class="alert alert-danger">Scales returned an error when trying to process your request. Scales said: '.$unirest->raw_body.' [HTTP/1.1 '.$unirest->code.']</div>');
			$response->redirect('/admin/server/view/'.$request->param('id'))->send();
			return;

		}

	} catch(Exception $e) {

		Debugger::log($e);
		$service->flash('<div class="alert alert-danger">An error occured while trying to connect to the remote node. Please check that Scales is running and try again.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	if($request->param('email_user')) {

		$core->email->buildEmail('admin_new_ftppass', array(
			'PASS' => $request->param('sftp_pass'),
			'SERVER' => $core->server->getData('name')
		))->dispatch($core->user->getData('email'), Settings::config()->company_name.' - Your SFTP Password was Reset');

	}

	$service->flash('<div class="alert alert-success">The SFTP password for this server has been successfully reset.</div>');
	$response->redirect('/admin/server/view/'.$request->param('id').'?tab=sftp_sett')->send();

});

$klein->respond('POST', '/admin/server/view/[i:id]/reset-token', function($request, $response) use($core) {

	$secret = $core->auth->generateUniqueUUID('servers', 'daemon_secret');

	ORM::get_db()->beginTransaction();

	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->daemon_secret = $secret;
	$server->save();

	try {

		$unirest = Request::put(
			'https://'.$core->server->nodeData('fqdn').':'.$core->server->nodeData('daemon_listen')."/server",
			array(
				"X-Access-Token" => $core->server->nodeData('daemon_secret'),
				"X-Access-Server" => $core->server->getData('hash')
			),
			array(
				"json" => json_encode(array(
					$core->server->getData('daemon_secret') => array(),
					$secret => array("s:ftp", "s:get", "s:power", "s:files", "s:files:get", "s:files:put", "s:query", "s:console", "s:console:send")
				)),
				"object" => "keys",
				"overwrite" => false
			)
		);

		if($unirest->code != 204) {

			$response->body("Error trying to update token. Scales said: {$unirest->raw_body} [HTTP/1.1 {$unirest->code}]")->send();
			return;

		}

		$response->body($secret)->send();

	} catch(Exception $e) {

		Debugger::log($e);
		ORM::get_db()->rollBack();

		$response->body("The server management daemon is not responding, we were unable to update the Scales Security Token.")->send();
		return;

	}

});

$klein->respond('POST', '/admin/server/view/[i:id]/startup', function($request, $response, $service) use ($core) {

	// Build Startup Variables
	$startup_variables = [];
	if($request->param('daemon_variables') && !empty($request->param('daemon_variables'))) {

		foreach(explode(PHP_EOL, $request->param('daemon_variables')) as $line => $contents) {

			if(strpos($contents, '|')) {

				list($var, $value) = explode('|', $contents);
				$startup_variables[$var] = str_replace(array("\n", "\r"), "", $value);

			}

		}

	}
	$startup_variables['memory'] = $core->server->getData('max_ram');

	ORM::get_db()->beginTransaction();

	$orm = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$orm->set(array(
		'daemon_startup' => $request->param('daemon_startup'),
		'daemon_variables' => $request->param('daemon_variables')
	));
	$orm->save();

	try {

		Request::put(
			"https://".$core->server->nodeData('fqdn').":".$core->server->nodeData('daemon_listen')."/server",
			array(
				"X-Access-Token" => $core->server->nodeData('daemon_secret'),
				"X-Access-Server" => $core->server->getData('hash')
			),
			array(
				"json" => json_encode(array(
					"command" => $request->param('daemon_startup'),
					"variables" => $startup_variables
				)),
				"object" => "startup",
				"overwrite" => true
			)
		);

		ORM::get_db()->commit();
		$service->flash('<div class="alert alert-success">Server startup command and variables have been updated.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	} catch(Exception $e) {

		Debugger::log($e);
		ORM::get_db()->rollBack();

		$service->flash('<div class="alert alert-danger">An error occured while trying to connect to the remote node. Please check that Scales is running and try again.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

});

$klein->respond('POST', '/admin/server/view/[i:id]/settings', function($request, $response, $service) use($core) {

	if(!is_numeric($request->param('alloc_mem')) || !is_numeric($request->param('cpu_limit'))) {

		$service->flash('<div class="alert alert-danger">Allocated memory and CPU limits must be an integer.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	if(!preg_match('/^[\w -]{4,35}$/', $request->param('server_name'))) {

		$service->flash('<div class="alert alert-danger">The server name did not meet server requirements. Server names must be between 4 and 35 characters and not contain any special characters.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

	ORM::get_db()->beginTransaction();

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

		Request::put(
			"https://".$core->server->nodeData('fqdn').":".$core->server->nodeData('daemon_listen')."/server",
			array(
				"X-Access-Token" => $core->server->nodeData('daemon_secret'),
				"X-Access-Server" => $core->server->getData('hash')
			),
			array(
				"json" => json_encode(array(
					"memory" => $request->param('alloc_mem')
				)),
				"object" => "startup:variables",
				"overwrite" => false
			)
		);

		Request::put(
			"https://".$core->server->nodeData('fqdn').":".$core->server->nodeData('daemon_listen')."/server",
			array(
				"X-Access-Token" => $core->server->nodeData('daemon_secret'),
				"X-Access-Server" => $core->server->getData('hash')
			),
			array(
				"json" => json_encode(array(
					"cpu" => (int) $request->param('cpu_limit')
				)),
				"object" => "build",
				"overwrite" => false
			)
		);

		ORM::get_db()->commit();

		$service->flash('<div class="alert alert-success">Server settings have been updated.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	} catch(Exception $e) {

		Debugger::log($e);
		ORM::get_db()->rollBack();

		$service->flash('<div class="alert alert-danger">An error occured while trying to connect to the remote node. Please check that Scales is running and try again.</div>');
		$response->redirect('/admin/server/view/'.$request->param('id'))->send();
		return;

	}

});

$klein->respond('GET', '/admin/server/new', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render(
		'admin/server/new.html',
		array(
			'locations' => ORM::forTable('locations')->findMany(),
			'plugins' => ORM::forTable('plugins')->findMany(),
			'flash' => $service->flashes()
		)
	))->send();

});

$klein->respond('GET', '/admin/server/accounts/[:email]', function($request, $response) use ($core) {

	$select = ORM::forTable('users')->where_raw('email LIKE ? OR username LIKE ?', array('%'.$request->param('email').'%', '%'.$request->param('email').'%'))->findMany();

	$resp = array();
	foreach($select as $select) {

		$resp = array_merge($resp, array(array(
			'email' => $select->email,
			'username' => $select->username,
			'hash' => md5($select->email)
		)));

	}

	$response->header('Content-Type', 'application/json');
	$response->body(json_encode(array('accounts' => $resp), JSON_PRETTY_PRINT))->send();

});

$klein->respond('POST', '/admin/server/new', function($request, $response, $service) use($core) {

	ORM::get_db()->beginTransaction();

	$node = ORM::forTable('nodes')->findOne($request->param('node'));

	if(!$node) {

		$service->flash('<div class="alert alert-danger">The selected node does not exist on the system.</div>');
		$response->redirect('/admin/server/new')->send();
		return;

	}

	if(!preg_match('/^[\w -]{4,35}$/', $request->param('server_name'))) {

		$service->flash('<div class="alert alert-danger">The name provided for the server did not meet server requirements. Server names must be between 4 and 3 characters long and contain no special characters.</div>');
		$response->redirect('/admin/server/new')->send();
		return;

	}

	$ips = json_decode($node->ips, true);
	$ports = json_decode($node->ports, true);

	if(
		!array_key_exists($request->param('server_ip'), $ips) ||
		!array_key_exists($request->param('server_port'), $ports[$request->param('server_ip')]) ||
		$ports[$request->param('server_ip')][$request->param('server_port')] == 0
	) {

		$service->flash('<div class="alert alert-danger">The selected IP or Port is currently in use or not avaliable.</div>');
		$response->redirect('/admin/server/new')->send();
		return;

	}

	$user = ORM::forTable('users')->select('id')->where('email', $request->param('email'))->findOne();

	if(!$user) {

		$service->flash('<div class="alert alert-danger">The email provided does not match any account in the system.</div>');
		$response->redirect('/admin/server/new')->send();
		return;

	}

	/*
	* Validate Disk & Memory
	*/
	if(
		!is_numeric($request->param('alloc_mem')) ||
		!is_numeric($request->param('alloc_disk')) ||
		!is_numeric($request->param('cpu_limit'))
	) {

		$service->flash('<div class="alert alert-danger">Allocated memory, disk, and CPU must all be integers.</div>');
		$response->redirect('/admin/server/new')->send();
		return;

	}

	if($request->param('sftp_pass') != $request->param('sftp_pass_2') || !$core->auth->validatePasswordRequirements($request->param('sftp_pass'))) {

		$service->flash('<div class="alert alert-danger">The SFTP password provided did not match in both fields or did not meet the server password requirements.</div>');
		$response->redirect('/admin/server/new')->send();
		return;

	}

	$sftp_username = Functions::generateFTPUsername($request->param('server_name'));
	$server_hash = $core->auth->generateUniqueUUID('servers', 'hash');
	$daemon_secret = $core->auth->generateUniqueUUID('servers', 'daemon_secret');

	// Build Startup Variables
	$startup_variables = [];
	if($request->param('daemon_variables') && !empty($request->param('daemon_variables'))) {

		foreach(explode(PHP_EOL, $request->param('daemon_variables')) as $line => $contents) {

			if(strpos($contents, '|')) {

				list($var, $value) = explode('|', $contents);
				$startup_variables[$var] = str_replace(array("\n", "\r"), "", $value);

			}

		}

	}
	$startup_variables['memory'] = $core->server->getData('max_ram');

	$server = ORM::forTable('servers')->create();
	$server->set(array(
		'hash' => $server_hash,
		'daemon_secret' => $daemon_secret,
		'node' => $request->param('node'),
		'name' => $request->param('server_name'),
		'plugin' => $request->param('installable'),
		'daemon_startup' => $request->param('daemon_startup'),
		'daemon_variables' => $request->param('daemon_variables'),
		'owner_id' => $user->id,
		'max_ram' => $request->param('alloc_mem'),
		'disk_space' => $request->param('alloc_disk'),
		'cpu_limit' => $request->param('cpu_limit'),
		'date_added' => time(),
		'server_ip' => $request->param('server_ip'),
		'server_port' => $request->param('server_port'),
		'sftp_user' => $sftp_username
	));
	$server->save();

	$ips[$request->param('server_ip')]['ports_free']--;
	$ports[$request->param('server_ip')][$request->param('server_port')]--;

	$node->ips = json_encode($ips);
	$node->ports = json_encode($ports);
	$node->save();

	/*
	* Build Call
	*/
	$data = array(
		"name" => $server_hash,
		"user" => $sftp_username,
		"build" => array(
			"pack" => ($request->param('installable') == "none") ? null : $request->param('installable'),
			"disk" => array(
				"hard" => ($request->param('alloc_disk') < 32) ? 32 : (int) $request->param('alloc_disk'),
				"soft" => ($request->param('alloc_disk') > 2048) ? (int) $request->param('alloc_disk') - 1024 : 32
			),
			"cpu" => (int) $request->param('cpu_limit')
		),
		"startup" => array(
			"command" => $request->param('daemon_startup'),
			"variables" => $startup_variables
		),
		"keys" => array(
			$daemon_secret => array(
				"s:ftp",
				"s:get",
				"s:power",
				"s:files",
				"s:files:get",
				"s:files:delete",
				"s:files:put",
				"s:files:zip",
				"s:query",
				"s:console",
				"s:console:send"
			)
		),
		"gameport" => (int) $request->param('server_port'),
		"gamehost" => $request->param('server_ip'),
		"plugin" => $request->param('plugin')
	);

	try {

		$unirest = Request::post(
			'https://'.$node->fqdn.':'.$node->daemon_listen.'/server',
			array(
				'X-Access-Token' => $node->daemon_secret,
				'X-Access-Server' => "hodor"
			),
			array(
				'settings' => json_encode($data),
				'password' => $request->param('sftp_pass')
			)
		);

		if($unirest->code !== 204) {
			throw new Exception("An error occured trying to add a server. (".$unirest->raw_body.") [HTTP 1.1/".$unirest->code."]");
		}

		ORM::for_db()->commit();

	} catch(Exception $e) {

		Debugger::log($e);
		ORM::for_db()->rollBack();

		$service->flash('<div class="alert alert-danger">An error occured while trying to connect to the remote node. Please check that Scales is running and try again.</div>');
		$response->redirect('/admin/server/new')->send();
		return;

	}

	$core->email->buildEmail('admin_new_server', array(
			'NAME' => $request->param('server_name'),
			'SFTP' => $node->fqdn.':'.$node->daemon_sftp,
			'SERVER_CONN' => $node->fqdn.':'.$request->param('server_port'),
			'USER' => $sftp_username,
			'PASS' => $request->param('sftp_pass')
	))->dispatch($request->param('email'), Settings::config()->company_name.' - New Server Added');

	$service->flash('<div class="alert alert-success">Server created successfully.</div>');
	$response->redirect('/admin/server/view/'.$server->id())->send();
	return;

});

$klein->respond('POST', '/admin/server/new/node-list', function($request, $response) use($core) {

	$response->body($core->twig->render(
		'admin/server/node-list.html',
		array(
			'nodes' => ORM::forTable('nodes')->where('location', $request->param('location'))->findMany()
		)
	))->send();

});

$klein->respond('POST', '/admin/server/new/ip-list', function($request, $response) use($core) {

	$node = ORM::forTable('nodes')->findOne($request->param('node'));
	$node->ips = json_decode($node->ips, true);
	$node->ports = json_decode($node->ports, true);

	$response->body($core->twig->render(
		'admin/server/ip-list.html',
		array(
			'node' => $node
		)
	))->send();

});