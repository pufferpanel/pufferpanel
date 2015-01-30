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

$klein->respond('GET', '/assets/[:type]/[*:file]', function($request, $response) {

	$file = APP_DIR.'assets/'.$request->param('type').'/'.$request->param('file');

	if(!file_exists($file)) {

		$response->code(404)->body("The requested asset does not exist on this server.")->send();

	}

	if($request->param('type') == 'css') {
		$response->header('Content-Type', 'text/'.$request->param('type'));
	} else if($request->param('type') == 'javascript') {
		$response->header('Content-Type', 'application/'.$request->param('type'));
	}

	$response->header('X-Content-Type-Options', 'nosniff');
	$response->header('Last-Modified', gmdate('D, d M Y H:i:s', filemtime($file)).' GMT');

	$response->body(file_get_contents($file))->send();

});