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
require_once('../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../index.php?login');
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(rand(12, 18)));

if(!isset($_GET['id']))
	Page\components::redirect('find.php?error=no_id');

$core->server->rebuildData($_GET['id']);
$core->user->rebuildData($core->server->getData('owner_id'));

if(!$core->server->getData('hash') || $core->server->getData('hash') === false)
	Page\components::redirect('find.php?error=invalid_id');

echo $twig->render('admin/server/view.html', array(
		'node' => $core->server->nodeData(),
		'decoded' => array('ips' => json_decode($core->server->nodeData('ips'), true), 'ports' => json_decode($core->server->nodeData('ports'), true)),
		'server' => $core->server->getData(),
		'user' => $core->user->getData(),
		'footer' => array(
			'queries' => Database\databaseInit::getCount(),
			'seconds' => number_format((microtime(true) - $pageStartTime), 4)
		)
	));
?>