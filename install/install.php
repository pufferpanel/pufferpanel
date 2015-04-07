<?php
/**
 *
 *                 __.............__
 *          .--""```                 ```""--.
 *           ':--..___             ___..--:'
 *             \      ```"""""""```      /
 *           .-`  ___.....-----.....___  '-.
 *         .:-""``     ~          ~    ``""-:.
 *        /`-..___ ~        ~         ~___..-'\
 *       /  ~    '`""---.........---""`        \
 *      ;                                       ;
 *     ; '::.   '          _,           _,       ;
 *     |   ':::    '     .' (    ~   .-'./    ~  |
 *     |~  .:'   .     _/..._'.    .'.-'/        |
 *     | .:'       .-'`      ` '-./.'_.'         |
 *     |  ':.     ( o)   ))      ;= <_           |
 *     ; '::.      '-.,\\__ __.-;`\'. '.  .      ;
 *      ;    ':         \) |`\ \)  '.'-.\       ;
 *       \.:'.:':.         \_/       '-._\     /
 *        \ ':.     ~                    `    /
 *         '. '::..  _ . - - -- .~ _      ~ .'
 *           '-._':'                 `'-_.-'
 *              (``''--..._____...--''``)
 *               `"--...__     __...--"`
 *                        `````
 *
 *
 * I did not waste hours and hours writing beautiful documentation
 * just so you could ignore it and do some voodo-magic. Please
 * do not try to edit this file, and if you even have this open
 * I am surprised. Shoo, this is a fish only domain. Use the installer
 * and read the documentation. pls.
 *
 * http://pufferpanel.com/
 *
 **/
use \PDO;

define("BASE_DIR", dirname(dirname(__FILE__)).'/');

$params = array();
parse_str(implode('&', array_splice($argv, 1)), $params);

if(empty($params)) {
	echo "You failed to read the docs. Go read them again";
	return;
}

$keyset = "abcdefghijklmABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*-=+[]()";
$hash = "";

for($i = 0; $i < 48; $i++) {
	$hash .= substr($keyset, rand(0, strlen($keyset) - 1), 1);
}

$pass = "";
for($i = 0; $i < 48; $i++) {
	$pass .= substr($keyset, rand(0, strlen($keyset) - 1), 1);
}

