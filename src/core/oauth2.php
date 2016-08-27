<?php

/*
  PufferPanel - A Minecraft Server Management Panel
  Copyright (c) 2016 Joshua Taylor <lordralex@gmail.com>

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

use \ORM as ORM;

class OAuthService {

    private static $oauthServer;

    /**
     * @return OAuth2Server
     */
    public static function Get() {
        if (is_null(self::$oauthServer)) {
            $server = new OAuthService();
            self::$oauthServer = $server;
        }
        return self::$oauthServer;
    }

    /**
     * 
     * @param String $clientId
     * @param String $clientSecret
     * @return String Json-reply
     */
    public function handleTokenCredentials($clientId, $clientSecret) {
        $pdo = ORM::get_db();
        $query = $pdo->prepare("SELECT id, client_secret, scopes FROM oauth_clients WHERE client_id = ? AND client_secret = ?");
        $query->execute(array($clientId, $clientSecret));
        $keys = $query->fetch(\PDO::FETCH_ASSOC);
        if ($keys === false || count($keys) == 0) {
            return array("error" => $clientId);
        }
        $accessToken = base64_encode(openssl_random_pseudo_bytes(16));
        $scopes = $keys['scopes'];
        $dbId = $keys['id'];
        $expire = time() + 3600;
        $pdo->prepare("INSERT INTO oauth_access_tokens VALUES (?, ?, ?, ?)")->execute(array(
            $accessToken,
            $dbId,
            date("Y-m-d H:i:s", $expire),
            $scopes
        ));
        return array(
            "access_token" => $accessToken,
            "expires" => $expire,
            "token_type" => "bearer",
            "scope" => $scopes
        );
    }

    /**
     * @param String $token
     * @return String Json-reply
     */
    public function handleInfoRequest($token) {
        $pdo = ORM::get_db();
        $stmt = $pdo->prepare("SELECT user_id, IFNULL(hash, '*') AS server_id, oat.scopes, expiretime, oc.client_id "
                . "FROM oauth_clients AS oc "
                . "INNER JOIN oauth_access_tokens AS oat ON oat.client_id = oc.id "
                . "LEFT JOIN servers AS s ON s.id = oc.server_id "
                . "WHERE access_token = ? AND expiretime > NOW()");
        $stmt->execute(array($token));
        $data = $stmt->fetch(\PDO::FETCH_ASSOC);
        if (count($data) === 0) {
            return array("active" => false);
        }
        return array(
            "active" => strtotime($data['expiretime']) > time(),
            "scope" => $data['scopes'],
            "client_id" => $data['client_id'],
            "username" => $data['user_id'],
            "server_id" => $data['server_id']
        );
    }

    /**
     * @return string Token
     */
    public function getPanelAccessToken() {
        $clientId = 'pufferpanel';
        $clientSecret = $this->getOrGenPanelSecret();
        $pdo = ORM::get_db();
        $query = $pdo->prepare("SELECT access_token FROM oauth_access_tokens AS oat "
                . "INNER JOIN oauth_clients AS oc ON oc.id = oat.client_id "
                . "WHERE user_id = 0 AND server_id = 0 AND expiretime > NOW()");
        $query->execute();
        $data = $query->fetch(\PDO::FETCH_ASSOC);
        if ($data === false || count($data) == 0) {
            \Tracy\Debugger::log("Creating new pair with " . $clientId . " -> " . $clientSecret);
            $newToken = $this->handleTokenCredentials($clientId, $clientSecret);
            \Tracy\Debugger::log(json_encode($newToken));
            return $newToken['access_token'];
        }
        return $data['access_token'];
    }

    private function getOrGenPanelSecret() {
        $pdo = ORM::get_db();
        $query = $pdo->prepare("SELECT client_secret FROM oauth_clients WHERE client_id = ?");
        $query->execute(array('pufferpanel'));
        $data = $query->fetch(\PDO::FETCH_ASSOC);
        if ($data === false || count($data) === 0) {
            $secret = base64_encode(openssl_random_pseudo_bytes(16));
            $pdo->prepare('INSERT INTO oauth_clients VALUES (:clientId, :clientSecret, 0, 0, :scopes, :name, :desc')->execute(array(
                ':clientId' => 'pufferpanel',
                ':clientSecret' => $secret,
                ':scopes' => 'pufferadmin',
                ':name' => 'pufferpanel',
                ':desc' => 'Pufferpanel auth'
            ));
            return getPanelToken();
        }
        return $data['client_secret'];
    }

}
