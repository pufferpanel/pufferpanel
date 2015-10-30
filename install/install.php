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

use \PDO as PDO;

$params = array();
parse_str(implode('&', array_splice($argv, 1)), $params);
define("BASE_DIR", __DIR__ . '/../');

if (empty($params)) {
	echo "You failed to read the docs. Go read them again\n";
	return;
}

$keyset = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*-=+[]()";
$hash = "";

for ($i = 0; $i < 48; $i++) {
	$hash .= substr($keyset, rand(0, strlen($keyset) - 1), 1);
}

$pass = "";
for ($i = 0; $i < 48; $i++) {
	$pass .= substr($keyset, rand(0, strlen($keyset) - 1), 1);
}

try {

	$fp = fopen(BASE_DIR . 'config.json', 'w');
	if ($fp === false) {
		throw new \Exception('Could not open config.json');
	}

	fwrite($fp, json_encode(array(
		'mysql' => array(
			'host' => $params['host'],
			'database' => 'pufferpanel',
			'username' => 'pufferpanel',
			'password' => $pass,
			'port' => 3306,
			'ssl' => array(
				'use' => false,
				'client-key' => '/path/to/key.pem',
				'client-cert' => '/path/to/cert.pem',
				'ca-cert' => '/path/to/ca-cert.pem'
			)
		),
		'hash' => $hash
	)));
	fclose($fp);

	if (!file_exists(BASE_DIR . 'config.json')) {
		throw new \Exception("Could not create config.json");
	}

	$mysql = new PDO('mysql:host=' . $params['host'] . ';port=' . $params['port'], $params['user'], $params['pass'], array(
		PDO::ATTR_PERSISTENT => true,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC,
		PDO::ATTR_ERRMODE => PDO::ERRMODE_EXCEPTION
	));

	$mysql->beginTransaction();
	$mysqlQueries = explode(";", file_get_contents(BASE_DIR . "install/install.sql"));
	foreach ($mysqlQueries as $query) {
		if (empty(trim($query))) {
			continue;
		}
		$mysql->exec($query);
	}

	$query = $mysql->prepare("INSERT INTO `acp_settings` (`setting_ref`, `setting_val`) VALUES
				('company_name', :cname),
				('master_url', :murl),
				('assets_url', :aurl),
				('main_website', :mwebsite),
				('transport_token', NULL),
				('transport_email', NULL),
				('transport_method','php'),
				('captcha_pub',NULL),
				('captcha_priv',NULL),
				('default_language', 'en'),
				('force_online', 0),
				('https', 0),
				('use_api', 0),
				('allow_subusers', 0)");

    $params['siteUrl'] = preg_replace("(^https?://)", "", $param['siteUrl']);
	$query->execute(array(
		':cname' => $params['companyName'],
		':murl' => 'http://' . $params['siteUrl'] . '/',
		':mwebsite' => 'http://' . $params['siteUrl'] . '/',
		':aurl' => 'http://' . $params['siteUrl'] . '/assets/'
	));

	$uuid = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0x0fff) | 0x4000, mt_rand(0, 0x3fff) | 0x8000, mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff));
	$mysql->prepare("INSERT INTO `users` VALUES(NULL, :uuid, :username, :email, :password, :language, :time, NULL, NULL, 1, 0, 1, 0, NULL)")->execute(array(
		':uuid' => $uuid,
		':username' => $params['adminName'],
		':email' => $params['adminEmail'],
		':password' => password_hash($params['adminPass'], PASSWORD_BCRYPT),
		':language' => 'en',
		':time' => time()
	));

	$hostquery = $mysql->prepare("SELECT host FROM information_schema.processlist WHERE ID=connection_id()");
	$hostquery->execute();
	$fullHost = parse_url($hostquery->fetchColumn(0));
	$host = isset($fullHost['host']) ? $fullHost['host'] : $fullHost['path'];

	try {
		$mysql->prepare("DROP USER 'pufferpanel'@:host")->execute(array(
			'host' => $host
		));
	} catch (\Exception $ex) {
		//ignoring because no user actually existed
	}

	$mysql->prepare("GRANT SELECT, UPDATE, DELETE, ALTER, INSERT ON pufferpanel.* TO 'pufferpanel'@:host IDENTIFIED BY :pass")->execute(array(
		'pass' => $pass,
		'host' => $host
	));
	echo "PufferPanel SQL user added as pufferpanel@" . $host . "\n";
	$mysql->commit();

	exit(0);
} catch (\Exception $ex) {

	echo $ex->getMessage() . "\n";
	if (isset($mysql) && $mysql->inTransaction()) {
		$mysql->rollBack();
	}
	exit(1);
}
