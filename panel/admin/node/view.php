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

if(!isset($_GET['id'])) {
	Components\Page::redirect('list.php');
}

/*
 * Select Node Information
 */
$node = ORM::forTable('nodes')->findOne($_GET['id']);

if(!$node) {

	if(!isset($_POST['gsd_reset'])) {
		Components\Page::redirect('list.php?error=no_node');
	} else {
		exit('Invalid data provided to function.');
	}

}

if(isset($_POST['gsd_reset'])) {

	if(!$core->gsd->avaliable($node->ip, $node->gsd_listen)) {
		$node->gsd_secret = $core->auth->generateUniqueUUID('nodes', 'gsd_secret');
		$node->save();
		exit($node->gsd_secret);
	} else {
		exit('You must stop GSD before performing this command.');
	}

}


echo $twig->render(
	'admin/node/view.html', array(
		'node' => $node,
		'portlisting' => json_decode($node->ports, true),
		'footer' => array(
			'seconds' => number_format((microtime(true) - $pageStartTime), 4)
		)
	));
?>