DROP DATABASE IF EXISTS `pufferpanel`;
CREATE DATABASE IF NOT EXISTS `pufferpanel`;
USE `pufferpanel`;

CREATE TABLE IF NOT EXISTS `locations` (
  `id` binary(16) NOT NULL,
  `code` varchar(16) NOT NULL,
  `name` varchar(100) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `nodes` (
	`id` BINARY(36) NOT NULL,
	`location_id` BINARY(36) NOT NULL,
	`code` VARCHAR(20) NOT NULL,
	`name` VARCHAR(100) NOT NULL,
	`external_ip` VARCHAR(45) NOT NULL,
	`internal_ip` VARCHAR(45) NOT NULL,
	`port` SMALLINT(6) NOT NULL,
	`sftp_port` SMALLINT(6) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT '',
	`updated_at` TIMESTAMP NULL DEFAULT NULL,
	PRIMARY KEY (`id`),
	INDEX `FK__locations` (`location_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `servers` (
  `id` binary(16) NOT NULL,
  `node_id` binary(16) NOT NULL,
  `name` varchar(100) NOT NULL,
  `sftp_id` varchar(10) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `FK__nodes` (`node_id`),
  CONSTRAINT `FK__nodes` FOREIGN KEY (`node_id`) REFERENCES `nodes` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `users` (
  `id` binary(12) NOT NULL,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` binary(60) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `clients` (
  `id` BINARY(36) NOT NULL,
  `secret` BINARY(60) NOT NULL,
  `user_id` BINARY(36) NOT NULL,
  `internal` BIT(1) NOT NULL DEFAULT b'0',
  `description` VARCHAR(200) NULL DEFAULT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `FK_clients_users` (`user_id`),
  CONSTRAINT `FK_clients_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `scopes` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `client_id` BINARY(36) NOT NULL,
  `server_id` BINARY(36) NULL DEFAULT NULL,
  `value` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT current_timestamp(),
  `updated_at` TIMESTAMP NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `FK_scopes_clients` (`client_id`),
  INDEX `FK_scopes_servers` (`server_id`),
  CONSTRAINT `FK_scopes_clients` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`) ON DELETE CASCADE,
  CONSTRAINT `FK_scopes_servers` FOREIGN KEY (`server_id`) REFERENCES `servers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;