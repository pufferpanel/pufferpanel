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

use \ORM,
    \Exception,
    \Tracy\Debugger,
    \Unirest\Request;

$klein->respond('GET', '/admin/server', function($request, $response, $service) use ($core) {

    $servers = ORM::forTable('servers')->select('servers.*')->select('nodes.name', 'node_name')->select('users.email', 'user_email')
            ->select('nodes.ip', 'daemon_host')->select('nodes.daemon_listen', 'daemon_listen')
            ->join('users', array('servers.owner_id', '=', 'users.id'))
            ->join('nodes', array('servers.node', '=', 'nodes.id'))
            ->orderByDesc('active')
            ->findArray();
    
    $bearer = OAuthService::Get()->getPanelAccessToken();
    $header = array(
        'Authorization' => 'Bearer ' . $bearer
    );

    $serverIds = array();
    $nodes = array();
    foreach($servers as $server) {
        $serverIds[] = $server['hash'];
        $nodes[] = $server["daemon_host"] . ":" . $server["daemon_listen"];
    }
    foreach($servers as $server) {
        $serverIds[] = $server['hash'];
    }

    $ids = implode(",", $serverIds);
    $nodeConnections = array_unique($nodes);
    $results = array();
    
    foreach ($nodeConnections as $nodeConnection) {
        try {
            $unirest = Request::get(vsprintf(Daemon::buildBaseUrlForNode(explode(":", $nodeConnection)[0], explode(":", $nodeConnection)[1]) .'/network?ids=%s', array(
                        $ids)),
                        $header
            );
            $results = array_merge($results, get_object_vars($unirest->body));
        } catch (\Exception $e) {
        }
    }
    
    $newServers = array();
    
    foreach($servers as $server) {
        foreach($results as $key => $value) {
            if ($server['hash'] == $key) {
                $server['connection'] = $value;
            }
        }
        $newServers[] = $server;
    }    

    $response->body($core->twig->render('admin/server/find.html', array(
                'flash' => $service->flashes(),
                'servers' => $newServers
            )))->send();
});

$klein->respond(array('GET', 'POST'), '/admin/server/view/[i:id]/[*]?', function($request, $response, $service, $app, $klein) use($core) {

    if (!$core->server->rebuildData($request->param('id'))) {

        if ($request->method('post')) {

            $response->body('A server by that ID does not exist in the system.')->send();
        } else {

            $service->flash('<div class="alert alert-danger">A server by that ID does not exist in the system.</div>');
            $response->redirect('/admin/server')->send();
        }

        $klein->skipRemaining();
        return;
    }

    if (!$core->user->rebuildData($core->server->getData('owner_id'))) {

        throw new Exception("This error should never occur. Attempting to access a server with an unknown user id.");
    }
});

$klein->respond('GET', '/admin/server/view/[i:id]', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('admin/server/view.html', array(
        'flash' => $service->flashes(),
        'node' => $core->server->nodeData(),
        'server' => $core->server->getData(),
        'user' => $core->user->getData())
    ))->send();
});

$klein->respond('POST', '/admin/server/view/[i:id]/delete/[:force]?', function($request, $response, $service) use ($core) {

    // Start Transaction so if the daemon errors we can rollback changes
    $bearer = OAuthService::Get()->getPanelAccessToken();
    ORM::get_db()->beginTransaction();

    $node = ORM::forTable('nodes')->findOne($core->server->getData('node'));

    ORM::forTable('subusers')->where('server', $core->server->getData('id'))->deleteMany();
    ORM::forTable('permissions')->where('server', $core->server->getData('id'))->deleteMany();
    ORM::forTable('downloads')->where('server', $core->server->getData('id'))->deleteMany();
    $clientIds = ORM::forTable('oauth_clients')->where('server_id', $core->server->getData('id'))->select('id')->findMany();
    foreach ($clientIds as $id) {
        ORM::forTable('oauth_access_tokens')->where('oauthClientId', $id->id)->deleteMany();
    }
    ORM::forTable('oauth_clients')->where('server_id', $core->server->getData('id'))->deleteMany();
    ORM::forTable('servers')->where('id', $core->server->getData('id'))->deleteMany();

    try {

        $header = array(
            'Authorization' => 'Bearer ' . $bearer
        );

        $updatedUrl = sprintf('%s/server/%s', Daemon::buildBaseUrlForNode($node->ip, $node->daemon_listen), $core->server->getData('hash'));

        try {
            $unirest = Request::delete($updatedUrl, $header);
        } catch (\Exception $ex) {
            throw $ex;
        }
        if ($unirest->code == 204 || $unirest->code == 200) {
            ORM::get_db()->commit();
            $service->flash('<div class="alert alert-success">The requested server has been deleted from PufferPanel.</div>');
            $response->redirect('/admin/server')->send();
        } else {
            throw new Exception('<div class="alert alert-danger">pufferd returned an error when trying to process your request. Daemon said: ' . $unirest->raw_body . ' [HTTP/1.1 ' . $unirest->code . ']</div>');
        }
    } catch (Exception $e) {

        Debugger::log($e);

        if ($request->param('force') && $request->param('force') === "force") {

            ORM::get_db()->commit();

            $service->flash('<div class="alert alert-danger">An error was encountered with the daemon while trying to delete this server from the system. <strong>Because you requested a force delete this server has been removed from the panel regardless of the reason for the error. This server and its data may still exist on the pufferd instance.</strong></div>');
            $response->redirect('/admin/server')->send();
            return;
        }

        ORM::get_db()->rollBack();
        $service->flash('<div class="alert alert-danger">An error was encountered with the daemon while trying to delete this server from the system.</div>');
        $response->redirect('/admin/server/view/' . $request->param('id') . '?tab=delete')->send();
        return;
    }
});

