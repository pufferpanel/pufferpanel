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
require '../../../../vendor/autoload.php';

use \Unirest;

$mysql = new PDO('mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'], $_INFO['sql_u'], $_INFO['sql_p'], array(
	PDO::ATTR_PERSISTENT => true,
	PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
));
$mysql->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

$select = $mysql->prepare("SELECT * FROM `nodes`");
$select->execute();

while($row = $select->fetch()) {

	// get GSD config
	try {

		$request = Unirest\Request::get(
			"http://".$row['ip'].":".$row['gsd_listen']."/gameservers/config",
			array(
				'X-Access-Token' => $row['gsd_secret']
			)
		);

		if($request->code !== 200) {

			throw new \Exception("HTTP Error attempting to update node.");

		}

		$config = json_decode($request->raw_body, true);

		// update all the core settings
		$config['daemon']['listenport'] = $row['gsd_listen'];
		$config['daemon']['consoleport'] = $row['gsd_console'];
		$config['interfaces']['rest']['authurl'] = str_replace("_ftp", "_download", $config['interfaces']['ftp']['authurl']);
		$config['interfaces']['ftp']['use_ssl'] = true;

		// loop each server and set permissions and keys
		foreach($config['servers'] as $id => $internals) {

			$server = $mysql->prepare("SELECT * FROM `servers` WHERE `node` = :node AND `gsd_id` = :gid");
			$server->execute(array(
				'node' => $row['id'],
				'gid' => $id
			));

			$s = $server->fetch();

			$config['servers'][$id]['build'] = array(
				"install_dir" => '/mnt/MC/CraftBukkit/',
				"disk" => array(
					"hard" => ($s['disk_space'] < 32) ? 32 : (int) $s['disk_space'],
					"soft" => ($s['disk_space'] > 2048) ? (int) $s['disk_space'] - 1024 : 32
				),
				"cpu" => (int) $s['cpu_limit']
			);

			$config['servers'][$id]['keys'] = array(
				$s['gsd_secret'] => array("s:ftp", "s:get", "s:power", "s:files", "s:files:get", "s:files:put", "s:files:zip", "s:query", "s:console", "s:console:send")
			);

			$config['servers'][$id]['gameport'] = (int) $s['server_port'];
			$config['servers'][$id]['gamehost'] = $s['server_ip'];

		}

		$putrequest = Unirest\Request::put(
			"http://".$row['ip'].":".$row['gsd_listen']."/gameservers/config",
			array(
				'X-Access-Token' => $row['gsd_secret']
			),
			array(
				"cfg" => json_encode($config)
			)
		);

	} catch(\Exception $e) {

		exit($e->getMessage());

	}

}

header('Location: ../finished.php');