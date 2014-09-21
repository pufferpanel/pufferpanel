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

if(isset($_GET['do']) && $_GET['do'] == 'redirect' && isset($_GET['node'])){

	$select = $mysql->prepare("SELECT `id` FROM `nodes` WHERE `node` = :name");
	$select->execute(array(':name' => $_GET['node']));

	if($select->rowCount() == 1) {
		$n = $select->fetch();
		Page\components::redirect('view.php?id='.$n['id']);
	}else{
		Page\components::redirect('list.php');
	}

}

if(!isset($_GET['id']))
	Page\components::redirect('list.php');

/*
 * Select Node Information
 */
$selectNode = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
$selectNode->execute(array(
	':id' => $_GET['id']
));

	if($selectNode->rowCount() != 1)
		Page\components::redirect('list.php?error=no_node');
	else
		$node = $selectNode->fetch();

echo $twig->render(
	'admin/node/view.html', array(
		'node' => $node,
		'portlisting' => json_decode($node['ports'], true),
		'footer' => array(
			'queries' => Database\databaseInit::getCount(),
			'seconds' => number_format((microtime(true) - $pageStartTime), 4)
		)
	));