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

function gen_UUID(){

	return sprintf('%04x%04x-%04x-%04x-%04x-%04x%04x%04x', mt_rand(0, 0xffff), mt_rand(0, 0xffff),
					mt_rand(0, 0xffff),
					mt_rand(0, 0x0fff) | 0x4000,
					mt_rand(0, 0x3fff) | 0x8000,
					mt_rand(0, 0xffff), mt_rand(0, 0xffff), mt_rand(0, 0xffff));

}

// update acp_settings table
$mysql->exec("ALTER TABLE users
				ADD COLUMN uuid varchar(36) NOT NULL AFTER whmcs_id
			");

// set a uuid for each user
$results = $mysql->prepare("SELECT * FROM `users`");
$results->execute();

$row = $results->fetch();
foreach($row as &$row){

	$mysql->exec("UPDATE `users` SET `uuid` = '".gen_UUID()."' WHERE `id` = ".$row['id']);

}

header('Location: 0.7.4.1-beta.php');
?>