try {

	$fp = fopen(BASE_DIR.'config.json', 'w');
	if($fp === false) {
		throw new \Exception('Could not open config.json');
	}

	fwrite($fp, json_encode(array(
		'mysql' => array(
			'host' => $params['mysqlHost'],
			'database' => 'pufferpanel',
			'username' => 'pufferpanel',
			'password' => $pass,
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

	if(!file_exists(BASE_DIR.'config.json')) {
		throw new \Exception("Could not create config.json");
	}

	$mysql = new PDO('mysql:host='.$params['mysqlHost'], $params['mysqlUser'], $params['mysqlPass'], array(
		PDO::ATTR_PERSISTENT => true,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
	));

	$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
	$mysql->beginTransaction();

	$mysql->exec("DROP DATABASE IF EXISTS `pufferpanel`");
	$mysql->exec("CREATE DATABASE `pufferpanel`");

	$mysql->exec("USE `pufferpanel`");

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
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
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
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
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
					) ENGINE=MEMORY DEFAULT CHARSET=utf8;");
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
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
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
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
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
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
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
					`allocate_disk` int(11) NOT NULL,
					`fqdn` tinytext NOT NULL,
					`ip` tinytext NOT NULL,
					`daemon_secret` char(36) DEFAULT NULL,
					`daemon_listen` int(1) DEFAULT '5656',
					`daemon_console` int(1) DEFAULT '5657',
					`daemon_upload` int(1) DEFAULT '5658',
					`daemon_sftp` int(1) DEFAULT '22',
					`daemon_base_dir` varchar(200) DEFAULT '/home/',
					`ips` text NOT NULL,
					`ports` text NOT NULL,
					`public` int(1) NOT NULL DEFAULT '1',
					PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
	echo "Table `nodes` created.\n";

	/*
	 * CREATE TABLE `servers`
	 */
	$mysql->exec("DROP TABLE IF EXISTS `servers`");
	$mysql->exec("CREATE TABLE `servers` (
					`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
					`hash` char(36) NOT NULL DEFAULT '',
					`daemon_secret` char(36) NOT NULL DEFAULT '',
					`node` int(11) NOT NULL,
					`name` varchar(200) NOT NULL DEFAULT '',
					`plugin` char(100) NOT NULL DEFAULT '',
					`server_jar` tinytext,
					`active` int(1) DEFAULT '1',
					`owner_id` int(11) NOT NULL,
					`max_ram` int(11) NOT NULL,
					`disk_space` int(11) NOT NULL,
					`cpu_limit` int(11) DEFAULT NULL,
					`date_added` int(15) NOT NULL,
					`server_ip` varchar(50) NOT NULL DEFAULT '',
					`server_port` int(11) NOT NULL,
					`sftp_user` tinytext NOT NULL
					PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
	echo "Table `servers` created.\n";

	/*
	 * CREATE TABLE `users`
	 */
	$mysql->exec("DROP TABLE IF EXISTS `users`");
	$mysql->exec("CREATE TABLE `users` (
					`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
					`whmcs_id` int(11) DEFAULT NULL,
					`uuid` char(36) NOT NULL,
					`username` varchar(50) NOT NULL DEFAULT '',
					`email` tinytext NOT NULL,
					`password` tinytext NOT NULL,
					`language` char(2) NOT NULL DEFAULT 'en',
					`register_time` int(15) NOT NULL,
					`session_id` varchar(12) DEFAULT '',
					`session_ip` varchar(50) DEFAULT '',
					`root_admin` int(1) NOT NULL DEFAULT '0',
					`notify_login_s` int(1) NOT NULL DEFAULT '1',
					`notify_login_f` int(1) NOT NULL DEFAULT '1',
					`use_totp` int(1) NOT NULL DEFAULT '0',
					`totp_secret` tinytext,
					PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8");
	echo "Table `users` created.\n";

	$mysql->exec("DROP TABLE IF EXISTS `permissions`");
	$mysql->exec("CREATE TABLE `permissions` (
					`id` int(1) unsigned NOT NULL AUTO_INCREMENT,
					`user` int(1) NOT NULL,
					`server` int(1) NOT NULL,
					`permission` varchar(200) NOT NULL DEFAULT '',
					PRIMARY KEY (`id`)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8");
	echo "Table `permissions` created.\n";

	$mysql->exec("DROP TABLE IF EXISTS `plugins`");
	$mysql->exec("CREATE TABLE `plugins` (
					`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
					`hash` char(36) NOT NULL DEFAULT '',
					`name` varchar(100) NOT NULL DEFAULT '',
					`description` text NOT NULL,
					`slug` varchar(100) NOT NULL DEFAULT '',
					PRIMARY KEY (`id`)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8;");
	echo "Table `plugins` created.\n";

	$mysql->exec("DROP TABLE IF EXISTS `subusers`");

	$mysql->exec("CREATE TABLE `subusers` (
					`id` int(1) unsigned NOT NULL AUTO_INCREMENT,
					`uuid` char(36) NOT NULL DEFAULT '',
					`user` int(1) NOT NULL,
					`server` int(1) NOT NULL,
					`daemon_secret` char(36) NOT NULL DEFAULT '',
					`daemon_permissions` text,
					`permissions` text,
					`pending` int(1) NOT NULL,
					`pending_email` varchar(200) DEFAULT NULL,
					PRIMARY KEY (`id`)
				) ENGINE=InnoDB DEFAULT CHARSET=utf8");
	echo "Table `subusers` created.\n";

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
		':murl' => $params['siteUrl'].'/',
		':mwebsite' => $params['siteUrl'].'/',
		':aurl' => '//'.$params['siteUrl'].'/assets/'
	));

	echo "Settings added\n";

	$uuid = sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0x0fff) | 0x4000, mt_rand(0, 0x3fff) | 0x8000, mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff));
	$mysql->prepare("INSERT INTO `users` VALUES(NULL, NULL, :uuid, :username, :email, :password, NULL, :language, :time, NULL, NULL, 1, 0, 1, 0, NULL)")->execute(array(
		':uuid' => $uuid,
		':username' => $params['adminName'],
		':email' => $params['adminEmail'],
		':password' => password_hash($params['adminPass'], PASSWORD_BCRYPT),
		':language' => 'en',
		':time' => time()
	));

	echo "Admin user added\n";

	try {
		$mysql->prepare("DROP USER 'pufferpanel'@'localhost'")->execute();
	} catch(\Exception $ex) {
		//ignoring because no user actually existed
	}
	$mysql->prepare("GRANT SELECT, UPDATE, DELETE, ALTER, INSERT ON pufferpanel.* TO 'pufferpanel'@'localhost' IDENTIFIED BY :pass")->execute(array(
		'pass' => $pass
	));
	echo "PufferPanel SQL user added\n";

	$mysql->commit();

	exit(0);
} catch(\Exception $ex) {

	echo $ex->getMessage()."\n";
	if(isset($mysql) && $mysql->inTransaction()) {
		$mysql->rollBack();
	}
	exit(1);
}
