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

require_once('../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Components\Page::redirect('../../../index.php');
}

if(!isset($_POST['node'], $_POST['port'], $_POST['ip']))
	exit('POST Only');

/*
 * Verify port is Real & Not in Use
 */
$node = ORM::forTable('nodes')->findOne($_POST['node']);

if($node === false)
	exit('Invalid Node');

$ips = json_decode($node->ips, true);
$ports = json_decode($node->ports, true);

if(array_key_exists($_POST['ip'], $ports) && array_key_exists($_POST['port'], $ports[$_POST['ip']]) && $ports[$_POST['ip']][$_POST['port']] == 1){

	unset($ports[$_POST['ip']][$_POST['port']]);
	$ips[$_POST['ip']]['ports_free'] = ($ips[$_POST['ip']]['ports_free'] - 1);

}else
	exit('No Port/IP or Port in Use');

$node->ips = json_encode($ips);
$node->ports = json_encode($ports);
$node->save();

die('Done');

?>