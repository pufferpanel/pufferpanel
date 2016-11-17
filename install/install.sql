-- Remove existing database and create new
CREATE DATABASE IF NOT EXISTS `pufferpanel` DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_unicode_ci;
USE `pufferpanel`;

-- Disable Foreign keys to avoid errors in dropping
SET FOREIGN_KEY_CHECKS = 0;

-- Create the tables needed
DROP TABLE IF EXISTS `acp_settings`;
CREATE TABLE `acp_settings` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `setting_ref` varchar(25) NOT NULL,
  `setting_val` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `setting_ref_unique` (`setting_ref`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `downloads`;
CREATE TABLE `downloads` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `server` char(36) NOT NULL,
  `token` char(32) NOT NULL,
  `path` varchar(5000) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MEMORY;

DROP TABLE IF EXISTS `locations`;
CREATE TABLE `locations` (
  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT,
  `short` varchar(10) NOT NULL,
  `long` varchar(500) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `short_unique` (`short`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL,
  `username` varchar(50),
  `email` varchar(255) NOT NULL,
  `password` text DEFAULT NULL,
  `language` char(2) NOT NULL DEFAULT 'en',
  `register_time` int(15) unsigned NOT NULL,
  `session_id` char(12) DEFAULT '',
  `session_ip` varchar(50) DEFAULT '',
  `root_admin` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `notify_login_s` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `notify_login_f` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `use_totp` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `totp_secret` char(16),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid_unique` (`uuid`),
  UNIQUE KEY `email_unique` (`email`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `account_change`;
CREATE TABLE `account_change` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned,
  `type` varchar(50) NOT NULL DEFAULT '',
  `content` mediumtext NOT NULL,
  `key` mediumtext NOT NULL,
  `time` int(15) unsigned NOT NULL,
  `verified` tinyint(1) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `FK_account_change_users` (`user_id`),
  CONSTRAINT `FK_account_change_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `autodeploy`;
CREATE TABLE `autodeploy` (
  `id` mediumint(10) unsigned NOT NULL AUTO_INCREMENT,
  `node` mediumint(10) unsigned NOT NULL,
  `code` char(36) NOT NULL DEFAULT '',
  `expires` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `actions_log`;
CREATE TABLE `actions_log` (
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

DROP TABLE IF EXISTS `nodes`;
CREATE TABLE `nodes` (
  `id` mediumint(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(15) NOT NULL,
  `location` mediumint(8) unsigned NOT NULL,
  `allocate_memory` mediumint(8) unsigned NOT NULL,
  `allocate_disk` int(10) unsigned NOT NULL,
  `fqdn` varchar(255) NOT NULL,
  `ip` varchar(45) NOT NULL,
  `daemon_secret` char(36) DEFAULT NULL,
  `daemon_listen` smallint(5) unsigned DEFAULT '5656',
  `daemon_sftp` smallint(5) unsigned DEFAULT '22',
  `daemon_base_dir` varchar(200) DEFAULT '/home/',
  `ips` mediumtext NOT NULL,
  `ports` mediumtext NOT NULL,
  `public` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `docker` tinyint(1) unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `FK_nodes_locations` (`location`),
  CONSTRAINT `FK_nodes_locations` FOREIGN KEY (`location`) REFERENCES `locations` (`id`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `servers`;
CREATE TABLE `servers` (
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

DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user` int(10) unsigned NOT NULL,
  `server` int(10) unsigned NOT NULL,
  `permission` varchar(200) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `FK_permissions_users` (`user`),
  KEY `FK_permissions_servers` (`server`),
  CONSTRAINT `FK_permissions_servers` FOREIGN KEY (`server`) REFERENCES `servers` (`id`),
  CONSTRAINT `FK_permissions_users` FOREIGN KEY (`user`) REFERENCES `users` (`id`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `subusers`;
CREATE TABLE `subusers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user` int(10) unsigned NOT NULL,
  `server` int(10) unsigned NOT NULL,
  `daemon_secret` char(36) NOT NULL,
  `daemon_permissions` mediumtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_server_key` (`user`, `server`),
  CONSTRAINT `FK_subusers_user` FOREIGN KEY (`user`) REFERENCES `users` (`id`),
  CONSTRAINT `FK_subusers_server` FOREIGN KEY (`server`) REFERENCES `servers` (`id`)
) ENGINE=InnoDB;

DROP TABLE IF EXISTS `oauth_scopes`;
CREATE TABLE `oauth_scopes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `scope` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT '0',
  `requiresserver` bit(1) NOT NULL DEFAULT b'1',
  `admin` bit(1) NOT NULL DEFAULT b'0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `scope_unique` (`scope`),
  KEY `scope_key` (`scope`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `oauth_clients`;
CREATE TABLE IF NOT EXISTS `oauth_clients` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `client_id` char(16) COLLATE utf8_unicode_ci NOT NULL,
  `client_secret` char(64) COLLATE utf8_unicode_ci NOT NULL,
  `user_id` int(10) unsigned NOT NULL,
  `server_id` int(10) unsigned NOT NULL,
  `name` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
  `description` varchar(1024) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `FK_oauth_clients_users` (`user_id`),
  KEY `FK_oauth_clients_servers` (`server_id`),
  KEY `client_id` (`client_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `oauth_clients_scopes`;
CREATE TABLE IF NOT EXISTS `oauth_clients_scopes` (
  `clientid` int(10) unsigned NOT NULL,
  `scopeid` int(10) unsigned NOT NULL,
  PRIMARY KEY (`clientid`,`scopeid`),
  KEY `clientid` (`clientid`),
  KEY `FK_oauth_scopes` (`scopeid`),
  CONSTRAINT `FK__oauth_clients` FOREIGN KEY (`clientid`) REFERENCES `oauth_clients` (`id`),
  CONSTRAINT `FK_oauth_scopes` FOREIGN KEY (`scopeid`) REFERENCES `oauth_scopes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `oauth_access_tokens`;
CREATE TABLE IF NOT EXISTS `oauth_access_tokens` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `access_token` char(128) COLLATE utf8_unicode_ci NOT NULL,
  `client_id` int(10) unsigned NOT NULL,
  `expiretime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_access_token` (`access_token`),
  KEY `key_access_token` (`access_token`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE IF EXISTS `oauth_access_token_scopes`;
CREATE TABLE IF NOT EXISTS `oauth_access_token_scopes` (
  `id` bigint(20) unsigned NOT NULL,
  `accessTokenId` bigint(20) unsigned NOT NULL,
  `scopeId` int(10) unsigned NOT NULL,
  KEY `FK_oauth_access_token_scopes_oauth_access_tokens` (`accessTokenId`),
  KEY `FK_oauth_access_token_scopes_oauth_scopes` (`scopeId`),
  CONSTRAINT `FK_oauth_access_token_scopes_oauth_access_tokens` FOREIGN KEY (`accessTokenId`) REFERENCES `oauth_access_tokens` (`id`),
  CONSTRAINT `FK_oauth_access_token_scopes_oauth_scopes` FOREIGN KEY (`scopeId`) REFERENCES `oauth_scopes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Enable foreign keys again
SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO `locations` (`id`, `short`, `long`)
VALUES (1, 'Localhost', 'Localhost');

INSERT INTO `oauth_scopes` (`id`, `scope`, `requiresserver`, `admin`) VALUES
	(1, 'server.start', b'1', b'0'),
	(2, 'server.stop', b'1', b'0'),
	(3, 'server.install', b'1', b'0'),
	(4, 'server.delete', b'1', b'1'),
	(5, 'server.edit', b'1', b'1'),
	(6, 'server.file.put', b'1', b'0'),
	(8, 'server.file.get', b'1', b'0'),
	(9, 'server.console', b'1', b'0'),
	(10, 'server.console.send', b'1', b'0'),
	(11, 'server.stats', b'1', b'0'),
	(12, 'server.reload', b'1', b'0'),
	(13, 'server.network', b'0', b'1');


DROP PROCEDURE IF EXISTS `oauthCreateClient`;
DELIMITER //
CREATE PROCEDURE `oauthCreateClient`(
	IN `clientId` VARCHAR(128),
	IN `userId` BIGINT,
	IN `serverId` BIGINT,
	IN `scopes` VARCHAR(2000),
	IN `displayname` VARCHAR(128),
	IN `description` VARCHAR(1024)


)
Main:BEGIN
	DECLARE clientSecret varchar(128);
	
	IF clientId IS NULL OR userId IS NULL
	THEN
		LEAVE Main;
	END IF;
	
	IF displayName IS NULL
	THEN
		SET displayName = clientId;
	END IF;
		
	IF description IS NULL
	THEN
		SET description = 'No description';
	END IF;
	
	INSERT INTO oauth_clients VALUES (NULL, clientId, clientSecret, userId, serverId, displayName, description);
	
	SELECT clientSecret AS secret;			
END//
DELIMITER ;

DROP FUNCTION IF EXISTS `oauthGenerateToken`;
DELIMITER //
CREATE FUNCTION `oauthGenerateToken`() RETURNS varchar(32) CHARSET utf8 COLLATE utf8_unicode_ci
BEGIN
	RETURN REPLACE(UUID(), '-', '');
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `oauthGetOrGenAccessToken`;
DELIMITER //
CREATE PROCEDURE `oauthGetOrGenAccessToken`(
	IN `userId` BIGINT,
	IN `serverId` BIGINT

)
BEGIN
	DECLARE accessToken varchar(128);
	DECLARE clientId varchar(128);
	DECLARE clientSecret varchar(128);
	DECLARE success BIT;

	SELECT access_token INTO accessToken FROM oauth_access_tokens AS oat
		INNER JOIN oauth_clients AS oc ON oc.id = oat.client_id
		WHERE user_id = userId AND server_id = serverId AND expiretime > NOW();
		
	IF accessToken IS NULL
	THEN
		SELECT NULL INTO accessToken;
		
		SELECT client_id, client_secret INTO clientId, clientSecret FROM oauth_clients WHERE user_id = userId AND server_id = serverId;
		CALL oauthHandleTokenCredentials(clientId, clientSecret, @success);
		IF success = 1
		THEN 
			SELECT access_token INTO accessToken FROM oauth_access_tokens AS oat
				INNER JOIN oauth_clients AS oc ON oc.id = oat.client_id
				WHERE user_id = userId AND server_id = serverId AND expiretime > NOW();
		END IF;
	END IF;
	
	SELECT accessToken AS accessToken;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `oauthHandleInfoRequest`;
DELIMITER //
CREATE PROCEDURE `oauthHandleInfoRequest`(IN `token` VARCHAR(128))
BEGIN
	DECLARE id bigint;
	DECLARE server varchar(128);
	DECLARE expiretime timestamp;
	DECLARE clientId varchar(128);

	SELECT user_id, IFNULL(hash, '*') AS server_id, expiretime, oc.client_id
		INTO id, server, expiretime, clientid
		FROM oauth_clients AS oc
		INNER JOIN oauth_access_tokens AS oat ON oat.client_id = oc.id
		LEFT JOIN servers AS s ON s.id = oc.server_id
		WHERE access_token = token AND expiretime > NOW();
		
	SELECT id AS userId, server AS serverId, expireTime AS expireTime, clientId AS clientId;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `oauthHandleResourceOwner`;
DELIMITER //
CREATE PROCEDURE `oauthHandleResourceOwner`(IN `email` VARCHAR(128), IN `password` TEXT, IN `server` VARCHAR(128), OUT `success` BIT)
Main:BEGIN
	DECLARE resultUserId bigint;
	DECLARE resultServerId bigint;
	DECLARE resultOauthId bigint;
	DECLARE token varchar(128);
	DECLARE expireTime timestamp;
	DECLARE tokenId bigint;
	
	SELECT 0 INTO success;
	
	SELECT id INTO resultUserId FROM users WHERE users.email = email AND users.password = password;
	
	IF resultUserId IS NULL
	THEN
		LEAVE Main;
	END IF;
	
	SELECT s.id INTO resultServerId FROM servers AS s LEFT JOIN subusers AS su ON su.server = s.id WHERE s.name = server AND (s.owner_id = resultUserId OR su.user = resultUserId);
	
	IF resultServerId IS NULL
	THEN
		LEAVE Main;
	END IF;
	
	SELECT id INTO resultOauthId FROM oauth_clients WHERE user_id = resultUserId AND server_id = resultServerId;
	
	
	IF resultOauthId IS NOT NULL
	THEN
		SELECT oauthGenerateToken() INTO token;
		SELECT ADDTIME(NOW(), '0 1:0:0.0') INTO expireTime;
		
		INSERT INTO oauth_access_tokens VALUES (NULL, token, resultId, expireTime);
		
		SELECT LAST_INSERT_ID() INTO tokenId;
		
		SELECT token AS accessToken, expireTime AS expireTime;
		SELECT 1 INTO success;
	END IF;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `oauthHandleTokenCredentials`;
DELIMITER //
CREATE PROCEDURE `oauthHandleTokenCredentials`(IN `clientid` CHAR(16), IN `clientsecret` CHAR(64), OUT `success` BIT)
BEGIN
	DECLARE resultId bigint;
	DECLARE token varchar(128);
	DECLARE expireTime timestamp;
	DECLARE tokenId bigint;
	
	SELECT id INTO resultId FROM oauth_clients WHERE client_id = clientid AND client_secret = clientsecret;
	
	IF resultId IS NOT NULL
	THEN
		SELECT oauthGenerateToken() INTO token;
		SELECT ADDTIME(NOW(), '0 1:0:0.0') INTO expireTime;

		INSERT INTO oauth_access_tokens VALUES (NULL, token, resultId, expireTime);
		
		SELECT token AS accessToken, expireTime AS expireTime;		
		SELECT LAST_INSERT_ID() INTO tokenId;
		
		SELECT 1 INTO success;
	ELSE
		SELECT 0 INTO success;
	END IF;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `oauthRemoveAccessToken`;
DELIMITER //
CREATE PROCEDURE `oauthRemoveAccessToken`(IN `accessToken` VARCHAR(128), IN `accessTokenId` BIGINT)
Main:BEGIN
	IF accessToken IS NULL AND accessTokenId IS NULL
	THEN
		LEAVE Main;
	END IF;	
	
	IF accessToken IS NOT NULL
	THEN 
		SELECT oat.id INTO accessTokenId FROM oauth_access_tokens AS oat WHERE oat.accessToken = accessToken;
		
		IF accessTokenId IS NULL
		THEN
			LEAVE Main;
		END IF;	
	END IF;
	
	DELETE FROM oauth_access_token_scopes WHERE oauth_access_token_scopes.accessTokenId = accessTokenId;
	DELETE FROM oauth_access_tokens WHERE id = accessTokenId;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `oauthRemoveClientByClientId`;
DELIMITER //
CREATE PROCEDURE `oauthRemoveClientByClientId`(IN `clientId` BIGINT)
Main:BEGIN	
	DECLARE done INT DEFAULT FALSE;
	DECLARE tokenId BIGINT;
	DECLARE curs CURSOR FOR SELECT id FROM oauth_access_tokens AS oat WHERE oat.client_id = clientId;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	IF clientId IS NULL
	THEN 
		LEAVE Main;
	END IF;
	
	OPEN curs;
	
	curs_loop: LOOP
		FETCH curs INTO tokenId;
		IF done THEN LEAVE curs_loop; END IF;
		CALL oauthRemoveAccessToken(null, tokenId);
	END LOOP;
	
	DELETE FROM oauth_clients_scopes WHERE oauth_clients_scopes.clientid = clientId;
	DELETE FROM oauth_clients WHERE oauth_clients.client_id = clientId;
END//
DELIMITER ;

DROP PROCEDURE IF EXISTS `oauthRemoveClientById`;
DELIMITER //
CREATE PROCEDURE `oauthRemoveClientById`(IN `id` BIGINT)
Main:BEGIN	
	DECLARE done INT DEFAULT FALSE;
	DECLARE tokenId BIGINT;
	DECLARE curs CURSOR FOR 
		SELECT oat.id FROM oauth_access_tokens AS oat 
		INNER JOIN oauth_clients AS oc
			ON oc.client_id = oat.client_id
		WHERE oc.id = id;
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

	IF id IS NULL
	THEN 
		LEAVE Main;
	END IF;
	
	OPEN curs;
	
	curs_loop: LOOP
		FETCH curs INTO tokenId;
		IF done THEN LEAVE curs_loop; END IF;
		CALL oauthRemoveAccessToken(null, tokenId);
	END LOOP;
	
	DELETE FROM oauth_clients WHERE oauth_clients.id = clientId;
END//
DELIMITER ;
