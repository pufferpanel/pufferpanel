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
    \Unirest;

$klein->respond('GET', '/account', function($request, $response, $service) use ($core) {

    $response->body($core->twig->render('panel/account.html', array(
                'xsrf' => $core->auth->XSRF(),
                'flash' => $service->flashes(),
                'notify_login_s' => $core->user->getData('notify_login_s'),
                'notify_login_f' => $core->user->getData('notify_login_f')
    )))->send();
});

$klein->respond('POST', '/account/update/[:action]', function($request, $response, $service) use ($core) {

    $core->routes = new Router\Router_Controller('Account', $core->user);
    $core->routes = $core->routes->loadClass();

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning">' . $core->language->render('error.xsrf') . '</div>');
        $response->redirect('/account')->send();
        return;
    }

    // Update Email Address
    if ($request->param('action') == "email") {

        if (!$core->auth->verifyPassword($core->user->getData('email'), $request->param('password'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_password') . '</div>');
            $response->redirect('/account')->send();
            return;
        }

        if ($request->param('newemail') == $core->user->getData('email')) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.account.same_email') . '</div>');
            $response->redirect('/account')->send();
            return;
        }

        if (!filter_var($request->param('newemail'), FILTER_VALIDATE_EMAIL)) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_email') . '</div>');
            $response->redirect('/account')->send();
            return;
        }

        $account = ORM::forTable('users')->findOne($core->user->getData('id'));
        $account->email = $request->param('newemail');
        $account->save();

        $service->flash('<div class="alert alert-success">' . printf($core->language->render('success.account.email'), $request->param('newemail')) . '</div>');
        $response->redirect('/account')->send();
    }

    // Update Account Password
    else if ($request->param('action') == "password") {

        if (!$core->auth->verifyPassword($core->user->getData('email'), $request->param('p_password'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_password') . '</div>');
            $response->redirect('/account')->send();
            return;
        }

        if ($request->param('p_password_new') != $request->param('p_password_new_2')) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.password_mismatch') . '</div>');
            $response->redirect('/account')->send();
            return;
        }

        if (!$core->auth->validatePasswordRequirements($request->param('p_password_new'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.password_strength') . '</div>');
            $response->redirect('/account')->send();
            return;
        }

        if (!$core->routes->updatePassword($request->param('p_password_new'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.unhandled') . '</div>');
            $response->redirect('/account')->send();
        } else {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('success.account.password') . '</div>');
            $response->redirect('/auth/login')->send();
        }
    }

    // Update Account Notitification Preferences
    else if ($request->param('action') == "notifications") {

        if (!$core->auth->verifyPassword($core->user->getData('email'), $request->param('password'))) {

            $service->flash('<div class="alert alert-danger">' . $core->language->render('error.invalid_password') . '</div>');
            $response->redirect('/account')->send();
            return;
        }

        $account = ORM::forTable('users')->findOne($core->user->getData('id'));
        $account->notify_login_s = $request->param('e_s');
        $account->notify_login_f = $request->param('e_f');
        $account->save();

        $service->flash('<div class="alert alert-success">' . $core->language->render('success.notifications.updated') . '</div>');
        $response->redirect('/account')->send();
    }
});

$klein->respond('GET', '/[|index:index]', function($request, $response, $service) use ($core) {

    if ($core->auth->isAdmin()) {

        $servers = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->orderByDesc('active')
                ->findArray();
    } else {

        $servers = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->where('servers.active', 1)
                ->where_in('servers.id', $core->permissions->listServers())
                ->findArray();
    }

    /*
     * List Servers
     */
    $response->body($core->twig->render('panel/index.html', array(
                'servers' => $servers,
                'user' => $core->user->getData(),
                'flash' => $service->flashes()
    )))->send();
});

$klein->respond('GET', '/index/[:goto]', function($request, $response) use ($core) {

    if (!$core->server->nodeRedirect($request->param('goto'))) {

        $response->code(403)->body($core->twig->render('errors/403.html'))->send();
        return;
    } else {

        if (!ORM::forTable('servers')->where(array(
                    'hash' => $request->param('goto')
                ))->findOne()) {

            $response->body($core->twig->render('errors/installing.html'))->send();
            return;
        } else {

            $response->cookie('pp_server_hash', $request->param('goto'), 0);
            $response->redirect('/node/index')->send();
            return;
        }
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
        $response->redirect(($request->server()["HTTP_REFERER"]) ? $request->server()["HTTP_REFERER"] : '/index')->send();
    } else {

        $response->redirect(($request->server()["HTTP_REFERER"]) ? $request->server()["HTTP_REFERER"] : '/index')->send();
    }
});

$klein->respond('GET', '/bulkcmd', function($request, $response, $service) use($core) {

    if ($core->auth->isAdmin()) {
        $servers = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->orderByDesc('active')
                ->findArray();
    } else {
        $servers = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->where('servers.active', 1)
                ->where('servers.owner_id', $core->user->getData('id'))
                ->findArray();
    }

    $response->body($core->twig->render('panel/bulkcmd.html', array(
                'servers' => $servers,
                'xsrf' => $core->auth->XSRF(),
                'flash' => $service->flashes()
    )))->send();
});

$klein->respond('POST', '/bulkcmd', function($request, $response, $service) use ($core) {

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning">' . $core->language->render('error.xsrf') . '</div>');
        $response->redirect('/bulkcmd')->send();
        return;
    }

    if ($core->auth->isAdmin()) {
        $serversArray = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->orderByDesc('active')
                ->findArray();
    } else {
        $serversArray = ORM::forTable('servers')
                ->select('servers.*')->select('nodes.name', 'node_name')->select('locations.long', 'location')
                ->join('nodes', array('servers.node', '=', 'nodes.id'))
                ->join('locations', array('nodes.location', '=', 'locations.id'))
                ->where('servers.active', 1)
                ->where('servers.owner_id', $core->user->getData('id'))
                ->findArray();
    }

    if (empty($serversArray)) {
        $service->flash('<div class="alert alert-warning">You do not have any servers.</div>');
        $response->redirect('/bulkcmd')->send();
        return;
    }

    $servers = $request->param('servers');

    if (!is_array($servers) || (is_array($servers) && count($servers) < 1)) {
        $service->flash('<div class="alert alert-warning">You selected an invalid server.</div>');
        $response->redirect('/bulkcmd')->send();
        return;
    }

    $command = $request->param('command');

    if (empty($command)) {
        $service->flash('<div class="alert alert-warning">You entered an invalid command.</div>');
        $response->redirect('/bulkcmd')->send();
        return;
    }

    foreach ($servers as $server) {
        $serverToSendCommand = null;

        foreach ($serversArray as $srv) {
            if ($srv['hash'] === $server) {
                $serverToSendCommand = $srv;
            }
        }

        if ($serverToSendCommand != null) {
            $node = ORM::forTable('nodes')
                            ->select('nodes.*')
                            ->where('id', $serverToSendCommand['node'])
                            ->findArray()[0];

            try {
                $bearer = OAuthService::Get()->getPanelAccessToken();
                $header = array(
                    'Authorization' => 'Bearer ' . $bearer
                );

                $unirest = Unirest\Request::put(sprintf('https://%s:%s/server/%s/console', array(
                            $core->server->nodeData('fqdn'),
                            $core->server->nodeData('daemon_listen'),
                            $core->server->getData('hash'))), $header, array(
                            "command" => $command
                                )
                );
            } catch (\Exception $e) {
                \Tracy\Debugger::log($e);
                $service->flash('<div class="alert alert-danger">A problem happened while sending the command to the target servers.</div>');
                $response->redirect('/bulkcmd')->send();
                return;
            }
        } else {
            $service->flash('<div class="alert alert-warning">One of the servers you selected was invalid.</div>');
            $response->redirect('/bulkcmd')->send();
            return;
        }
    }

    $service->flash('<div class="alert alert-success">Your command has been sent successfully.</div>');
    $response->redirect('/bulkcmd')->send();
});

$klein->respond('GET', '/totp', function($request, $response, $service) use($core) {

    $response->body($core->twig->render('panel/totp.html', array(
                'totp' => $core->user->getData('use_totp'),
                'xsrf' => $core->auth->XSRF(),
                'flash' => $service->flashes()
    )))->send();
});

$klein->respond('POST', '/totp', function($request, $response, $service) use($core) {

    if (!$core->auth->XSRF($request->param('xsrf'))) {

        $service->flash('<div class="alert alert-warning">' . $core->language->render('error.xsrf') . '</div>');
        $response->redirect('/totp')->send();
        return;
    }

    if (!$core->auth->validateTOTP($request->param('token'), $core->user->getData('totp_secret'))) {

        $service->flash('<div class="alert alert-danger">Unable to validate that TOTP token for this account.</div>');
        $response->redirect('/totp')->send();
        return;
    }

    $user = ORM::forTable('users')->findOne($core->user->getData('id'));
    $user->use_totp = 0;
    $user->totp_secret = null;
    $user->save();

    $service->flash('<div class="alert alert-warning">TOTP has been disabled for this account.</div>');
    $response->redirect('/totp')->send();
});
