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

$klein->respond('GET', '/admin/node', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render(
		'admin/node/list.html',
		array(
			'flash' => $service->flashes(),
			'nodes' => ORM::forTable('nodes')
				->select('nodes.*')->select('locations.long', 'l_location')
				->join('locations', array('nodes.location', '=', 'locations.short'))
				->findMany()
		)
	))->send();

});

$klein->respond('GET', '/admin/node/add', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render('admin/node/add.html'))->send();

});

$klein->respond('GET', '/admin/node/delete', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render('admin/node/delete.html'))->send();

});

$klein->respond('GET', '/admin/node/locations', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render('admin/node/locations.html'))->send();

});

$klein->respond('GET', '/admin/node/view/[i:id]', function($request, $response, $service) use ($core) {

	$node = ORM::forTable('nodes')->findOne($request->param('id'));

	if(!$node) {

		$service->flash('<div class="alert alert-danger">A node by that ID does not exist in the system.</div>');
		$response->redirect('/admin/node')->send();
		return;

	}

	$response->body($core->twig->render(
		'admin/node/view.html',
		array(
			'flash' => $service->flashes(),
			'node' => $node,
			'portlisting' => json_decode($node->ports, true),
		)
	))->send();

});