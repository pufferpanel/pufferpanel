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
session_start();
require_once('../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) !== true){
	Page\components::redirect('index.php?login');
	exit();
}

/*
 * Redirect
 */
if(isset($_GET['goto']) && !empty($_GET['goto']))
	$core->server->nodeRedirect($_GET['goto'], $core->user->getData('id'), $core->user->getData('root_admin'));

/*
 * Get the Servers
 */
if($core->user->getData('root_admin') == '1'){
	$query = $mysql->prepare("SELECT * FROM `servers` ORDER BY `active` DESC");
	$query->execute();
}else{
	$query = $mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :oid ORDER BY `active` DESC");
	$query->execute(array(':oid' => $core->user->getData('id')));
}

/*
 * Build Array
 */
$servers = array();
while($row = $query->fetch()){

	$servers = array_merge($servers, array(
		"id" => $row['id'],
		"hash" => $row['hash'],
		"node" => $core->settings->nodeName($row['node']),
		"server_ip" => $row['server_ip'],
		"server_port" => $row['server_port'],
		"name" => $row['name']
	));

}

/*
 * List Servers
 */
echo $twig->render(
		'panel/servers.html', array(
			'servers' => array($servers),
			'footer' => array(
				'queries' => Database\databaseInit::getCount(),
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>