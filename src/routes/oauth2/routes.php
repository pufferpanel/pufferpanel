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

$klein->respond('POST', '/oauth2/token/request', function($req, $res) {

    $header = $req->headers()['Authorization'];

    $id = NULL;
    $secret = NULL;

    if ($header != NULL) {
        $prefix = 'Basic ';

        if (substr($prefix, 0, strlen($prefix)) == $prefix) {
            $header = trim(substr($header, strlen($prefix)));
        } else {
            $res->code(400);
            $res->json(array("error" => "invalid_request"));
            return;
        }
        $decoded = base64_decode($header, true);

        if ($decoded === false) {
            $res->code(400);
            $res->json(array("error" => "invalid_request"));
            return;
        }

        $parts = explode(":", $decoded, 2);

        if (count($parts) != 2) {   $res->code(400);
            $res->code(400);
            $res->json(array("error" => "invalid_request"));
            return;
        }

        $id = $parts[1];
        $secret = $parts[2];
    }

    $grantType = $req->param("grant_type");

    switch ($grantType) {
        case 'client_credentials': {
                $clientId = $id == NULL ? $req->param("client_id") : $id;
                $clientSecret = $secret === NULL ? $req->param("client_secret") : $secret;
                $internal = '.internal';
                $length = strlen($internal);

                if ($clientId === false || $clientSecret === false || $clientId == 'pufferpanel' || substr($clientId, 0, $length) === $internal) {
                    $res->code(400);
                    $res->json(array("error" => "invalid_request"));
                    return;
                }

                $server = OAuthService::Get();
                $response = $server->handleTokenCredentials($clientId, $clientSecret);
                if (array_key_exists("error", $response)) {
                    $res->code(400);
                } else {
                    $res->code(200);
                }

                $res->json($response);
                break;
            }
        case 'password': {
                $username = $id == NULL ? $req->param('username') : $id;
                $password = $secret == NULL ? $req->param('password') : $secret;
                
                if ($username === false || $password === false) {
                    $res->code(400);
                    $res->json(array("error" => "invalid_request"));
                    return;
                }

                $server = OAuthService::Get();
                $response = $server->handleResourceOwner($username, $password);
                if (array_key_exists("error", $response)) {
                    $res->code(400);
                } else {
                    $res->code(200);
                }

                $res->json($response);
                $res;
                break;
            }
        default: {
                $res->code(400);
                $res->json(array("error" => "unsupported_grant_type"));
                break;
            }
    }
});

$klein->respond('POST', '/oauth2/token/info', function($req, $res) {
    $authHeader = trim($req->headers()['Authorization']);
    $parsedHeader = explode(' ', $authHeader);
    if ($authHeader === '' || count($parsedHeader) != 2 || $parsedHeader[0] !== 'Bearer') {
        $res->code(401);
        $res->json(array("error" => "invalid_token"));
        return;
    }
    
    $node = \ORM::forTable('nodes')->where_equal('daemon_secret', $parsedHeader[1])->count();
    
    if ($node !== 1) {
        $res->code(401);
        $res->json(array("error" => "invalid_token"));
        return;
    }
    
    $token = $req->param('token');    
    if ($token === false || $token == null) {
        $res->code(400);
        $res->json(array("error" => "invalid_request"));
        return;
    }
    $server = OAuthService::Get();
    $response = $server->handleInfoRequest($token);
    $res->json($response);
});
