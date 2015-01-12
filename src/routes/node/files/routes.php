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

$klein->respond('GET' , '/node/files', function($request, $response, $service) use($core) {

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

$klein->respond('GET', '/node/files/edit/[*:file]', function($request, $response, $service) use ($core) {

	if(!$core->user->hasPermission('files.edit')) {

		$response->code(403);
		$response->body($core->twig->render('node/403.html'))->send();
		return;

	}

	$file = (object) pathinfo($request->param('file'));
	if(!in_array($file->extension, $core->files->editable())) {

		$service->flash('<div class="alert alert-danger">You do not have permission to edit files with that extension.</div>');
		$response->redirect('/node/files')->send();
		return;

	}

	if(in_array($file->dirname, array(".", "./", "/"))) {
		$file->dirname = "";
	}

	$file->dirname = trim($file->dirname, '/');

	try {

		$unirest = \Unirest::get(
			"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/file/".$file->dirname.$file->basename,
			array(
				"X-Access-Token" => $core->server->getData('gsd_secret')
			)
		);

		if($unirest->code != 200 || !isset($unirest->body->contents)) {

			$service->flash('<div class="alert alert-danger">An error was encountered when trying to retrieve this file for editing. [HTTP\1.1 '.$unirest->code.']</div>');
			$response->redirect('/node/files')->send();
			return;

		}

		/*
		* Display Page
		*/
		$response->body($core->twig->render('node/files/edit.html', array(
			'server' => $core->server->getData(),
			'xsrf' => $core->auth->XSRF(),
			'file' => $request->param('file'),
			'extension' => $file->extension,
			'directory' => $file->dirname,
			'contents' => $unirest->body->contents
		)))->send();

	} catch(\Exception $e) {

		\Tracy\Debugger::log($e);
		$service->flash('<div class="alert alert-danger">The daemon does not appear to be online currently. Please try again.</div>');
		$response->redirect('/node/files')->send();
		return;

	}

});