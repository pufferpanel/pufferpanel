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
    \Unirest\Request;

$klein->respond('GET', '/account', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('panel/account.html', array(
                'xsrf' => $core->auth->XSRF(),
                'flash' => $service->flashes(),
                'notify_login_s' => $core->user->getData('notify_login_s'),
                'notify_login_f' => $core->user->getData('notify_login_f')
    )));
});

$klein->respond('POST', '/account/update/[:action]', function($request, $response, $service) use ($core) {

    $core->routes = new Router\Router_Controller('Account', $core->user);
    $core->routes = $core->routes->loadClass();

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning">' . $core->language->render('error.xsrf') . '</div>');
        $response->redirect('/account');
        return;
    }

    // Update Email Address
    if ($request->param('action') == "email") {

        if (!$core->auth->verifyPassword($core->user->getData('email'), $request->param('password'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_password') . '</div>');
            $response->redirect('/account');
            return;
        }

        if ($request->param('newemail') == $core->user->getData('email')) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.account.same_email') . '</div>');
            $response->redirect('/account');
            return;
        }

        if (!filter_var($request->param('newemail'), FILTER_VALIDATE_EMAIL)) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_email') . '</div>');
            $response->redirect('/account');
            return;
        }

        $account = ORM::forTable('users')->findOne($core->user->getData('id'));
        $account->email = $request->param('newemail');
        $account->save();

        $service->flash('<div class="alert alert-success">' . printf($core->language->render('success.account.email'), $request->param('newemail')) . '</div>');
        $response->redirect('/account');
    }

    // Update Account Password
    else if ($request->param('action') == "password") {

        if (!$core->auth->verifyPassword($core->user->getData('email'), $request->param('p_password'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_password') . '</div>');
            $response->redirect('/account');
            return;
        }

        if ($request->param('p_password_new') != $request->param('p_password_new_2')) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.password_mismatch') . '</div>');
            $response->redirect('/account');
            return;
        }

        if (!$core->auth->validatePasswordRequirements($request->param('p_password_new'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.password_strength') . '</div>');
            $response->redirect('/account');
            return;
        }

        if (!$core->routes->updatePassword($request->param('p_password_new'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.unhandled') . '</div>');
            $response->redirect('/account');
        } else {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('success.account.password') . '</div>');
            $response->redirect('/auth/login');
        }
    }

    // Update Account Notitification Preferences
    else if ($request->param('action') == "notifications") {

        if (!$core->auth->verifyPassword($core->user->getData('email'), $request->param('password'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_password') . '</div>');
            $response->redirect('/account');
            return;
        }

        $account = ORM::forTable('users')->findOne($core->user->getData('id'));
        $account->notify_login_s = $request->param('e_s');
        $account->notify_login_f = $request->param('e_f');
        $account->save();

        $service->flash('<div class="alert alert-success">' . $core->language->render('success.notifications.updated') . '</div>');
        $response->redirect('/account');
    }
});

$klein->respond('GET', '/[|index:index]', function($request, $response, $service) use ($core) {

    if ($core->auth->isAdmin()) {

        $servers = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('nodes.ip', 'daemon_host')->select('nodes.daemon_listen', 'daemon_listen')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->orderByDesc('active')
                ->findArray();
    } else {

        $servers = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('nodes.ip', 'daemon_host')->select('nodes.daemon_listen', 'daemon_listen')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->where('servers.active', 1)
                ->where_in('servers.id', $core->permissions->listServers())
                ->findArray();
    }

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

            if ($unirest->code == 200) {
                $results = array_merge($results, get_object_vars($unirest->body));
            }
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

    /*
     * List Servers
     */
    $response->body($core->twig->render('panel/index.html', array(
        'servers' => $newServers,
        'user' => $core->user->getData(),
        'flash' => $service->flashes()
    )));
});

$klein->respond('GET', '/index/[:goto]', function($request, $response) use ($core) {

    if (!$core->server->nodeRedirect($request->param('goto'))) {
        $response->code(403)->body($core->twig->render('errors/403.html'));
    } else {
        $response->cookie('pp_server_hash', $request->param('goto'), 0);
        $response->redirect('/node/index');
    }
});

$klein->respond('GET', '/language/[:language]', function($request, $response) use ($core) {

    if (file_exists(APP_DIR . 'languages/' . $request->param('language') . '.json')) {

        if ($core->auth->isLoggedIn()) {

            $account = ORM::forTable('users')->findOne($core->user->getData('id'));
            $account->set(array(
                'language' => $request->param('language')
            ));
            $account->save();
        }

        $response->cookie("pp_language", $request->param('language'), time() + 2678400);
    }

    $response->redirect(($request->server()["HTTP_REFERER"]) ? $request->server()["HTTP_REFERER"] : '/index');
});

$klein->respond('GET', '/totp', function($request, $response, $service) use($core) {

    $response->body($core->twig->render('panel/totp.html', array(
                'totp' => $core->user->getData('use_totp'),
                'xsrf' => $core->auth->XSRF(),
                'flash' => $service->flashes()
    )));
});

$klein->respond('POST', '/totp', function($request, $response, $service) use($core) {

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning">' . $core->language->render('error.xsrf') . '</div>');
        $response->redirect('/totp');
        return;
    }

    if (!$core->auth->validateTOTP($request->param('token'), $core->user->getData('totp_secret'))) {

        $service->flash('<div class="alert alert-danger">Unable to validate that TOTP token for this account.</div>');
        $response->redirect('/totp');
        return;
    }

    $user = ORM::forTable('users')->findOne($core->user->getData('id'));
    $user->use_totp = 0;
    $user->totp_secret = null;
    $user->save();

    $service->flash('<div class="alert alert-warning">TOTP has been disabled for this account.</div>');
    $response->redirect('/totp');
});
