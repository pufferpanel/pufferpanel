<?php

/*
  PufferPanel - A Minecraft Server Management Panel
  Copyright (c) 2015 PufferPanel

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

use \PDO,
	Components\Config,
	\PufferPanel\Core\Components\Url;

require_once(BASE_DIR . 'vendor/autoload.php');
require_once(SRC_DIR . 'core/autoloader.php');

$klein = new \Klein\Klein();
$twig = new \Twig_Environment(new \Twig_Loader_Filesystem(APP_DIR . 'views/'), array(
	'cache' => false,
	'debug' => true
		));

$klein->respond('/install/[:progress]', function($request, $response) {

	//validate the page the installer should be at is what they are at
	if (isset($_SESSION['installer_progress']) && $_SESSION['installer_progress'] !== $request->param('progress')) {
		$response->redirect('/install/' . $_SESSION['installer_progress'])->send();
	} else if(!isset($_SESSION['installer_progress'])) {
		$response->redirect('/install')->send();
	}

});

$klein->respond('GET', '/install', function($request, $response) use ($twig) {
	include(PANEL_DIR . 'install/install/index.php');
});

$klein->respond('GET', '/install/start', function($request, $response) use ($twig) {
	$response->body($twig->render('install/start.html'))->send();
});

$klein->respond('POST', '/install/start', function($request, $response) {

	try {

		$host = $request->param('sql_h');
		$db = $request->param('sql_db');
		$user = $request->param('sql_u');
		$pass = $request->param('sql_p');

		$database = new PDO('mysql:host=' . $host . ';dbname=' . $db, $user, $pass, array(
			PDO::ATTR_PERSISTENT => true,
			PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
		));

		$database->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

		/* Make File */
		$keyset = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*-=+[]()";
		$hash = "";

		for ($i = 0; $i < 48; $i++) {
			$hash .= substr($keyset, rand(0, strlen($keyset) - 1), 1);
		}

		if (file_exists(BASE_DIR . 'config.json')) {
			throw \Exception('Cannot override existing config');
		}

		$fp = fopen(BASE_DIR . 'config.json', 'w+');
		fwrite($fp, json_encode(array(
			'mysql' => array(
				'host' => $host,
				'database' => $db,
				'username' => $user,
				'password' => $pass
			),
			'hash' => $hash
		)));
		fclose($fp);
		$_SESSION['installer_progress'] = 'tables';
		$response->body(true)->send();

	} catch (\PDOException $e) {
		$response->body($e->getMessage())->send();
	}

});

$klein->respond('GET', '/install/tables', function($request, $response) use ($twig) {
	$response->body($twig->render('install/tables.html'))->send();
});

