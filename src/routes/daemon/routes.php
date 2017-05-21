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
        
        $result = $unireq->raw_body;

        $headers = $unireq->headers;
        foreach ($headers as $key => $value) {
            $response->header($key, $value);
        }
        $response->code($unireq->code)->body($result);
}

$klein->respond('/daemon/server/[:serverid]/[**:path]', function($request, $response) use ($core, $klein) {

    $oldHeaders = $request->headers();
    $header = null;
    $server = $request->param('serverid', $request->cookies()['pp_server_hash']);

    if ($server === false || $server == '' || $server === 'undefined') {
        $response->code(404);
        return;
    }

    $serverObj = ORM::forTable('servers')->selectMany(array('id', 'node'))->where('hash', $server)->findOne();

    if ($serverObj === false) {
        $response->code(401);
        return;
    }

    $nodeObj = ORM::forTable('nodes')->where('id', $serverObj->node)->findOne();

    if ($nodeObj === false) {
        $response->code(401);
        return;
    }

    if ($oldHeaders->Authorization) {
        $header = array (
            'Authorization' => $oldHeaders->Authorization
        );
    } else {
        $auth = $request->cookies()['pp_auth_token'];

        if ($auth === false || $auth == '') {
            $response->code(403);
            return;
        }

        $userObj = ORM::forTable('users')->where('session_id', $auth)->findOne();

        if ($userObj === false) {
            $response->code(401);
            return;
        }

        $pdo = ORM::get_db();
        $query = $pdo->prepare('SELECT access_token FROM oauth_access_tokens AS oat '
            . 'INNER JOIN oauth_clients AS oc ON oc.id = oat.oauthClientId '
            . 'WHERE user_id = ? AND server_id = ? AND expiretime > NOW() AND oat.scopes NOT LIKE ?');
        $query->execute(array($userObj->id, $serverObj->id, 'sftp %'));
        $data = $query->fetch(\PDO::FETCH_ASSOC);
        $bearer = null;
        if ($data === false || count($data) == 0) {
            $clientInfo = $pdo->prepare('SELECT client_id, client_secret FROM oauth_clients '
                . 'WHERE user_id = ? AND server_id = ?');
            $clientInfo->execute(array($userObj->id, $serverObj->id));
            $info = $clientInfo->fetch(\PDO::FETCH_ASSOC);
            $bearerArr = OAuthService::Get()->handleTokenCredentials($info['client_id'], $info['client_secret']);
            if (array_key_exists('error', $bearerArr)) {
                $response->code(403);
                return;
            }
            $bearer = $bearerArr['access_token'];
        } else {
            $bearer = $data['access_token'];
        }

        $header = array(
            'Authorization' => 'Bearer ' . $bearer
        );
    }

    $updatedUrl = sprintf('%s/server/%s/%s', Daemon::buildBaseUrlForNode($nodeObj->ip, $nodeObj->daemon_listen), $server, $request->param('path'));

    foreach($oldHeaders as $k => $v) {
        if ($v !== '') {
            $header[$k] = $v;
        }
    }

    try {
        handleProxy($updatedUrl, $header, $request, $response);
    } catch (\Exception $ex) {
        $response->code(500)->json(array(
                'error' => $ex->getMessage()
            ));
    } catch (\Throwable $ex) {
        $response->code(500)->json(array(
            'error' => $ex->getMessage()
        ));
    }
});
