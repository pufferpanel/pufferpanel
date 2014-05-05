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
require_once('../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Page\components::redirect($core->settings->get('master_url').'index.php?login');
	exit();
	
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(rand(12, 18)));


$packs = $mysql->prepare("SELECT hash, name, version, deleted FROM `modpacks` WHERE `hash` = :hash");
$packs->execute(array(
	':hash' => $core->server->getData('modpack')
));

$pack = $packs->fetch();
$isDeleted = ($pack['deleted'] == 1) ? '[DELETED] ' : null;
$currentPack = $isDeleted.$pack['name'].' ('.$pack['version'].')';


$packs = $mysql->prepare("SELECT hash, name, version, server_jar FROM `modpacks` WHERE `deleted` = 0 AND `min_ram` <= :ram");
$packs->execute(array(
	':ram' => $core->server->getData('max_ram')
));

$allpacks = array();
while($row = $packs->fetch()){
	
	$allpacks = array_merge($allpacks, array($row));
	
}

/*
 * Display Page
 */
echo $twig->render(
		'node/settings.html', array(
			'server' => array(
				'server_jar' => (str_replace(".jar", "", $core->server->getData('server_jar'))),
				'ftp_user' => $core->server->getData('ftp_user'),
				'modpack' => $core->server->getData('modpack')
			),
			'node' => array(
				'sftp_ip' => $core->server->nodeData('sftp_ip')
			),
			'modpack' => array(
				'current' => $currentPack,
				'packs' => $allpacks
			),
			'footer' => array(
				'queries' => Database\databaseInit::getCount(),
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>