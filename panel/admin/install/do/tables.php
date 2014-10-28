<?php
/*
    PufferPanel - A Minecraft Server Management Panel
    Copyright (c) 2013 Dane Everitt

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

if(file_exists('../install.lock'))
	exit('Installer is Locked.');
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<link rel="stylesheet" href="../../../assets/css/bootstrap.css">
	<title>PufferPanel Installer</title>
</head>
<body>
	<div class="container">
		<div class="alert alert-danger">
			<strong>WARNING:</strong> Do not run this version on a live environment! There are known security holes that we are working on getting patched. This is extremely beta software and this version is to get the features in place while we work on security enhancements.
		</div>
		<div class="navbar navbar-default">
			<div class="navbar-header">
				<a class="navbar-brand" href="#">Install PufferPanel - Database Content</a>
			</div>
		</div>
		<div class="col-12">
			<div class="row">
				<div class="col-2"></div>
				<div class="col-8">
					<pre>
                    <?php

                    	if(!file_exists('../../../../src/framework/configuration.php'))
                    		echo '<div class="alert alert-danger">The configuration file was not found.</div>';
                    	else {

	                        include('../../../../src/core/database.php');
	                        $mysql = Components\Database::connect();

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
	                        echo "\nTable `account_change` created.\n";

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
	                         * CREATE TABLE `acp_settings`
	                         */
	                        $mysql->exec("DROP TABLE IF EXISTS `acp_settings`");
	                        $mysql->exec("CREATE TABLE `acp_settings` (
	                          `setting_ref` char(25) NOT NULL DEFAULT '',
	                          `setting_val` tinytext
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
	                         * CREATE TABLE `nodes`
	                         */
	                        $mysql->exec("DROP TABLE IF EXISTS `nodes`");
	                        $mysql->exec("CREATE TABLE `nodes` (
	                          `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	                          `node` char(15) NOT NULL DEFAULT '',
	                          `fqdn` tinytext NOT NULL,
	                          `ip` tinytext NOT NULL,
	                          `gsd_secret` char(32) DEFAULT NULL,
	                          `ips` text NOT NULL,
	                          `ports` text NOT NULL,
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

						}

                    ?>
					</pre>
					<form action="settings.php">
					    <input type="submit" class="btn btn-primary" value="Continue &rarr;">
					</form>
				</div>
				<div class="col-2"></div>
			</div>
		</div>
		<div class="footer">
            <div class="col-8 nopad"><p>PufferPanel is licensed under a <a href="https://github.com/DaneEveritt/PufferPanel/blob/master/LICENSE">GPL-v3 License</a>.<br />Running <?php echo trim(file_get_contents('../../../../src/versions/current')).' ('.substr(trim(file_get_contents('../../../../.git/HEAD')), 0, 8).')'; ?> distributed by <a href="http://pufferpanel.com">PufferPanel Development</a>.</p></div>
		</div>
	</div>
</body>
</html>
