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

$klein->respond(array('POST', 'GET'), '/node/users/[*]?', function($request, $response, $service, $app, $klein) use ($core) {

	if($core->settings->get('allow_subusers') != 1) {

		$response->code(403);
		$response->body($core->twig->render('node/403.html'))->send();
		$klein->skipRemaining();

	}

});

$klein->respond('GET', '/node/users', function($request, $response, $service) use ($core) {

	$response->body($core->twig->render('node/users/index.html', array(
		'flash' => $service->flashes(),
		'users' => $core->server->listAffiliatedUsers(),
		'server' => $core->server->getData()
	)))->send();

});

$klein->respond('GET', '/node/users/[:action]/[:id]?', function($request, $response, $service) use ($core) {

	if($request->param('action') == 'add') {

		$response->body($core->twig->render('node/users/add.html', array(
			'flash' => $service->flashes(),
			'xsrf' => $core->auth->XSRF(),
			'server' => $core->server->getData()
		)))->send();

	} else if($request->param('action') == 'edit' && $request->param('id')) {

		$user = ORM::forTable('users')->selectMany('permissions', 'email', 'uuid')->where('uuid', $request->param('id'))->findOne();

		if(!$user || empty($user->permissions) || !is_array(json_decode($user->permissions, true))) {

			$service->flash('<div class="alert alert-danger">An error occured when trying to access that subuser.</div>');
			$response->redirect('/node/users')->send();
			return;

		}

		$permissions = json_decode($user->permissions, true);
		if(!array_key_exists($core->server->getData('hash'), $permissions)) {

			$service->flash('<div class="alert alert-danger">An error occured when trying to access that subuser.</div>');
			$response->redirect('/node/users')->send();
			return;

		}

		$response->body($core->twig->render('node/users/edit.html', array(
			'flash' => $service->flashes(),
			'server' => $core->server->getData(),
			'permissions' => $core->user->twigListPermissions($permissions[$core->server->getData('hash')]['perms']),
			'user' => array('email' => $user->email, 'uuid' => $user->uuid),
			'xsrf' => $core->auth->XSRF()
		)))->send();

	} else if($request->param('action') == 'revoke' && $request->param('id')) {

		$core->routes = new Router\Router_Controller('Node\Users', $core->server);
		$core->routes = $core->routes->loadClass();

		$query = ORM::forTable('account_change')->where(array('key' => $request->param('id'), 'verified' => 0))->findOne();
		if(!$query) {

			$query = ORM::forTable('users')->where('uuid', $request->param('id'))->findOne();
			if(!$query) {

				$service->flash('<div class="alert alert-danger">Unable to locate the requested user for revoking.</div>');
				$response->redirect('/node/users')->send();
				return;

			} else {

				if(!$core->routes->revokeActiveUserPermissions($query)) {

					$service->flash('<div class="alert alert-danger">Unable to revoke permissions for this user. ('.$core->routes->retrieveLastError(false).')</div>');
					$response->redirect('/node/users')->send();
					return;

				} else {

					$service->flash('<div class="alert alert-success">Permissions have been successfully revoked for the requested user.</div>');
					$response->redirect('/node/users')->send();

				}

			}

		} else {

			if(!array_key_exists($core->server->getData('hash'), json_decode($query->content, true))) {

				$service->flash('<div class="alert alert-danger">Unable to locate correct permissions node for this user.</div>');
				$response->redirect('/node/users')->send();
				return;

			} else {

				$query->delete();

				$permissions = json_decode($core->server->getData('subusers'), true);
				unset($permissions[$request->param('id')]);

				$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
				$server->subusers = json_encode($permissions);
				$server->save();

				$service->flash('<div class="alert alert-success">Permissions have been successfully revoked for the requested user.</div>');
				$response->redirect('/node/users')->send();

			}

		}

	}

});

$klein->respond('POST', '/node/users/add', function($request, $response, $service) use ($core) {

	$core->routes = new Router\Router_Controller('Node\Users', $core->server);
	$core->routes = $core->routes->loadClass();

	if(!$core->auth->XSRF($request->param('xsrf'))) {

		$service->flash('<div class="alert alert-warning"> The XSRF token recieved was not valid. Please make sure cookies are enabled and try your request again.</div>');
		$response->redirect('/account')->send();

	}

	if(!filter_var($request->param('email'), FILTER_VALIDATE_EMAIL)) {

		$service->flash('<div class="alert alert-danger">The email provided is invalid.</div>');
		$response->redirect('/node/users/add')->send();

	}

	if($_POST['email'] == $core->user->getData('email')) {

		$service->flash('<div class="alert alert-danger">You cannot add yourself as a subuser.</div>');
		$response->redirect('/node/users/add')->send();

	}

	if(!$response->isLocked()) {

		if(!$core->routes->addSubuser($_POST)) {

			$service->flash('<div class="alert alert-danger">Something appears to have gone wrong when trying to add this subuser. Please try again.</div>');
			$response->redirect('/node/users/add')->send();
			return;

		} else {

			$service->flash('<div class="alert alert-success">Successfully added subuser.</div>');
			$response->redirect('/node/users')->send();

		}

	}

});

$klein->respond('POST', '/node/users/edit', function($request, $response, $service) use ($core) {

});