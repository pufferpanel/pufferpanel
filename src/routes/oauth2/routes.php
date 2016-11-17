<?php

/*
  PufferPanel - A Game Server Management Panel
  Copyright (c) 2015 Joshua Taylor <lordralex@gmail.com>

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

use \PufferPanel\Core\OAuthService as OAuthService;

$klein->respond('/oauth2/token/request', function($req, $res) {

    $grantType = $req->param("grant_type");

    switch ($grantType) {
        case 'client_credentials': {
                $clientId = $req->param("client_id");
                $clientSecret = $req->param("client_secret");

                if ($clientId === false || $clientSecret === false) {
                    $res->code(400);
                    $res->json(array("error" => "invalid_request"));
                    $res->send();
                }

                $server = OAuthService::Get();
                $response = $server->handleTokenCredentials($clientId, $clientSecret);
                if (array_key_exists("error", $response)) {
                    $res->code(400);
                } else {
                    $res->code(200);
                }

                $res->json($response);
                $res->send();
                break;
            }
        case 'password': {
                $username = $req->param('username');
                $password = $req->param('password');
                
                if ($username === false || $password === false) {
                    $res->code(400);
                    $res->json(array("error" => "invalid_request"));
                    $res->send();
                }

                $server = OAuthService::Get();
                $response = $server->handleResourceOwner($username, $password);
                if (array_key_exists("error", $response)) {
                    $res->code(400);
                } else {
                    $res->code(200);
                }

                $res->json($response);
                $res->send();
                break;
            }
        default: {
                $res->code(400);
                $res->json(array("error" => "unsupported_grant_type"));
                $res->send();
                break;
            }
    }
});

$klein->respond('POST', '/oauth2/token/info', function($req, $res) {
    //TODO: ADD SECURITY SO ONLY DAEMON CAN VALIDATE
    $token = $req->param('token');
    if ($token === false || $token == null) {
        $res->code(400);
        $res->json(array("error" => "invalid_request"));
        $res->send();
        return;
    }
    $server = OAuthService::Get();
    $response = $server->handleInfoRequest($token);
    $res->json($response);
    //$res->send();    
});
