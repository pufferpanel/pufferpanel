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
    
    if ($grantType != "client_credentials") {
        $res->code(400);
        $res->json(array("error" => "unsupported_grant_type"));
        $res->send();
        return;
    }
    
    $clientId = $req->param("client_id");
    $clientSecret = $req->param("client_secret");
    
    if ($clientId === false || $clientSecrete === false) {        
        $res->code(400);
        $res->json(array("error" => "invalid_request"));
        $res->send();
    }
    
    $server = OAuthService::Get();
    $response = $server->handleTokenCredentials($clientId, $clientSecret);
    if (array_key_exists("error", $response )) {
        $res->code(400);
    } else {
        $res->code(200);
    }
    
    $res->json($response);
    $res->send();
    
});

$klein->respond('/oauth2/token/info', function($req, $res) {
    
    //TODO: ADD SECURITY SO ONLY DAEMON CAN VALIDATE
    $token = $req->param('token');
    if ($token === false) {        
        $res->code(400);
        $res->json(array("error" => "invalid_request"));
        $res->send();
    }
    $server = OAuthService::Get();
    $response = $server->handleInfoRequest($token);
    $res->json($response);
    $res->send();
    
});