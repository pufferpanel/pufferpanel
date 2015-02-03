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

$klein->respond('POST', '/ajax/status', function($request, $response) use ($core){

	if(!$core->auth->isLoggedIn()) {
		$response->body('#FF9900')->send();
	} else {

		if($request->param('server')) {

			$status = ORM::forTable('servers')
				->select('servers.gsd_id')->select('nodes.ip')->select('nodes.gsd_secret')->select('nodes.gsd_listen')
				->join('nodes', array('servers.node', '=', 'nodes.id'))
				->where('servers.id', $request->param('server'))
				->findOne();

			if(!$status) {
				$response->body('#FF9900')->send();
				return;
			}

			if(!$core->gsd->check_status($status->ip, $status->gsd_listen, $status->gsd_id, $status->gsd_secret)) {
				$response->body('#E33200')->send();
			} else {
				$response->body('#53B30C')->send();
			}

		} else {
			$response->body('#FF9900')->send();
		}

	}

});

$klein->respond('POST', '/ajax/status/node', function($request, $response) use ($core){

	if(!$core->auth->isLoggedIn()) {
		$response->body('#FF9900')->send();
	} else {

		if($request->param('node')) {

			$status = ORM::forTable('nodes')->findOne($request->param('node'));

			if(!$status) {
				$response->body('#FF9900')->send();
				return;
			}

			if(!$core->gsd->avaliable($status->ip, $status->gsd_listen, 1)) {
				$response->body('#E33200')->send();
			} else {
				$response->body('#53B30C')->send();
			}

		} else {
			$response->body('#FF9900')->send();
		}

	}

});

include 'account/routes.php';