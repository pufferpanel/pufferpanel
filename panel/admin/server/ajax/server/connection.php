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
require_once('../../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php');
}

if(!isset($_POST['sid']))
	Page\components::redirect('../../find.php');

$_POST['server_port'] = $_POST['server_port_'.str_replace('.', '_', $_POST['server_ip'])];

if(!isset($_POST['server_ip'], $_POST['server_port'], $_POST['nid']))
	Page\components::redirect('../../view.php?id='.$_POST['sid']);

$select = $mysql->prepare("SELECT `ports`, `ips` FROM `nodes` WHERE `id` = :nid");
$select->execute(array(':nid' => $_POST['nid']));
$node = $select->fetch();

$select = $mysql->prepare("SELECT `server_port`, `server_ip` FROM `servers` WHERE `id` = :sid");
$select->execute(array(':sid' => $_POST['sid']));
$server = $select->fetch();

$ports = json_decode($node['ports'], true);
$ips = json_decode($node['ips'], true);

if(!array_key_exists($_POST['server_ip'], $ports))
	Page\components::redirect('../../view.php?id='.$_POST['sid'].'&error=server_ip&disp=no_ip');

if(!array_key_exists($_POST['server_port'], $ports[$_POST['server_ip']]))
	Page\components::redirect('../../view.php?id='.$_POST['sid'].'&error=server_port&disp=no_port');

if($ports[$_POST['server_ip']][$_POST['server_port']] == 0 && $_POST['server_port'] != $server['server_port'])
	Page\components::redirect('../../view.php?id='.$_POST['sid'].'&error=server_port&disp=port_in_use');

$mysql->prepare("UPDATE `servers` SET `server_ip` = :ip, `server_port` = :port WHERE `id` = :sid")->execute(array(
	':ip' => $_POST['server_ip'],
	':port' => $_POST['server_port'],
	':sid' => $_POST['sid']
));

/*
 * Update Old
 */
$ports[$server['server_ip']][$server['server_port']] = 1;
$ips[$server['server_ip']]['ports_free']++;

/*
 * Update Old
 */
$ports[$_POST['server_ip']][$_POST['server_port']] = 0;
$ips[$_POST['server_ip']]['ports_free']--;

$mysql->prepare("UPDATE `nodes` SET `ports` = :ports, `ips` = :ips WHERE `id` = :nid")->execute(array(
	':ports' => json_encode($ports),
	':ips' => json_encode($ips),
	':nid' => $_POST['nid']
));

Page\components::redirect('../../view.php?id='.$_POST['sid']);
