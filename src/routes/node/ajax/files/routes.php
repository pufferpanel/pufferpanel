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

$klein->respond('POST', '/node/ajax/files/[*]', function() use($core) {

	$core->routes = new Router\Router_Controller('Node\Ajax\Files', $core->server);
	$core->routes = $core->routes->loadClass();

});

$klein->respond('POST', '/node/ajax/files/directory', function($request, $response) use($core) {

	if(!$core->user->hasPermission('files.view')) {

		$response->code(403);
		$response->body("You are not authorized to perform this action.")->send();
		return;

	}

	$previous_directory = array();
	if(!empty($request->param('dir'))){
		$previous_directory['first'] = true;
	}

	$go_back = explode('/', ltrim(rtrim($request->param('dir'), '/'), '/'));
	if(count($go_back) > 1 && !empty($go_back[1])) {

		$previous_directory['show'] = true;
		$previous_directory['link'] = str_replace(end($go_back), "", trim(implode('/', $go_back), '/'));
		$previous_directory['link_show'] = rtrim($previous_directory['link'], "/");

	}

	if(!$core->routes->buildContents($request->params())) {

		$response->body($core->routes->retrieveLastError())->send();

	} else {

		$response->body($core->twig->render('node/files/ajax/list_dir.html', array(
			'files' => $core->routes->getFiles(),
			'folders' => $core->routes->getFolders(),
			'extensions' => array('txt', 'yml', 'log', 'conf', 'html', 'json', 'properties', 'props', 'cfg', 'lang'),
			'zip_extensions' => array('zip', 'tar.gz', 'tar', 'gz'),
			'directory' => $previous_directory,
			'header_dir' => ltrim($request->param('dir'), '/')
		)))->send();

	}

});