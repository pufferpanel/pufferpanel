<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2016 Joshua Taylor <lordralex@ae97.net>

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

function handleProxy($url, $header, $request, $response) {
         $unireq = null;

        switch ($request->method()) {
            default:
            case 'GET': {
                    $unireq = Unirest\Request::get($url, $header);
                    break;
                }
            case 'POST': {
                    $unireq = Unirest\Request::post($url, $header, $request->body());
                    break;
                }
            case 'DELETE': {
                    $unireq = Unirest\Request::delete($url, $header);
                    break;
                }
            case 'PUT': {
                    $unireq = Unirest\Request::put($url, $header, $request->body());
                    break;
                }
        }
        
        $result = $unireq->body;

        $response->code($unireq->code)->json($result);   
}

$klein->respond('/daemon/[**:path]', function($request, $response) use ($core, $klein) {

    $server = $request->cookies()['pp_server_hash'];
    $auth = $request->cookies()['pp_auth_token'];

    if ($server === false || $auth === false) {
        $response->code(401)->send();
        return;
    }

    $serverObj = ORM::forTable('servers')->selectMany(array('id', 'node'))->where('hash', $server)->findOne();
    $userObj = ORM::forTable('users')->where('session_id', $auth)->findOne();

    if ($serverObj === false || $userObj === false) {
        $response->code(402)->send();
        return;
    }

    $nodeObj = ORM::forTable('nodes')->where('id', $serverObj->node)->findOne();

    $pdo = ORM::get_db();
    $query = $pdo->prepare("SELECT access_token FROM oauth_access_tokens AS oat "
            . "INNER JOIN oauth_clients AS oc ON oc.id = oat.oauthClientId "
            . "WHERE user_id = ? AND server_id = ? AND expiretime > NOW()");
    $query->execute(array($userObj->id, $serverObj->id));
    $data = $query->fetch(\PDO::FETCH_ASSOC);
    $bearer = null;
    if ($data === false || count($data) == 0) {
        $clientInfo = $pdo->prepare("SELECT client_id, client_secret FROM oauth_clients "
                . "WHERE user_id = ? AND server_id = ?");
        $clientInfo->execute(array($userObj->id, $serverObj->id));
        $info = $clientInfo->fetch(\PDO::FETCH_ASSOC);
        $bearerArr = OAuthService::Get()->handleTokenCredentials($info['client_id'], $info['client_secret']);
        if (array_key_exists('error', $bearerArr)) {
            $response->code(403)->send();
            return;
        }
        $bearer = $bearerArr['access_token'];
    } else {
        $bearer = $data['access_token'];
    }

    $header = array(
        'Authorization' => 'Bearer ' . $bearer
    );

    $updatedUrl = sprintf("%s/%s", Daemon::buildBaseUrlForNode($nodeObj->ip, $nodeObj->daemon_listen), $request->param('path'));

    try {
        handleProxy($updatedUrl, $header, $request, $response);
    } catch (\Exception $ex) {
        $response->code(500)->json(array(
                'error' => $ex->getMessage()
            ))->send();
    }
    $klein->skipRemaining();
});
