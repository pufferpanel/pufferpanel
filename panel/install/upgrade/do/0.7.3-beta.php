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

// upgrader for version 0.7.4 Beta to 0.7.4.1 Beta
require '../../../../src/core/configuration.php';

$mysql = new PDO('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'], array(
	PDO::ATTR_PERSISTENT => true,
	PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
));
$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

// remove old databases
$mysql->exec("DROP TABLE IF EXISTS `acp_email_templates`");
$mysql->exec("DROP TABLE IF EXISTS `modpacks`");

// rename columns on nodes
$mysql->exec("ALTER TABLE nodes CHANGE node_ip fqdn tinytext, sftp_ip ip tinytext");

// add column to servers
$mysql->exec("ALTER TABLE servers ADD COLUMN subusers tinytext AFTER owner_id");

// Update users table
$mysql->exec("ALTER TABLE users MODIFY language char(2) NOT NULL DEFAULT 'en'");
$mysql->exec("ALTER TABLE users ADD COLUMN permissions text AFTER password");
$mysql->exec("ALTER TABLE users DROP COLUMN position");

header('Location: ../0.7.4-beta.php');
?>