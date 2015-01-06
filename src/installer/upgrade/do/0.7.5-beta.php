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

// upgrader for version 0.7.5 Beta to 0.7.6 Beta
require '../../../../src/core/configuration.php';

function uuid() {
	return sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 0xffff), mt_rand(0, 0xffff),
	mt_rand(0, 0xffff),
	mt_rand(0, 0x0fff) | 0x4000,
	mt_rand(0, 0x3fff) | 0x8000,
	mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff));
}

$mysql = new PDO('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'], array(
	PDO::ATTR_PERSISTENT => true,
	PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
));
$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

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

$mysql->exec("INSERT INTO `locations` VALUES(NULL, 'def', 'Default Location')");

// update nodes table
$mysql->exec("ALTER TABLE nodes
	ADD COLUMN location varchar(500) NOT NULL AFTER node,
	ADD COLUMN allocate_memory int(11) NOT NULL AFTER location,
	ADD COLUMN allocate_disk int(11) NOT NULL AFTER allocate_memory,
	ADD COLUMN public int(1) NOT NULL DEFAULT '1' AFTER ports");

$select = $mysql->prepare("SELECT * FROM `nodes`");
$select->execute();

while($row = $select->fetch()){

	$mysql->exec("UPDATE nodes SET `location` = 'def', `allocate_memory` = 1024, `allocate_disk` = 10240, `public` = 1 WHERE `id` = ".$row['id']);

}

$select = $mysql->prepare("SELECT * FROM `servers`");
$select->execute();

while($row = $select->fetch()){

	$mysql->exec("UPDATE servers SET `subusers` = NULL, `gsd_secret` = '".uuid()."' WHERE `id` = ".$row['id']);

}

$select = $mysql->prepare("SELECT * FROM `users`");
$select->execute();

while($row = $select->fetch()){

	$mysql->exec("UPDATE users SET `permissions` = NULL WHERE `id` = ".$row['id']);

}

header('Location: 0.7.5-beta-gsd.php');
