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

$klein->respond(array('GET', 'POST'), '/node/files/[*]', function($request, $response, $service, $app, $klein) use($core) {

	if(!$core->user->hasPermission('files.view')) {

		$response->code(403);
		$response->body($core->twig->render('node/403.html'))->send();
		$klein->skipRemaining();

	}

});

$klein->respond('GET', '/node/files', function($request, $response, $service) use($core) {

	$response->body($core->twig->render('node/files/index.html', array(
		'server' => $core->server->getData(),
		'flash' => $service->flashes()
	)))->send();

});

$klein->respond('GET', '/node/files/download/[*:file]', function($request, $response, $service) use($core) {

	if(!$core->user->hasPermission('files.download')) {

		$response->code(403);
		$response->body($core->twig->render('node/403.html'))->send();
		return;

	}

	if(!$request->param('file')) {

		$response->code(403);
		$response->body($core->twig->render('node/403.html'))->send();
		return;

	} else {

		if(!$core->gsd->avaliable($core->server->nodeData('ip'), $core->server->nodeData('gsd_listen'))) {

			$service->flash('<div class="alert alert-danger">Unable to access the server daemon to process file downloads.</div>');
			$response->redirect('/node/files')->send();
			return;

		}

		$downloadToken = $core->auth->keygen(32);

		$download = ORM::forTable('downloads')->create();
		$download->set(array(
			'server' => $core->server->getData('gsd_id'),
			'token' => $downloadToken,
			'path' => str_replace("../", "", $request->param('file'))
		));
		$download->save();

		$response->redirect("http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/download/".$downloadToken)->send();

	}

});