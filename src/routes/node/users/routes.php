<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2015 Dane Everitt

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

    if (Settings::config('allow_subusers') != 1 || !$core->permissions->has('users.view')) {

        $response->code(403);
        $response->body($core->twig->render('node/403.html'))->send();
        $klein->skipRemaining();
    }
});

$klein->respond('GET', '/node/users', function($request, $response, $service) use ($core) {

    $select = ORM::forTable('subusers')
                    ->raw_query("SELECT subusers.*, users.email, users.username, GROUP_CONCAT(permissions.permission SEPARATOR ', ') as user_permissions FROM subusers LEFT JOIN permissions ON subusers.user = permissions.user AND subusers.server = permissions.server LEFT JOIN users ON subusers.user = users.id WHERE subusers.server = :server GROUP BY subusers.id", array(
                        'server' => $core->server->getData('id')
                    ))->findArray();

    $response->body($core->twig->render('node/users/index.html', array(
                'flash' => $service->flashes(),
                'users' => $select,
                'server' => $core->server->getData(),
                'node' => $core->server->nodeData()
    )))->send();
});

$klein->respond('GET', '/node/users/[:action]/[:email]?', function($request, $response, $service) use ($core) {

    if ($request->param('action') == 'add') {

        $response->body($core->twig->render('node/users/add.html', array(
                    'flash' => $service->flashes(),
                    'xsrf' => $core->auth->XSRF(),
                    'server' => $core->server->getData(),
                    'node' => $core->server->nodeData()
        )))->send();
    } else if ($request->param('action') == 'edit' && $request->param('email')) {

        $user = ORM::forTable('subusers')
                        ->raw_query("SELECT subusers.*, users.email, users.id as user_id, GROUP_CONCAT(permissions.permission) as user_permissions FROM subusers
				LEFT JOIN users ON subusers.user = users.id
				LEFT JOIN permissions ON subusers.user = permissions.user AND subusers.server = permissions.server
				WHERE users.email = :email AND subusers.server = :server
				GROUP BY subusers.id", array(
                            "email" => $request->param('email'),
                            "server" => $core->server->getData('id')
                        ))->findOne();

        if (!$user) {

            $service->flash('<div class="alert alert-danger">An error occured when trying to access that subuser.</div>');
            $response->redirect('/node/users')->send();
            return;
        }

        $response->body($core->twig->render('node/users/edit.html', array(
                    'flash' => $service->flashes(),
                    'server' => $core->server->getData(),
                    'permissions' => array_flip(explode(',', str_replace('.', '_', $user->user_permissions))),
                    'user' => array('email' => $user->email, 'user_id' => $user->user_id),
                    'node' => $core->server->nodeData(),
                    'xsrf' => $core->auth->XSRF()
        )))->send();
    } else if ($request->param('action') == 'revoke' && $request->param('email')) {

        $core->routes = new Router\Router_Controller('Node\Users', $core->server);
        $core->routes = $core->routes->loadClass();

        $query = ORM::forTable('subusers')
                        ->raw_query("SELECT subusers.*, users.email FROM subusers
				LEFT JOIN users ON subusers.user = users.id
				WHERE users.email = :email AND subusers.server = :server", array(
                            "email" => $request->param('email'),
                            "server" => $core->server->getData('id')
                        ))->findOne();

        if (!$query) {

            $service->flash('<div class="alert alert-danger">Unable to locate the requested user for revoking.</div>');
            $response->redirect('/node/users')->send();
            return;
        }

        if (!$core->routes->revokeActiveUserPermissions($query)) {

            $service->flash('<div class="alert alert-danger">Unable to revoke permissions for this user. (' . $core->routes->retrieveLastError(false) . ')</div>');
            $response->redirect('/node/users')->send();
            return;
        } else {

            $service->flash('<div class="alert alert-success">Permissions have been successfully revoked for the requested user.</div>');
            $response->redirect('/node/users')->send();
        }
    }
});

$klein->respond('POST', '/node/users/add', function($request, $response, $service) use ($core) {

    $core->routes = new Router\Router_Controller('Node\Users', $core->server);
    $core->routes = $core->routes->loadClass();

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning"> The XSRF token recieved was not valid. Please make sure cookies are enabled and try your request again.</div>');
        $response->redirect('/node/users')->send();
    }

    if (!filter_var($request->param('email'), FILTER_VALIDATE_EMAIL)) {

        $service->flash('<div class="alert alert-danger">The email provided is invalid.</div>');
        $response->redirect('/node/users/add')->send();
    }

    if ($request->param('email') == $core->user->getData('email')) {

        $service->flash('<div class="alert alert-danger">You cannot add yourself as a subuser.</div>');
        $response->redirect('/node/users/add')->send();
    }

    if (!$response->isLocked()) {

        if (!$core->routes->addSubuser($request->paramsPost())) {

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

    $core->routes = new Router\Router_Controller('Node\Users', $core->server);
    $core->routes = $core->routes->loadClass();

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning"> The XSRF token recieved was not valid. Please make sure cookies are enabled and try your request again.</div>');
        $response->redirect('/node/users')->send();
    }

    if (!$response->isLocked()) {

        if (!$core->routes->modifySubuser($request->paramsPost())) {

            $service->flash('<div class="alert alert-danger">Something appears to have gone wrong when trying to modify this subuser. (' . $core->routes->retrieveLastError(false) . ')</div>');
            $response->redirect('/node/users')->send();
            return;
        } else {

            $service->flash('<div class="alert alert-success">Successfully modified subuser.</div>');
            $response->redirect('/node/users')->send();
        }
    }
});
