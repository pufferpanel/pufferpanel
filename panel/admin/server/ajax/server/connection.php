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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Components\Page::redirect('../../../index.php');
}

if(!isset($_POST['sid']))
	Components\Page::redirect('../../find.php');

$_POST['server_port'] = $_POST['server_port_'.str_replace('.', '_', $_POST['server_ip'])];

if(!isset($_POST['server_ip'], $_POST['server_port'], $_POST['nid']))
	Components\Page::redirect('../../view.php?id='.$_POST['sid']);

$core->server->rebuildData($_POST['sid']);
$core->user->rebuildData($core->server->getData('owner_id'));

$ports = json_decode($core->server->nodeData('ports'), true);
$ips = json_decode($core->server->nodeData('ips'), true);

if(!array_key_exists($_POST['server_ip'], $ports))
	Components\Page::redirect('../../view.php?id='.$_POST['sid'].'&error=server_ip&disp=no_ip');

if(!array_key_exists($_POST['server_port'], $ports[$_POST['server_ip']]))
	Components\Page::redirect('../../view.php?id='.$_POST['sid'].'&error=server_port&disp=no_port');

if($ports[$_POST['server_ip']][$_POST['server_port']] == 0 && $_POST['server_port'] != $core->server->getData('server_port'))
	Components\Page::redirect('../../view.php?id='.$_POST['sid'].'&error=server_port&disp=port_in_use');

$mysql->prepare("UPDATE `servers` SET `server_ip` = :ip, `server_port` = :port WHERE `id` = :sid")->execute(array(
	':ip' => $_POST['server_ip'],
	':port' => $_POST['server_port'],
	':sid' => $_POST['sid']
));

/*
 * Update Old
 */
$ports[$core->server->getData('server_ip')][$core->server->getData('server_port')] = 1;
$ips[$core->server->getData('server_ip')]['ports_free']++;

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

/*
* Build the Data
*/
$url = "http://".$core->server->nodeData('ip').":8003/gameservers/".$core->server->getData('gsd_id');
$data = json_encode(array(
	"gameport" => (int)$_POST['server_port'],
	"gamehost" => $_POST['server_ip']
));

/*
* Run Query Aganist GSD
*/
$curl = curl_init($url);
curl_setopt($curl, CURLOPT_CUSTOMREQUEST, "PUT");
curl_setopt($curl, CURLOPT_HTTPHEADER, array(
	'X-Access-Token: '.$core->server->nodeData('gsd_secret')
));
curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($data));

Components\Page::redirect('../../view.php?id='.$_POST['sid']);