$klein->respond('POST', '/install/tables', function($request, $response) {

	if (!file_exists(SRC_DIR . 'core/configuration.php')) {

		echo "The configuration file was not found.\n";
		echo "false";
		return;

	}

	try {

		include(SRC_DIR . 'core/configuration.php');
		$mysql = new PDO('mysql:host=' . Config::getGlobal('mysql')->host . ';dbname=' . Config::getGlobal('mysql')->database, Config::getGlobal('mysql')->username, Config::getGlobal('mysql')->password, array(
			PDO::ATTR_PERSISTENT => true,
			PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
		));

		$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
		$mysql->beginTransaction();

		/*
		 * CREATE TABLE `account_change`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `account_change`");
		$mysql->exec("CREATE TABLE `account_change` (
						`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
						`user_id` int(11) DEFAULT NULL,
						`type` varchar(50) NOT NULL DEFAULT '',
						`content` text NOT NULL,
						`key` varchar(50) NOT NULL DEFAULT '',
						`time` int(15) NOT NULL,
						`verified` int(1) NOT NULL DEFAULT '0',
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `account_change` created.\n";

		/*
		 * CREATE TABLE `api`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `api`");
		$mysql->exec("CREATE TABLE `api` (
						`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
						`key` char(36) NOT NULL DEFAULT '',
						`permissions` tinytext NOT NULL,
						`request_ips` tinytext NOT NULL,
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `api` created.\n";

		/*
		 * CREATE TABLE `downloads`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `downloads`");
		$mysql->exec("CREATE TABLE `downloads` (
						`id` int(1) unsigned NOT NULL AUTO_INCREMENT,
						`server` int(1) NOT NULL,
						`token` char(32) NOT NULL DEFAULT '',
						`path` varchar(5000) NOT NULL DEFAULT '',
						PRIMARY KEY (`id`)
					) ENGINE=MEMORY DEFAULT CHARSET=latin1;");
		echo "Table `downloads` created.\n";

		/*
		 * CREATE TABLE `acp_settings`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `acp_settings`");
		$mysql->exec("CREATE TABLE `acp_settings` (
						`id` int(1) unsigned NOT NULL AUTO_INCREMENT,
						`setting_ref` char(25) NOT NULL DEFAULT '',
						`setting_val` tinytext,
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `acp_settings` created.\n";

		/*
		 * CREATE TABLE `actions_log`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `actions_log`");
		$mysql->exec("CREATE TABLE `actions_log` (
						`id` int(1) unsigned NOT NULL AUTO_INCREMENT,
						`priority` int(1) NOT NULL,
						`viewable` int(1) NOT NULL DEFAULT '0',
						`user` int(1) DEFAULT NULL,
						`time` int(1) NOT NULL,
						`ip` char(100) NOT NULL DEFAULT '',
						`url` tinytext NOT NULL,
						`action` char(100) NOT NULL DEFAULT '',
						`desc` text NOT NULL,
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `actions_log` created.\n";

		/*
		 * CREATE TABLE `locations`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `locations`");
		$mysql->exec("CREATE TABLE `locations` (
						`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
						`short` varchar(10) NOT NULL DEFAULT '',
						`long` varchar(500) NOT NULL DEFAULT '',
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `locations` created.\n";

		/*
		 * CREATE TABLE `nodes`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `nodes`");
		$mysql->exec("CREATE TABLE `nodes` (
						`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
						`node` char(15) NOT NULL DEFAULT '',
						`location` varchar(500) NOT NULL,
						`allocate_memory` int(11) NOT NULL,
						`allocate_disk` int(11) DEFAULT NULL,
						`fqdn` tinytext NOT NULL,
						`ip` tinytext NOT NULL,
						`gsd_secret` char(32) DEFAULT NULL,
						`gsd_listen` int(1) DEFAULT '8003',
						`gsd_console` int(1) DEFAULT '8031',
						`gsd_server_dir` tinytext,
						`ips` text NOT NULL,
						`ports` text NOT NULL,
						`public` int(1) NOT NULL DEFAULT '1',
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `nodes` created.\n";

		/*
		 * CREATE TABLE `servers`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `servers`");
		$mysql->exec("CREATE TABLE `servers` (
						`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
						`gsd_id` int(11) DEFAULT NULL,
						`whmcs_id` int(11) DEFAULT '0',
						`hash` char(36) NOT NULL DEFAULT '',
						`gsd_secret` char(36) NOT NULL DEFAULT '',
						`encryption_iv` tinytext NOT NULL,
						`node` int(11) NOT NULL,
						`name` varchar(200) NOT NULL DEFAULT '',
						`modpack` char(16) NOT NULL,
						`server_jar` tinytext,
						`active` int(1) DEFAULT '1',
						`owner_id` int(11) NOT NULL,
						`subusers` tinytext,
						`max_ram` int(11) NOT NULL,
						`disk_space` int(11) NOT NULL,
						`cpu_limit` int(11) DEFAULT NULL,
						`date_added` int(15) NOT NULL,
						`server_ip` varchar(50) NOT NULL DEFAULT '',
						`server_port` int(11) NOT NULL,
						`ftp_user` tinytext NOT NULL,
						`ftp_pass` tinytext NOT NULL,
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `servers` created.\n";

		/*
		 * CREATE TABLE `users`
		 */
		$mysql->exec("DROP TABLE IF EXISTS `users`");
		$mysql->exec("CREATE TABLE `users` (
		  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
						`whmcs_id` int(11) DEFAULT NULL,
						`uuid` varchar(36) NOT NULL,
						`username` varchar(50) NOT NULL DEFAULT '',
						`email` tinytext NOT NULL,
						`password` tinytext NOT NULL,
						`permissions` text,
						`language` char(2) NOT NULL DEFAULT 'en',
						`register_time` int(15) NOT NULL,
						`session_id` varchar(12) DEFAULT '',
						`session_ip` varchar(50) DEFAULT '',
						`root_admin` int(1) NOT NULL DEFAULT '0',
						`notify_login_s` int(1) DEFAULT '1',
						`notify_login_f` int(1) DEFAULT '1',
						`use_totp` int(1) NOT NULL DEFAULT '0',
						`totp_secret` tinytext,
						PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=latin1");
		echo "Table `users` created.\n";

		$mysql->commit();

		$_SESSION['installer_progress'] = 'settings';
		echo "true";

	} catch (\Exception $ex) {

		echo $ex->getMessage() . "\n";
		echo "false";
		if (isset($mysql) && $mysql->inTransaction()) {
			$mysql->rollBack();
		}
		
	}
});

$klein->respond('GET', '/install/settings', function($request, $response, $service) use ($twig) {
	$response->body($twig->render('install/settings.html', array('flash' => $service->flashes())))->send();
});

$klein->respond('POST', '/install/settings', function($request, $response, $service) {

	try {

		$service->validateParam('company_name', 'A company name must be provided')->notNull();
		$service->validateParam('master_url', 'A master URL must be provided')->notNull();
		$service->validateParam('main_website', 'A main website must be provided')->notNull();
		$service->validateParam('assets_url', 'An assets URL must be provided')->notNull();

		$company = $request->param('company_name');
		$masterUrl = Url::addTrailing(Url::stripHttp($request->param('master_url'), true));
		$mainWebsite = Url::addTrailing(Url::stripHttp($request->param('main_website'), true));
		$assetsUrl = Url::addTrailing(Url::stripHttp($request->param('assets_url'), true));

		include(SRC_DIR . 'core/configuration.php');
		$mysql = new PDO('mysql:host=' . Config::getGlobal('mysql')->host . ';dbname=' . Config::getGlobal('mysql')->database, Config::getGlobal('mysql')->username, Config::getGlobal('mysql')->password, array(
			PDO::ATTR_PERSISTENT => true,
			PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
		));

		$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

		$prepare = $mysql->prepare("INSERT INTO `acp_settings` (`setting_ref`, `setting_val`) VALUES
		('company_name', :cname),
		('master_url', :murl),
		('assets_url', :aurl),
		('main_website', :mwebsite),
		('postmark_api_key', NULL),
		('mandrill_api_key', NULL),
		('mailgun_api_key', NULL),
		('sendgrid_api_key', NULL),
		('sendmail_email', :smail),
		('sendmail_method','php'),
		('captcha_pub','6LdSzuYSAAAAAHkmq8LlvmhM-ybTfV8PaTgyBDII'),
		('captcha_priv','6LdSzuYSAAAAAISSAYIJrFGGGJHi5a_V3hGRvIAz'),
		('default_language', 'en'),
		('force_online', 0),
		('https', 0),
		('use_api', 0),
		('allow_subusers', 0)");

		$prepare->execute(array(
			':cname' => $company,
			':murl' => $masterUrl,
			':smail' => $request->param('sendmail_email'),
			':mwebsite' => $mainWebsite,
			':aurl' => $assetsUrl
		));
		$_SESSION['installer_progress'] = 'account';
		$response->body(true)->send();

	} catch (\Exception $ex) {
		$response->body($ex->getMessage())->send();
	}

});

$klein->respond('GET', '/install/account', function($request, $response) use ($twig) {
	$response->body($twig->render('install/account.html'))->send();
});

$klein->respond('POST', '/install/account', function($request, $response) {

	include(SRC_DIR . 'core/configuration.php');

	$mysql = new PDO('mysql:host=' . Config::getGlobal('mysql')->host . ';dbname=' . Config::getGlobal('mysql')->database, Config::getGlobal('mysql')->username, Config::getGlobal('mysql')->password, array(
		PDO::ATTR_PERSISTENT => true,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
	));

	$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

	$prepare = $mysql->prepare("INSERT INTO `users` VALUES(NULL, NULL, :uuid, :username, :email, :password, NULL, :language, :time, NULL, NULL, 1, 0, 1, 0, NULL)");
	$uuid = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0x0fff) | 0x4000, mt_rand(0, 0x3fff) | 0x8000, mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff));

	$prepare->execute(array(
		':uuid' => $uuid,
		':username' => $request->param('username'),
		':email' => $request->param('email'),
		':password' => password_hash($request->param('password'), PASSWORD_BCRYPT),
		':language' => 'en',
		':time' => time()
	));
	rename(SRC_DIR . 'install.lock.dist', SRC_DIR . 'install.lock');
	$response->redirect('/');

});

$klein->onHttpError(function($code, $klein) {

	if ($code == '404') {
		$klein->response()->redirect('/install');
	}
	
});

$klein->dispatch();
