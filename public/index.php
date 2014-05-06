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
require_once('../src/framework/framework.core.php');

use Aura\Router\RouterFactory;

$router_factory = new RouterFactory;
$router = $router_factory->newInstance();

/*
 * Base Panel Routes
 */
$router->add(null, '/{controller}/{action}{/do}/?')->addValues(array(
	'controller' => 'core'
));

$router->add(null, '/node/{action}{/do}/?')->addValues(array(
	'controller' => 'node'
));

$router->add(null, '/ajax/{action}{/do}/?')->addValues(array(
	'controller' => 'ajax'
));

$base = dirname($_SERVER['PHP_SELF']);
$path = (ltrim($base, '/')) ? substr($_SERVER['REQUEST_URI'], strlen($base)) : parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$path = ($path[0] != '/') ? '/'.$path : $path;

$route = $router->match($path, $_SERVER);

if(!$route) {
    throw new Exception("Error unable to establish a router instance.");
}

if($route->params['controller'] == "core"){

	if(file_exists(dirname(__DIR__).'/panel/'.$route->params['action'].'.php'))
		include(dirname(__DIR__).'/panel/'.$route->params['action'].'.php');
	else	
		throw new Exception("404 File Not Found.");
	
}

if($route->params['controller'] == "ajax"){

	if(file_exists(dirname(__DIR__).'/panel/ajax/'.$route->params['action'].'.php'))
		include(dirname(__DIR__).'/panel/ajax/'.$route->params['action'].'.php');
	else	
		throw new Exception("404 AJAX File Not Found.");

}

if($route->params['controller'] == "node"){

	/*
	 * Handle Authentication
	 */
	if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false)
		Page\components::redirect('core/index');

	if(file_exists(dirname(__DIR__).'/panel/node/'.$route->params['action'].'.php'))
		include(dirname(__DIR__).'/panel/node/'.$route->params['action'].'.php');
	else	
		throw new Exception("404 Node File Not Found.");

}
?>