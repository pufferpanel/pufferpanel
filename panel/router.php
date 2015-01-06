<?php

/*
	PufferPanel - A Minecraft Server Management Panel
	Copyright (c) 2014 PufferPanel

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
use \ORM, \PDO, \Twig_Autoloader, \Twig_Environment, \Twig_Loader_Filesystem, \Tracy\Debugger, \Unirest, \stdClass, \PufferPanel\Core\Config;
session_start();

if(!ini_get('date.timezone')) {
	date_default_timezone_set('UTC');
}

define('BASE_DIR', dirname(__DIR__).'/');
define('APP_DIR', BASE_DIR.'app/');
define('PANEL_DIR', BASE_DIR.'panel/');
define('SRC_DIR', BASE_DIR.'src/');

/*
 * Handle Cloudflare usage
 */
$_SERVER['REMOTE_ADDR'] = (isset($_SERVER['HTTP_CF_CONNECTING_IP'])) ? $_SERVER['HTTP_CF_CONNECTING_IP'] : $_SERVER['REMOTE_ADDR'];

/*
* Has Installer been run?
*/
if(!file_exists(BASE_DIR.'config.json') || (strpos($_SERVER['REQUEST_URI'], '/install') == 0 && !file_exists(SRC_DIR.'install.lock'))) {

	if(file_exists(BASE_DIR.'vendor/autoload.php')) {
		include SRC_DIR.'routes/install/router.php';
	} else {
		throw new Exception("You must install the dependencies before using PufferPanel.");
	}
	return;

}

require_once SRC_DIR.'core/autoloader.php';

Twig_Autoloader::register();
Unirest::timeout(5);

/*
 * Set Debugger::DETECT to Debugger::DEVELOPMENT to force errors to be displayed.
 * This should NEVER be done on a live environment. In most cases Debugger is smart
 * enough to figure out if it is a local or development environment.
 */
Debugger::enable(Debugger::DEVELOPMENT, SRC_DIR.'/logs');
Debugger::$strictMode = TRUE;

/*
* MySQL PDO Connection Engine
*/
ORM::configure(array(
	'connection_string' => 'mysql:host='.Config::getGlobal('mysql')->host.';dbname='.Config::getGlobal('mysql')->database,
	'username' => Config::getGlobal('mysql')->username,
	'password' => Config::getGlobal('mysql')->password,
	'driver_options' => array(
		PDO::ATTR_PERSISTENT => true,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
	)
));

/*
 * Initalize Global core
 */
$core = new stdClass();
$klein = new \Klein\Klein();

/*
 * Initalize cores
 */
$core->settings = new Settings();

/*
* Require HTTPS Connection
*/
if($core->settings->get('https') == 1 && (!isset($_SERVER['HTTPS']) || empty($_SERVER['HTTPS']))) {
	exit(header('Location: https://'.$core->settings->get('master_url').$_SERVER['REQUEST_URI']));
}

$core->auth = new Authentication();
$core->user = new User();
$core->server = new Server();
$core->email = new Email();
$core->log = new Log($core->user->getData('id'));
$core->gsd = new Query($core->server->getData('id'));
$core->files = new Files();
$core->twig = new Twig_Environment(new Twig_Loader_Filesystem(APP_DIR.'views/'), array(
	'cache' => false,
	'debug' => true
));

/*
 * Check Language Settings
 */
if(!$core->user->getData('language')) {
	if(empty($_COOKIE['pp_language'])) {
		$_l = new Language($core->settings->get('default_language'));
	} else {
		$_l = new Language($_COOKIE['pp_language']);
	}
} else {
	$_l = new Language($core->user->getData('language'));
}

/*
 * Twig Setup
 */
$core->twig->addGlobal('lang', $_l->loadTemplates());
$core->twig->addGlobal('settings', $core->settings->get());
$core->twig->addGlobal('get', Components\Page::twigGET());
$core->twig->addGlobal('permission', $core->user->twigListPermissions());
$core->twig->addGlobal('fversion', trim(file_get_contents(SRC_DIR.'versions/current')));
if($core->user->getData('root_admin') == 1) {
	$core->twig->addGlobal('admin', true);
}

$klein->respond('!@^(/auth/|/langauge/|/api/)', function($request, $response, $service, $app, $klein) use ($core) {

	if(!$core->auth->isLoggedIn()) {

		if(!strpos($request->pathname(), "/ajax/")) {

			$service->flash('<div class="alert alert-danger">You must be logged in to access that page.</div>');
			$response->redirect('/auth/login')->send();

		} else {

			$response->code(403);
			$response->body('Not Authenticated.')->send();

		}

		$klein->skipRemaining();

	}

});

$klein->respond('@^/auth/', function($request, $response, $service, $app, $klein) use ($core) {

	if($core->auth->isLoggedIn() && !in_array($request->pathname(), array(
		"/auth/logout",
		"/auth/gsd/download",
		"/auth/gsd/ftp"
	))) {

		$response->redirect('/index')->send();
		$klein->skipRemaining();

	}

});

$klein->respond('/node/[*]', function($request, $response, $service, $app, $klein) use($core) {

	if(!$core->auth->isServer()) {

		$service->flash('<div class="alert alert-danger">You seem to have accessed that page in an invalid manner.</div>');
		$response->redirect('/index')->send();
		$klein->skipRemaining();

	}

});

$klein->respond('/admin/[*]', function($request, $response, $service, $app, $klein) use($core) {

	if(!$core->auth->isAdmin()) {

		$response->redirect('/index')->send();
		$klein->skipRemaining();

	}

});

include SRC_DIR.'routes/admin/routes.php';
include SRC_DIR.'routes/ajax/routes.php';
include SRC_DIR.'routes/auth/routes.php';
include SRC_DIR.'routes/panel/routes.php';
include SRC_DIR.'routes/node/routes.php';
include SRC_DIR.'routes/api/routes.php';

$klein->respond('*', function($request, $response) use ($core) {

	if(!$response->isSent()) {

		$response->code(404);
		$response->body($core->twig->render('errors/404.html'))->send();

	}

});

$klein->dispatch();