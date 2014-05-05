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
require_once(dirname(__DIR__).'/src/framework/framework.core.php');

//if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) === true)
//	Page\components::redirect('servers.php');

$klein = new \Klein\Klein();

$klein->respond('GET', '/[index|password:action]?', function($request, $response) use ($core, $twig) {

	return $twig->render(
			'panel/'.$request->action.'.html', array(
				'footer' => array(
					'queries' => Database\databaseInit::getCount(),
					'seconds' => number_format((microtime(true) - $core->startTime), 4)
				)
		));

});

$klein->respond('GET', '/servers/[:goto]?', function($request, $response) use ($core, $twig, $mysql) {
	
	if($request->goto){
	
		$core->server->nodeRedirect($request->goto, $core->user->getData('id'), $core->user->getData('root_admin'));
	
	}else{
	
		include('../panel/servers.php');
		return $twig->render(
			'panel/servers.html', array(
				'servers' => array($servers),
				'footer' => array(
					'queries' => Database\databaseInit::getCount(),
					'seconds' => number_format((microtime(true) - $core->startTime), 4)
				)
		));
		
	}
	
});

$klein->dispatch();