$klein->respond('GET', '/admin/server/new', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('admin/server/new.html', array(
        'locations' => ORM::forTable('locations')->findMany(),
        'flash' => $service->flashes())
    ))->send();
});

$klein->respond('GET', '/admin/server/accounts/[:email]', function($request, $response) use ($core) {

    $select = ORM::forTable('users')->where_raw('email LIKE ? OR username LIKE ?', array('%' . $request->param('email') . '%', '%' . $request->param('email') . '%'))->findMany();

    $resp = array();
    foreach ($select as $select) {

        $resp = array_merge($resp, array(array(
                'email' => $select->email,
                'username' => $select->username,
                'hash' => md5($select->email)
        )));
    }

    $response->header('Content-Type', 'application/json');
    $response->body(json_encode(array('accounts' => $resp), JSON_PRETTY_PRINT))->send();
});

$klein->respond('POST', '/admin/server/new', function($request, $response, $service) use($core) {

    setcookie('__temporary_pp_admin_newserver', base64_encode(json_encode($_POST)), time() + 60);
    $bearer = OAuthService::Get()->getPanelAccessToken();
    ORM::get_db()->beginTransaction();

    $node = ORM::forTable('nodes')->findOne($request->param('node'));

    if (!$node) {

        $service->flash('<div class="alert alert-danger">The selected node does not exist on the system.</div>');
        $response->redirect('/admin/server/new')->send();
        return;
    }

    if (!preg_match('/^[\w -]{4,35}$/', $request->param('server_name'))) {

        $service->flash('<div class="alert alert-danger">The name provided for the server did not meet server requirements. Server names must be between 4 and 35 characters long and contain no special characters.</div>');
        $response->redirect('/admin/server/new')->send();
        return;
    }

    $user = ORM::forTable('users')->select('id')->where('email', $request->param('email'))->findOne();

    if (!$user) {

        $service->flash('<div class="alert alert-danger">The email provided does not match any account in the system.</div>');
        $response->redirect('/admin/server/new')->send();
        return;
    }

    $server_hash = $core->auth->generateUniqueUUID('servers', 'hash');
    $daemon_secret = $core->auth->generateUniqueUUID('servers', 'daemon_secret');

    $server = ORM::forTable('servers')->create();
    $server->set(array(
        'hash' => $server_hash,
        'daemon_secret' => $daemon_secret,
        'node' => $request->param('node'),
        'name' => $request->param('server_name'),
        'owner_id' => $user->id,
        'date_added' => time(),
    ));
    $server->save();

    OAuthService::Get()->create(ORM::get_db(),
            $user->id(),
            $server->id(),
            '.internal_' . $user->id . '_' . $server->id,
            OAuthService::getUserScopes(),
            'internal_use',
            'internal_use'
    );

    /*
     * Build Call
     */
    $data = array(
        "name" => $server_hash,
        "type" => $request->param('plugin')
    );

    foreach($request->paramsPost() as $k => $value) {
        if($k === 'plugin' || $k === 'server_name' || $k === 'email') {
            continue;
        }
        $data[$k] = $value;
    }

    try {

        $header = array(
            'Authorization' => 'Bearer ' . $bearer
        );

        $unirest = Request::put(Daemon::buildBaseUrlForNode($node->ip, $node->daemon_listen) . '/server/' . $server_hash, $header, json_encode($data));

        if ($unirest->code !== 204 && $unirest->code !== 200) {
            throw new \Exception("An error occurred trying to add a server. (" . $unirest->raw_body . ") [HTTP " . $unirest->code . "]");
        }

        ORM::get_db()->commit();
    } catch (\Exception $e) {

        ORM::get_db()->rollBack();

        $service->flash('<div class="alert alert-danger">An error occurred while trying to connect to the remote node. Please check that the daemon is running and try again.<br />' . $e->getMessage() . '</div>');
        $response->redirect('/admin/server/new')->send();
        return;
    }

    $service->flash('<div class="alert alert-success">Server created successfully.</div>');
    $response->redirect('/admin/server/view/' . $server->id())->send();
    
    //have daemon install server
    try {
        Request::post(Daemon::buildBaseUrlForNode($node->ip, $node->daemon_listen) . '/server/' . $server_hash . '/install', $header, json_encode($data));
    } catch (\Exception $ex) {
    }
    return;
});

$klein->respond('POST', '/admin/server/new/node-list', function($request, $response) use($core) {

    $response->body($core->twig->render('admin/server/node-list.html', array(
                'nodes' => ORM::forTable('nodes')->where('location', $request->param('location'))->findMany()
    )))->send();
});

$klein->respond('GET', '/admin/server/new/plugins', function($request, $response) {

    $node = ORM::forTable('nodes')->findOne($request->param('node'));

    $unirest = Request::get(Daemon::buildBaseUrlForNode($node->ip, $node->daemon_listen) . '/templates');

    $response->json($unirest->body)->send();
});
