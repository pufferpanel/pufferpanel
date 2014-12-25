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
use \ORM;

$klein->respond('GET', '/account', function() use ($core) {
});

$klein->respond('POST', '/account', function() use ($core) {
});

$klein->respond('GET', '/index', function($request, $response, $service) use ($core) {

	if($core->user->getData('root_admin') == '1') {

		$servers = ORM::forTable('servers')
			->select('servers.*')->select('nodes.node', 'node_name')->select('locations.long', 'location')
			->join('nodes', array('servers.node', '=', 'nodes.id'))
			->join('locations', array('nodes.location', '=', 'locations.short'))
			->orderByDesc('active')
			->findArray();

	} else {

		$servers = ORM::forTable('servers')
			->select('servers.*')->select('nodes.node', 'node_name')->select('locations.long', 'location')
			->join('nodes', array('servers.node', '=', 'nodes.id'))
			->join('locations', array('nodes.location', '=', 'locations.short'))
			->where(array('servers.owner_id' => $core->user->getData('id'), 'servers.active' => 1))
			->where_raw('servers.owner_id = ? OR servers.hash IN(?)', array($core->user->getData('id'), join(',', $core->user->listServerPermissions())))
			->findArray();

	}

	/*
	* List Servers
	*/
	$response->body($core->twig->render('panel/index.html', array(
		'servers' => $servers,
		'flash' => $service->flashes()
	)))->send();


});

$klein->respond('GET', '/index/[:goto]', function($request, $response, $service) use ($core) {

	if(!$core->server->nodeRedirect($request->param('goto'), $core->user->getData('id'), $core->user->getData('root_admin'))) {

		$service->flash('<div class="alert alert-danger">The requested server or function does not exist, or you do not have permission to access that server or function.</div>');
		$response->redirect('/index')->send();

	} else {

		$response->cookie('pp_server_hash', $request->param('goto'), 0);
		$response->redirect('/node/index')->send();

	}

});

$klein->respond('GET', '/language/[:language]', function($request, $response, $service, $app) use ($core) {

	if(file_exists(SRC_DIR.'lang/'.$request->param('language').'.json')) {

		if($app->isLoggedIn) {

			$account = ORM::forTable('users')->findOne($core->user->getData('id'));
			$account->set(array(
				'language' => $request->param('language')
			));
			$account->save();

		}

		$response->cookie("pp_language", $request->param('language'), time() + 2678400);
		$response->redirect(($request->server()["HTTP_REFERER"]) ? $request->server()["HTTP_REFERER"] : '/servers')->send();

	} else {

		$response->redirect(($request->server()["HTTP_REFERER"]) ? $request->server()["HTTP_REFERER"] : '/servers')->send();

	}

});