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
use PufferPanel\Core, \ORM;

$klein->respond('GET', '/account', function() use ($core) {
});

$klein->respond('POST', '/account', function() use ($core) {
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

$klein->respond('GET', '/servers', function() use ($core) {
});

$klein->respond('POST', '/servers', function() use ($core) {
});