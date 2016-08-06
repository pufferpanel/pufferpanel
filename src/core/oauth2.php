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
    private $pdo;

    public function __construct() {
        $this->pdo = ORM::get_db();
    }

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
        $query = $this->pdo->prepare("SELECT id, client_secret, scopes FROM oauth_clients WHERE client_id = ?");
        $query->execute(array($clientId));
        $keys = $query->fetchAll(\PDO::FETCH_ASSOC);
        if (count($keys) == 0) {
            return array("error" => $clientId);
        }

        $matchingKey = -1;
        foreach ($keys as $key => $value) {
            if (strcasecmp($value['client_secret'], $clientSecret) == 0) {
                $matchingKey = $key;
            }
        }

        if ($matchingKey == -1) {
            return $matchingKey;
        }
        $query->closeCursor();
        $accessToken = base64_encode(openssl_random_pseudo_bytes(16));
        $scopes = $keys[$matchingKey]['scopes'];
        $dbId = $keys[$matchingKey]['id'];
        $expire = time() + 3600;
        $this->pdo->prepare("INSERT INTO oauth_access_tokens VALUES (?, ?, ?, ?)")->execute(array(
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
        $pdo = $this->pdo;
        $stmt = $pdo->prepare("SELECT id, user_id, server_id, oat.scopes, expiretime, oc.client_id "
                . "FROM oauth_clients AS oc "
                . "INNER JOIN oauth_access_tokens AS oat ON oat.client_id = oc.id "
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
        $pdo = $this->pdo;
        $query = $pdo->prepare("SELECT access_token FROM oauth_access_tokens AS oat "
                . "INNER JOIN oauth_clients AS oc ON oc.id = oat.client_od "
                . "WHERE user_id = 0 AND server_id = 0 AND expiretime > NOW()");
        $query->execute();
        $data = $query->fetch(\PDO::FETCH_ASSOC);
        if (count($data) == 0) {
            return $this->handleTokenCredentials($clientId, $clientSecret)['access_token'];
        }
        return $data['access_token'];
    }

    private function getOrGenPanelSecret() {
        $pdo = $this->pdo;
        $query = $pdo->prepare("SELECT client_secret FROM oauth_clients WHERE client_id = ?");
        $query->execute(array('pufferpanel'));
        $data = $query->fetch(\PDO::FETCH_ASSOC);
        if (count($data) === 0) {
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
