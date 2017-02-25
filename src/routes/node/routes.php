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

$klein->respond('GET', '/node/[*]', function ($request, $response) use ($core) {
    $response->cookie('accessToken', OAuthService::Get()->getAccessToken($core->user->getData('id'), $core->server->getData('id')));
});

$klein->respond('GET', '/node/index', function ($request, $response, $service) use ($core) {

    $protocol = "wss";

    try {
        $protocol = Daemon::doesNodeUseHTTPS($core->server->nodeData()['fqdn'], $core->server->nodeData()['daemon_listen']) ? "wss" : "ws";
    } catch (\Exception $ex) {
        $service->flash('<div class="alert alert-danger">The daemon does not report it is online, functionality is reduced until it is restarted</div>');
    }

    $response->body($core->twig->render('node/index.html', array(
        'server' => array_merge($core->server->getData(), array(
            'daemon_secret' => ($core->permissions->get('daemon_secret')) ? $core->permissions->get('daemon_secret') : $core->server->getData('daemon_secret'),
            'node' => $core->server->nodeData('node')
        )),
        'node' => array_merge($core->server->nodeData(),
            array(
                'protocol' => $protocol
            )),
        'flash' => $service->flashes(),
        'user' => $core->user->getData(),
        'oauth' => OAuthService::Get()->getFor($core->user->getData('id'), $core->server->getData('id'))
    )));
});

$klein->respond('DELETE', '/node/oauth/[i:id]', function ($request, $response, $service) use ($core) {
    $service->validateParam('id')->isInt();
    $id = $request->param('id');
    if (OAuthService::Get()->hasAccess($id, $core->user->getData('id'))) {
        OAuthService::Get()->revoke($id);
        $response->code(200);
        return;
    }
    $response->code(401);
});

$klein->respond('POST', '/node/oauth', function ($request, $response, $service) use ($core) {
    $service->validateParam('oauthId')->notNull();
    $id = $request->param('oauthId');
    $name = $request->param('oauthName');
    $desc = $request->param('oauthDesc');

    if (strncmp($id, '.', 1) === 0) {
        $service->flash('<div class="alert alert-danger">OAuth ids cannot start with \'.\'</div>');
        $response->redirect('/node/index');
        return;
    }


    if ($name === 'undefined' || $name === '') {
        $name = $id;
    }

    if ($desc === 'undefined' || $desc === '') {
        $desc = $name;
    }

    $scopeDb = ORM::for_table('oauth_clients')->where(array(
            'user_id' => $core->user->getData('id'),
            'server_id' => $core->server->getData('id'),
            'client_id' => '.internal_' . $core->user->getData('id') . '_' . $core->server->getData('id'))
    )->select('scopes')->find_one();

    if (count($scopeDb) == 1) {
        $scopes = $scopeDb['scopes'];
    } else if ($core->server->getData('owner_id') == $core->user->getData('id')) {
        $scopes = OAuthService::getUserScopes();
    } else {
        $service->flash('<div class="alert alert-danger">You are not the owner of the server nor have a defined OAuth setup, please contact your server owner to refresh your permissions</div>');
        $response->redirect('/node/index');
        return;
    }

    $pdo = ORM::get_db();
    $secret = OAuthService::Get()->create($pdo, $core->user->getData('id'), $core->server->getData('id'), $id, $scopes, $name, $desc);
    $service->flash('<div class="alert alert-danger">Secret key generated: ' . $secret . '</div>');

    $response->redirect('/node/index');
});

include 'ajax/routes.php';
include 'settings/routes.php';
include 'users/routes.php';
