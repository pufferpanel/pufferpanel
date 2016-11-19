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

use ORM;

$klein->respond('GET', '/node/[*]', function($request, $response, $service) use ($core) {
    $response->cookie('accessToken', OAuthService::Get()->getAccessToken($core->user->getData('id'), $core->server->getData('id')));
});

$klein->respond('GET', '/node/index', function($request, $response, $service) use ($core) {
    $response->body($core->twig->render('node/index.html', array(
                'server' => array_merge($core->server->getData(), array(
                    'daemon_secret' => ($core->permissions->get('daemon_secret')) ? $core->permissions->get('daemon_secret') : $core->server->getData('daemon_secret'),
                    'node' => $core->server->nodeData('node')
                )),
                'node' => $core->server->nodeData(),
                'flash' => $service->flashes(),
                'user' => $core->user->getData(),
                'oauth' => OAuthService::Get()->getFor($core->user->getData('id'), $core->server->getData('id'))
    )))->send();
});

$klein->respond('DELETE', '/node/oauth/[i:id]', function($request, $response, $service) use ($core) {
    $service->validateParam('id')->isInt();
    $id = $request->param('id');
    if(OAuthService::Get()->hasAccess($id, $core->user->getData('id'))) {
        OAuthService::Get()->revoke($id);
        $response->code(200)->send();
        return;
    }
    $response->code(401)->send();
});

$klein->respond('POST', '/node/oauth', function($request, $response, $service) use ($core) {
    $service->validateParam('oauthId')->notNull();
    $id = $request->param('oauthId');
    $name = $request->param('oauthName');
    $desc = $request->param('oauthDesc');
    
    if ($name === 'undefined' || $name === '') {
        $name = $id;
    }
    
    if ($desc === 'undefined' || $desc === '') {
        $desc = $name;
    }
    
    $pdo = ORM::get_db();
    $secret = OAuthService::Get()->create($pdo, $core->user->getData('id'), $core->server->getData('id'), $id, OAuthService::Get()->getAllScopes(), $name, $desc);
    $service->flash('<div class="alert alert-danger">Secret key generated: ' . $secret . '</div>');

    $response->redirect('/node/index')->send();
    //$response->code(200)->send();
    return;
});

include 'ajax/routes.php';
include 'settings/routes.php';
include 'users/routes.php';
