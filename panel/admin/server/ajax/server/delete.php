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

require_once('../../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true)
	Components\Page::redirect('../../../index.php');

if(!isset($_POST['sid']))
	Components\Page::redirect('../../find.php');

$core->server->rebuildData($_POST['sid']);

$node = ORM::forTable('nodes')->findOne($core->server->getData('node'));

if(!@fsockopen($node->ip, 8003, $num, $error, 3))
	Components\Page::redirect('../../view.php?id='.$_POST['sid'].'&disp=gsd_offline');

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, 'http://'.$node->ip.':8003/gameservers/'.$core->server->getData('gsd_id'));
curl_setopt($ch, CURLOPT_HTTPHEADER, array(
	'X-Access-Token: '.$node->gsd_secret
));
curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 5);
curl_setopt($ch, CURLOPT_TIMEOUT, 5);
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, "DEL");
curl_setopt($ch, CURLOPT_HEADER, false);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_exec($ch);
curl_close($ch);

// delete server record
$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
$server->delete();

// update other servers with higher GSD id on node
$servers = ORM::forTable('servers')->where('node', $core->server->getData('node'))->where_gt('gsd_id', $core->server->getData("gsd_id"))->findMany();
foreach($servers as $s){
	$s->gsd_id--;
	$s->save();
}

// remove all subusers
$json = json_decode($core->server->getData('subusers'), true);

if(is_array($json)){

	foreach($json as $uid => $toss){

		$user = ORM::forTable('users')->findOne($uid);
		$updateJson = json_decode($user->permissions, true);
		unset($updateJson[$core->server->getData('hash')]);
		$user->permissions = json_encode($updateJson);
		$user->save();

	}

}

$ips = json_decode($node->ips, true);
$ports = json_decode($node->ports, true);

$ips[$core->server->getData('server_ip')]['ports_free']++;
$ports[$core->server->getData('server_ip')][$core->server->getData('server_port')]++;

$node->ips = json_encode($ips);
$node->ports = json_encode($ports);

$node->save();

Components\Page::redirect('../../find.php?deletedServer');

?>