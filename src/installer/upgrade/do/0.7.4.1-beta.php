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

// upgrader for version 0.7.4.1 Beta to 0.7.5 Beta
require '../../../../src/core/configuration.php';

$mysql = new PDO('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'], array(
	PDO::ATTR_PERSISTENT => true,
	PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
));
$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

// update acp_settings table
$mysql->exec("ALTER TABLE acp_settings
		ADD COLUMN id int(1) unsigned NOT NULL AUTO_INCREMENT FIRST,
		ADD PRIMARY KEY (id)
	");

// update nodes table
$mysql->exec("ALTER TABLE nodes
		ADD COLUMN daemon_listen int(1) DEFAULT '8003' AFTER daemon_secret,
		ADD COLUMN daemon_console int(1) DEFAULT '8031' AFTER daemon_listen,
		ADD COLUMN daemon_base_dir tinytext AFTER daemon_console
	");

$select = $mysql->prepare("SELECT `id`, `daemon_base_dir` FROM `nodes`");
$select->execute();

while($row = $select->fetch()) {
	
	$mysql->exec("UPDATE nodes SET `daemon_base_dir` = '/home/' WHERE `id` = ".$row['id']);
	
}

header('Location: ../finished.php');
?>
