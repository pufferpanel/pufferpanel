# Dump of table account_change
# ------------------------------------------------------------

DROP TABLE IF EXISTS `account_change`;

CREATE TABLE `account_change` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `type` varchar(50) NOT NULL DEFAULT '',
  `content` mediumtext NOT NULL,
  `key` mediumtext NOT NULL,
  `time` int(15) NOT NULL,
  `verified` int(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table acp_settings
# ------------------------------------------------------------

DROP TABLE IF EXISTS `acp_settings`;

CREATE TABLE `acp_settings` (
  `id` int(1) unsigned NOT NULL AUTO_INCREMENT,
  `setting_ref` char(25) NOT NULL DEFAULT '',
  `setting_val` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table actions_log
# ------------------------------------------------------------

DROP TABLE IF EXISTS `actions_log`;

CREATE TABLE `actions_log` (
  `id` int(1) unsigned NOT NULL AUTO_INCREMENT,
  `priority` int(1) NOT NULL,
  `viewable` int(1) NOT NULL DEFAULT '0',
  `user` int(1) DEFAULT NULL,
  `time` int(1) NOT NULL,
  `ip` char(100) NOT NULL DEFAULT '',
  `url` text NOT NULL,
  `action` char(100) NOT NULL DEFAULT '',
  `desc` mediumtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table downloads
# ------------------------------------------------------------

DROP TABLE IF EXISTS `downloads`;

CREATE TABLE `downloads` (
  `id` int(1) unsigned NOT NULL AUTO_INCREMENT,
  `server` char(36) NOT NULL DEFAULT '',
  `token` char(32) NOT NULL DEFAULT '',
  `path` varchar(5000) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=MEMORY DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table locations
# ------------------------------------------------------------

DROP TABLE IF EXISTS `locations`;

CREATE TABLE `locations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `short` varchar(10) NOT NULL DEFAULT '',
  `long` varchar(500) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table nodes
# ------------------------------------------------------------

DROP TABLE IF EXISTS `nodes`;

CREATE TABLE `nodes` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `node` char(15) NOT NULL DEFAULT '',
  `location` varchar(500) NOT NULL,
  `allocate_memory` int(11) NOT NULL,
  `allocate_disk` int(11) NOT NULL,
  `fqdn` text NOT NULL,
  `ip` text NOT NULL,
  `daemon_secret` char(36) DEFAULT NULL,
  `daemon_listen` int(1) DEFAULT '5656',
  `daemon_console` int(1) DEFAULT '5657',
  `daemon_upload` int(1) DEFAULT '5658',
  `daemon_sftp` int(1) DEFAULT '22',
  `daemon_base_dir` varchar(200) DEFAULT '/home/',
  `ips` mediumtext NOT NULL,
  `ports` mediumtext NOT NULL,
  `public` int(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table permissions
# ------------------------------------------------------------

DROP TABLE IF EXISTS `permissions`;

CREATE TABLE `permissions` (
  `id` int(1) unsigned NOT NULL AUTO_INCREMENT,
  `user` int(1) NOT NULL,
  `server` int(1) NOT NULL,
  `permission` varchar(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table plugins
# ------------------------------------------------------------

DROP TABLE IF EXISTS `plugins`;

CREATE TABLE `plugins` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `hash` char(36) NOT NULL DEFAULT '',
  `name` varchar(100) NOT NULL DEFAULT '',
  `description` text NOT NULL,
  `slug` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

LOCK TABLES `plugins` WRITE;
/*!40000 ALTER TABLE `plugins` DISABLE KEYS */;

INSERT INTO `plugins` (`id`, `hash`, `name`, `description`, `slug`)
VALUES
	(2,'37d8949d-5da2-4390-a28f-27ac1babc4da','Minecraft','Minecraft is a game about breaking and placing blocks. At first, people built structures to protect against nocturnal monsters, but as the game grew players worked together to create wonderful, imaginative things. This version of the plugin is ment for versions of the game <strong>greater than 1.7.0</strong>.','minecraft'),
	(3,'b4b90feb-6adb-499c-a9f8-09b6e80c9d16','Minecraft (pre 1.7)','Minecraft is a game about breaking and placing blocks. At first, people built structures to protect against nocturnal monsters, but as the game grew players worked together to create wonderful, imaginative things. This version of the plugin is ment for versions of the game <strong>less than 1.7.0</strong>.','minecraft-pre');

/*!40000 ALTER TABLE `plugins` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table servers
# ------------------------------------------------------------

DROP TABLE IF EXISTS `servers`;

CREATE TABLE `servers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `hash` char(36) NOT NULL DEFAULT '',
  `daemon_secret` char(36) NOT NULL DEFAULT '',
  `node` int(11) NOT NULL,
  `name` varchar(200) NOT NULL DEFAULT '',
  `plugin` char(100) NOT NULL DEFAULT '',
  `server_jar` text,
  `active` int(1) DEFAULT '1',
  `owner_id` int(11) NOT NULL,
  `max_ram` int(11) NOT NULL,
  `disk_space` int(11) NOT NULL,
  `cpu_limit` int(11) DEFAULT NULL,
  `date_added` int(15) NOT NULL,
  `server_ip` varchar(50) NOT NULL DEFAULT '',
  `server_port` int(11) NOT NULL,
  `sftp_user` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table subusers
# ------------------------------------------------------------

DROP TABLE IF EXISTS `subusers`;

CREATE TABLE `subusers` (
  `id` int(1) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL DEFAULT '',
  `user` int(1) NOT NULL,
  `server` int(1) NOT NULL,
  `daemon_secret` char(36) NOT NULL DEFAULT '',
  `daemon_permissions` mediumtext,
  `permissions` mediumtext,
  `pending` int(1) NOT NULL,
  `pending_email` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `whmcs_id` int(11) DEFAULT NULL,
  `uuid` char(36) NOT NULL,
  `username` varchar(50) NOT NULL DEFAULT '',
  `email` text NOT NULL,
  `password` text NOT NULL,
  `language` char(2) NOT NULL DEFAULT 'en',
  `register_time` int(15) NOT NULL,
  `session_id` varchar(12) DEFAULT '',
  `session_ip` varchar(50) DEFAULT '',
  `root_admin` int(1) NOT NULL DEFAULT '0',
  `notify_login_s` int(1) NOT NULL DEFAULT '1',
  `notify_login_f` int(1) NOT NULL DEFAULT '1',
  `use_totp` int(1) NOT NULL DEFAULT '0',
  `totp_secret` text,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;