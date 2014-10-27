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
namespace PufferPanel\Core;
use \ORM as ORM;

require_once('../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Components\Page::redirect('../../index.php?login');
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(rand(12,18)));

if(!isset($_GET['id']))
	Components\Page::redirect('find.php');

/*
 * Select User Information
 */
$select = $mysql->prepare("SELECT * FROM `users` WHERE `id` = :id LIMIT 1");
$select->execute(array(
	':id' => $_GET['id']
));

	if($select->rowCount() != 1)
		Components\Page::redirect('find.php?error=no_user');
	else
		$user = $select->fetch();

$date1 = new DateTime(date('Y-m-d', $user['register_time']));
$date2 = new DateTime(date('Y-m-d', time()));

$user['register_time'] = date('F j, Y g:ia', $user['register_time']).' ('.$date2->diff($date1)->format("%a Days Ago").')';

/*
 * Select Servers Owned by the User
 */
$select = $mysql->prepare("SELECT * FROM `servers` WHERE `owner_id` = :id ORDER BY `active` DESC");
$select->execute(array(
	':id' => $user['id']
));

/*
 * Iterate through Servers
 */
$servers = array();
while($row = $select->fetch()){

	$servers = array_merge($servers, array(array(
		"id" => $row['id'],
		"hash" => $row['hash'],
		"node" => $core->settings->nodeName($row['node']),
		"ip" => $row['server_ip'],
		"port" => $row['server_port'],
		"name" => $row['name'],
		"active" => ($row['active'] == '1') ? true : false,
	)));

}

echo $twig->render(
		'admin/account/view.html', array(
			'user' => $user,
			'servers' => $servers,
			'footer' => array(
				'queries' => Database_Initiator::getCount(),
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
		));
?>