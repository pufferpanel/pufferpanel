-- Remove existing database and create new
CREATE DATABASE IF NOT EXISTS `pufferpanel` DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_unicode_ci;
USE `pufferpanel`;

-- Disable Foreign keys to avoid errors in dropping
SET FOREIGN_KEY_CHECKS = 0;

-- Create the tables needed
CREATE TABLE IF NOT EXISTS `acp_settings` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `setting_ref` varchar(25) NOT NULL,
  `setting_val` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `setting_ref_unique` (`setting_ref`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `downloads` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `server` char(36) NOT NULL,
  `token` char(32) NOT NULL,
  `path` varchar(5000) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MEMORY;

CREATE TABLE IF NOT EXISTS `locations` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `short` varchar(10) NOT NULL,
  `long` varchar(500) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `short_unique` (`short`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL,
  `username` varchar(50) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` text DEFAULT NULL,
  `language` char(2) NOT NULL DEFAULT 'en',
  `register_time` int(15) unsigned NOT NULL,
  `session_id` char(12) DEFAULT NULL,
  `session_ip` varchar(50) DEFAULT NULL,
  `root_admin` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `notify_login_s` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `notify_login_f` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `use_totp` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `totp_secret` char(16),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid_unique` (`uuid`),
  UNIQUE KEY `email_unique` (`email`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `account_change` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `type` varchar(50) NOT NULL DEFAULT '',
  `content` mediumtext NOT NULL,
  `key` mediumtext NOT NULL,
  `time` int(15) unsigned NOT NULL,
  `verified` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `FK_account_change_users` (`user_id`),
  CONSTRAINT `FK_account_change_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `autodeploy` (
  `id` mediumint(10) unsigned NOT NULL AUTO_INCREMENT,
  `node` mediumint(10) unsigned NOT NULL,
  `code` char(36) NOT NULL DEFAULT '',
  `expires` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `actions_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `priority` tinyint(1) NOT NULL,
  `viewable` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `user` int(10) unsigned DEFAULT NULL,
  `time` int(15) unsigned NOT NULL,
  `ip` varchar(100) NOT NULL,
  `url` text NOT NULL,
  `action` varchar(100) NOT NULL,
  `desc` mediumtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_actions_log_users` (`user`),
  CONSTRAINT `FK_actions_log_users` FOREIGN KEY (`user`) REFERENCES `users` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `nodes` (
  `id` mediumint(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(15) NOT NULL,
  `location` mediumint(8) unsigned NOT NULL,
  `allocate_memory` mediumint(8) unsigned NOT NULL,
  `allocate_disk` int(10) unsigned NOT NULL,
  `fqdn` varchar(255) NOT NULL,
  `ip` varchar(45) NOT NULL,
  `daemon_secret` char(36) NOT NULL,
  `daemon_listen` smallint(5) unsigned DEFAULT '5656',
  `daemon_sftp` smallint(5) unsigned DEFAULT '5657',
  `ips` mediumtext NOT NULL,
  `ports` mediumtext NOT NULL,
  `public` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `docker` tinyint(1) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `FK_nodes_locations` (`location`),
  CONSTRAINT `FK_nodes_locations` FOREIGN KEY (`location`) REFERENCES `locations` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `servers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `hash` char(36) NOT NULL,
  `daemon_secret` char(36) NOT NULL,
  `node` mediumint(8) unsigned NOT NULL,
  `name` varchar(200) NOT NULL,
  `active` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `owner_id` int(10) unsigned NOT NULL,
  `date_added` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_servers_users` (`owner_id`),
  KEY `FK_servers_nodes` (`node`),
  CONSTRAINT `FK_servers_nodes` FOREIGN KEY (`node`) REFERENCES `nodes` (`id`),
  CONSTRAINT `FK_servers_users` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`)
 ) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `permissions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user` int(10) unsigned NOT NULL,
  `server` int(10) unsigned NOT NULL,
  `permission` varchar(200) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `FK_permissions_users` (`user`),
  KEY `FK_permissions_servers` (`server`),
  CONSTRAINT `FK_permissions_servers` FOREIGN KEY (`server`) REFERENCES `servers` (`id`),
  CONSTRAINT `FK_permissions_users` FOREIGN KEY (`user`) REFERENCES `users` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `subusers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user` int(10) unsigned NOT NULL,
  `server` int(10) unsigned NOT NULL,
  `daemon_permissions` mediumtext NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_server_key` (`user`, `server`),
  CONSTRAINT `FK_subusers_user` FOREIGN KEY (`user`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_subusers_server` FOREIGN KEY (`server`) REFERENCES `servers` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `oauth_clients` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`client_id` CHAR(16) NOT NULL,
	`client_secret` CHAR(64) NOT NULL,
	`user_id` INT(10) UNSIGNED NOT NULL,
	`server_id` INT(10) UNSIGNED NOT NULL,
	`scopes` VARCHAR(1000) NOT NULL DEFAULT '',
	`name` VARCHAR(128) NOT NULL,
	`description` VARCHAR(1024) NOT NULL DEFAULT 'No description',
	PRIMARY KEY (`id`),
	INDEX `FK_oauth_clients_users` (`user_id`),
	INDEX `FK_oauth_clients_servers` (`server_id`),
	INDEX `client_id` (`client_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `oauth_access_tokens` (
  `access_token` CHAR(128) NOT NULL,
  `oauthClientId` INT(10) UNSIGNED NOT NULL,
  `expiretime` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `scopes` VARCHAR(1000) NOT NULL DEFAULT '',
  PRIMARY KEY (`access_token`),
  UNIQUE INDEX `access_token` (`access_token`),
  INDEX `FK_oauth_access_tokens_oauth_clients` (`oauthClientId`),
  CONSTRAINT `FK_oauth_access_tokens_oauth_clients` FOREIGN KEY (`oauthClientId`) REFERENCES `oauth_clients` (`id`)
) ENGINE=InnoDB;



-- Enable foreign keys again
SET FOREIGN_KEY_CHECKS = 1;

INSERT IGNORE INTO `locations` (`id`, `short`, `long`)
VALUES (1, 'Localhost', 'Localhost');

DROP EVENT IF EXISTS `oauthTokenCleaner`;
CREATE EVENT `oauthTokenCleaner`
  ON SCHEDULE
    EVERY 12 HOUR
  ON COMPLETION NOT PRESERVE
ENABLE
  COMMENT 'Cleans up the oauth access tokens'
DO DELETE FROM oauth_access_tokens WHERE expireTime < NOW()