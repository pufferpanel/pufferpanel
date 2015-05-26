<?php

/*
  PufferPanel - A Minecraft Server Management Panel
  Copyright (c) 2014 PufferPanel

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

use \PDO;

try {

	$mysql = new PDO('mysql:host=' . $params['mysqlHost'], $params['mysqlUser'], $params['mysqlPass'], array(
		PDO::ATTR_PERSISTENT => true,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
	));

	$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
	$mysql->beginTransaction();

	$query = $mysql->prepare("INSERT INTO `acp_settings` (`setting_ref`, `setting_val`) VALUES
					('company_name', :cname),
					('master_url', :murl),
					('assets_url', :aurl),
					('main_website', :mwebsite),
					('postmark_api_key', NULL),
					('mandrill_api_key', NULL),
					('mailgun_api_key', NULL),
					('sendgrid_api_key', NULL),
					('sendmail_email', NULL),
					('sendmail_method','php'),
					('captcha_pub',NULL),
					('captcha_priv',NULL),
					('default_language', 'en'),
					('force_online', 0),
					('https', 0),
					('use_api', 0),
					('allow_subusers', 0)");

	$query->execute(array(
		':cname' => $params['companyName'],
		':murl' => 'http://' . $params['siteUrl'] . '/',
		':mwebsite' => 'http://' . $params['siteUrl'] . '/',
		':aurl' => '//' . $params['siteUrl'] . '/assets/'
	));

	echo "Settings added\n";

	$uuid = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0x0fff) | 0x4000, mt_rand(0, 0x3fff) | 0x8000, mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff));
	$mysql->prepare("INSERT INTO `users` VALUES(NULL, :uuid, :username, :email, :password, :language, :time, NULL, NULL, 1, 0, 1, 0, NULL)")->execute(array(
		':uuid' => $uuid,
		':username' => $params['adminName'],
		':email' => $params['adminEmail'],
		':password' => password_hash($params['adminPass'], PASSWORD_BCRYPT),
		':language' => 'en',
		':time' => time()
	));

	echo "Admin user added\n";
} catch (\Exception $ex) {

	echo $ex->getMessage() . "\n";
	if (isset($mysql) && $mysql->inTransaction()) {
		$mysql->rollBack();
	}
	exit(1);
}