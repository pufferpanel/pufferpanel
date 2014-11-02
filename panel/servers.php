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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) !== true){
	Components\Page::redirect('index.php?login');
	exit();
}

/*
 * Redirect
 */
if(isset($_GET['goto']) && !empty($_GET['goto']))
	$core->server->nodeRedirect($_GET['goto'], $core->user->getData('id'), $core->user->getData('root_admin'));

/*
 * Get the Servers
 */
if($core->user->getData('root_admin') == '1')
	$servers = ORM::forTable('servers')->select('servers.*')->select('nodes.node', 'node_name')
				->join('nodes', array('servers.node', '=', 'nodes.id'))
				->orderByDesc('active')
				->findArray();
else
	$servers = ORM::forTable('servers')->select('servers.*')->select('nodes.node', 'node_name')
				->join('nodes', array('servers.node', '=', 'nodes.id'))
				->where(array('servers.owner_id' => $core->user->getData('id'), 'servers.active' => 1))
				->where_raw('servers.owner_id = ? OR servers.hash IN(?)', array($core->user->getData('id'), join(',', $core->user->listServerPermissions())))
				->findArray();

/*
 * List Servers
 */
echo $twig->render(
		'panel/servers.html', array(
			'servers' => $servers,
			'footer' => array(
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>
