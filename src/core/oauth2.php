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
     * @return OAuthService
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
        $query->execute(array($clientId, password_hash($clientSecret, PASSWORD_BCRYPT)));
        $keys = $query->fetch(\PDO::FETCH_ASSOC);
        if ($keys === false || count($keys) == 0) {
            return array("error" => $clientId);
        }
        $scopes = $keys['scopes'];
        $tokenId = $keys['id'];
        return self::generateAccessToken($tokenId, $scopes);
    }

    public function handleResourceOwner($username, $password) {
        $email = explode('|', $username)[0];
        $serverName = explode('|', $username)[1];
        $pdo = ORM::get_db();
        $userQuery = $pdo->prepare("SELECT id, password FROM users WHERE email = ?");
        $userQuery->execute(array($email));
        $user = $userQuery->fetch(\PDO::FETCH_ASSOC);
        if ($user === false || count($user) == 0) {
            return array("error" => $username);
        }

        if (!password_verify($password, $user['password'])) {
            return array("error" => $username);
        }

        $serverQuery = $pdo->prepare("SELECT s.id, hash FROM servers AS s LEFT JOIN subusers AS su ON su.server = s.id WHERE s.name = ? AND (s.owner_id = ? OR su.user = ?) LIMIT 1");
        $serverQuery->execute(array($serverName, $user['id'], $user['id']));
        $server = $serverQuery->fetch(\PDO::FETCH_ASSOC);
        if ($server === false || count($server) == 0) {
            return array("error" => $username);
        }

        $query = $pdo->prepare("SELECT id, scopes FROM oauth_clients WHERE user_id = ? AND server_id = ?");
        $query->execute(array($user['id'], $server['id']));
        $keys = $query->fetch(\PDO::FETCH_ASSOC);
        if ($keys === false || count($keys) == 0) {
            return array("error" => $username);
        }

        if(!in_array('sftp', explode(' ', $keys['scopes']))) {
            return array('error' => $username);
        }

        $tokenId = $keys['id'];
        return self::generateAccessToken($tokenId, 'sftp '. $server['hash']);
    }

    /**
     * @param String $token
     * @return String Json-reply
     */
    public function handleInfoRequest($token) {
        $pdo = ORM::get_db();
        $stmt = $pdo->prepare("SELECT user_id, IFNULL(hash, '*') AS server_id, oat.scopes, expiretime, oc.client_id "
                . "FROM oauth_clients AS oc "
                . "INNER JOIN oauth_access_tokens AS oat ON oat.oauthClientId = oc.id "
                . "LEFT JOIN servers AS s ON s.id = oc.server_id "
                . "WHERE access_token = ? AND expiretime > NOW()");
        $stmt->execute(array($token));
        $data = $stmt->fetchAll(\PDO::FETCH_ASSOC);
        if (count($data) === 0) {
            return array("active" => false);
        }
        $res = $data[0];
        return array(
            "active" => true,
            "scope" => $res['scopes'],
            "client_id" => $res['client_id'],
            "username" => $res['user_id'],
            "server_id" => $res['server_id']
        );
    }

    public function getAccessToken($userid, $serverid) {
        $pdo = ORM::get_db();
        $query = $pdo->prepare("SELECT access_token FROM oauth_access_tokens AS oat "
                . "INNER JOIN oauth_clients AS oc ON oc.id = oat.oauthClientId "
                . "WHERE user_id = ? AND server_id = ? AND expiretime > NOW() ");
        $query->execute(array($userid, $serverid));
        $data = $query->fetch(\PDO::FETCH_ASSOC);
        if ($data === false || count($data) == 0 || $data['access_token'] == '') {
            $newquery = $pdo->prepare("SELECT id, scopes FROM oauth_clients "
                    . "WHERE user_id = ? AND server_id = ? AND client_id = ?");
            $newquery->execute(array($userid, $serverid, $userid == 0 && $serverid == 0 ? 'pufferpanel' : '.internal_' . $userid . '_' . $serverid));
            $result = $newquery->fetch(\PDO::FETCH_ASSOC);
            $newToken = $this->generateAccessToken($result['id'], $result['scopes']);
            return $newToken['access_token'];
        }
        return $data['access_token'];
    }

    /**
     * @return string Token
     */
    public function getPanelAccessToken() {
        $this->getOrGenPanelSecret();
        return $this->getAccessToken(0, 0);
    }

    public function getFor($userId, $serverId) {
        $pdo = ORM::get_db();
        $query = $pdo->prepare("SELECT id, client_id, name, description FROM oauth_clients WHERE user_id = ? AND server_id = ? AND client_id NOT LIKE '\.internal%'");
        $query->execute(array($userId, $serverId));
        return $query->fetchAll(\PDO::FETCH_ASSOC);
    }

    public function hasAccess($id, $userId) {
        return 1 <= ORM::forTable('oauth_clients')
                        ->where('user_id', $userId)
                        ->where('id', $id)
                        ->count();
    }

    public function revoke($id) {
        $pdo = ORM::get_db();
        $pdo->prepare("DELETE FROM oauth_access_tokens WHERE oauthClientId = ?")->execute(array($id));
        $pdo->prepare("DELETE FROM oauth_clients WHERE id = ?")->execute(array($id));
    }

    /**
     * 
     * @return String secret key
     */
    public function create($pdo, $userId, $serverId, $id, $scopes = '', $name = '', $description = '') {
        $secret = self::generateSecret();
        $pdo->prepare('INSERT INTO oauth_clients VALUES (NULL, ?, ?, ?, ?, ?, ?, ?)')->execute(array(
            $id,
            password_hash($secret, PASSWORD_BCRYPT),
            $userId,
            $serverId,
            $scopes,
            $name,
            $description
        ));
        return $secret;
    }

    private function getOrGenPanelSecret() {
        $pdo = ORM::get_db();
        $query = $pdo->prepare("SELECT client_secret FROM oauth_clients WHERE client_id = ?");
        $query->execute(array('pufferpanel'));
        $data = $query->fetch(\PDO::FETCH_ASSOC);
        if ($data === false || count($data) === 0) {
            $this->create($pdo, null, null, 'pufferpanel', self::getUserScopes() . ' ' . self::getAdminScopes(), 'pufferpanel', 'PufferPanel Internal Auth');
            return $this->getPanelAccessToken();
        }
        return $data['client_secret'];
    }

    /**
     * @return String
     */
    public static function generateSecret() {
        return bin2hex(openssl_random_pseudo_bytes(16));
    }

    public static function getUserScopes() {
        return 'server.start server.stop server.install server.file.get server.file.put server.console server.console.send server.stats server.network sftp';
    }

    public static function getAdminScopes() {
        return 'server.create server.delete server.edit server.reload node.stop';
    }

    private static function generateAccessToken($tokenId, $scopes) {
        $accessToken = self::generateSecret();
        $pdo = ORM::get_db();
        $pdo->prepare("INSERT INTO oauth_access_tokens (access_token, oauthClientId, scopes) VALUES (?, ?, ?)")->execute(array(
            $accessToken,
            $tokenId,
            $scopes
        ));
        $pdo->prepare("UPDATE oauth_access_tokens SET expiretime = DATE_ADD(expiretime, INTERVAL 1 HOUR) WHERE access_token = ?")->execute(array(
            $accessToken
        ));
        $expireTimeQuery = $pdo->prepare("SELECT expiretime FROM oauth_access_tokens WHERE access_token = ?");
        $expireTimeQuery->execute(array($accessToken));
        $expire = strtotime($expireTimeQuery->fetch(\PDO::FETCH_ASSOC)['expiretime']);
        return array(
            "access_token" => $accessToken,
            "expires" => $expire,
            "token_type" => "bearer",
            "scope" => $scopes
        );
    }
}
