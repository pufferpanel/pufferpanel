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
	Components\Page::redirect('../../index.php?login');
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(rand(12,18)));

if(!isset($_GET['id']))
	Components\Page::redirect('find.php');

/*
 * Select User Information
 */
$core->user->rebuildData($_GET['id']);

if(!$core->user->getData('id') || $core->user->getData('id') === false)
	Components\Page::redirect('find.php?error=no_user');

$date1 = new \DateTime(date('Y-m-d', $core->user->getData('register_time')));
$date2 = new \DateTime(date('Y-m-d', time()));

$user = $core->user->getData();
$user['register_time'] = date('F j, Y g:ia', $core->user->getData('register_time')).' ('.$date2->diff($date1)->format("%a Days Ago").')';

/*
 * Select Servers Owned by the User
 */
$servers = ORM::forTable('servers')->select('servers.*')->select('nodes.node', 'node_name')
			->join('nodes', array('servers.node', '=', 'nodes.id'))
			->where(array('servers.owner_id' => $core->user->getData('id'), 'servers.active' => 1))
			->findArray();

echo $twig->render(
		'admin/account/view.html', array(
			'user' => $user,
			'servers' => $servers,
			'footer' => array(
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
		));
?>