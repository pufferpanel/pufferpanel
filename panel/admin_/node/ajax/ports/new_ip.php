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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true) {

	Components\Page::redirect('../../../index.php');

}

if(!isset($_POST['ip_port'])) {

	Components\Page::redirect('../../view.php?id='.$_POST['nid'].'&disp=missing_args');

}

/*
 * Clean Inputs
 */
$_POST['ip_port'] = str_replace(" ", "", $_POST['ip_port']);

$lines = explode("\r\n", $_POST['ip_port']);

/*
 * Unpack Variables
 */
$node = ORM::forTable('nodes')->findOne($_POST['nid']);

$IPA = json_decode($node->ips, true);
$IPP = json_decode($node->ports, true);

foreach($lines as $id => $values) {

	list($ip, $ports) = explode('|', $values);

	$IPA = array_merge($IPA, array($ip => array()));
	$IPP = array_merge($IPP, array($ip => array()));

	$portList = [];
	$portList = Components\Functions::processPorts($ports);

	for($l=0; $l<count($portList); $l++) {

			$IPP[$ip][$portList[$l]] = 1;

	}

	/*
	 * Make sure Ports are in the array
	 */
	if(count($IPP[$ip]) > 0) {

		$IPA[$ip] = array_merge($IPA[$ip], array("ports_free" => count($IPP[$ip])));

	} else {

		Components\Page::redirect('../../view.php?id='.$_POST['nid'].'&error=ip_port&disp=ip_port_space');

	}

}

$node->ips = json_encode($IPA);
$node->ports = json_encode($IPP);
$node->save();

Components\Page::redirect('../../view.php?id='.$_POST['nid']);
?>
