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
	Page\components::redirect('../../../../index.php');
}

if(!isset($_POST['node'], $_POST['port'], $_POST['ip']))
	exit('POST Only');

/*
 * Verify port is Real & Not in Use
 */
$select = $mysql->prepare("SELECT `ips`, `ports` FROM `nodes` WHERE `id` = :nid");
$select->execute(array(
	':nid' => $_POST['node']
));

if($select->rowCount() != 1)
	exit('Invalid Node');


$data = $select->fetch();

$ips = json_decode($data['ips'], true);
$ports = json_decode($data['ports'], true);

if(array_key_exists($_POST['ip'], $ports) && array_key_exists($_POST['port'], $ports[$_POST['ip']]) && $ports[$_POST['ip']][$_POST['port']] == 1){

	unset($ports[$_POST['ip']][$_POST['port']]);
	$ips[$_POST['ip']]['ports_free'] = ($ips[$_POST['ip']]['ports_free'] - 1);

}else{

	exit('No Port/IP or Port in Use');

}


$update = $mysql->prepare("UPDATE `nodes` SET `ips` = :ips, `ports` = :ports WHERE `id` = :nid");
$update->execute(array(
	':nid' => $_POST['node'],
	':ips' => json_encode($ips),
	':ports' => json_encode($ports)
));

echo 'Done';

?>
