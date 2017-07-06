<?php

/*
	PufferPanel - A Game Server Management Panel
	Copyright (c) 2015 Dane Everitt

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
use \ORM, \PDO, \Twig_Environment, \Twig_Loader_Filesystem, \Tracy\Debugger, \stdClass, \Klein\Klein;
session_start();

if(!ini_get('date.timezone')) {
	date_default_timezone_set('UTC');
}

$baseDir = str_replace('\\', '/', dirname(__FILE__).'/');
if (!file_exists($baseDir.'pufferpanel')) {
    $baseUrl = dirname($_SERVER['PHP_SELF']);
    $baseDir = str_replace('\\', '/', rtrim(dirname(dirname(dirname(__FILE__))), '\\/') . '/pufferpanel' . '/');
}

define('BASE_DIR', $baseDir);
define('APP_DIR', BASE_DIR.'app/');
define('PANEL_DIR', BASE_DIR.'panel/');
define('SRC_DIR', BASE_DIR.'src/');
define('BASE_URL', isset($baseUrl) ? $baseUrl : null);

require_once SRC_DIR.'core/autoloader.php';

$logsDir = Config::config('logsDirectory');
if ($logsDir == null) {
    $logsDir = BASE_DIR.'/logs';
}

define('LOGS_DIR', $logsDir);

/*
 * Handle Cloudflare usage
 */
$_SERVER['REMOTE_ADDR'] = (isset($_SERVER['HTTP_CF_CONNECTING_IP'])) ? $_SERVER['HTTP_CF_CONNECTING_IP'] : $_SERVER['REMOTE_ADDR'];

/*
 * Set Debugger::DETECT to Debugger::DEVELOPMENT to force errors to be displayed.
 * This should NEVER be done on a live environment. In most cases Debugger is smart
 * enough to figure out if it is a local or development environment.
 */
if(Config::config('debugging') === true) {
	Debugger::enable(Debugger::DEVELOPMENT, LOGS_DIR);
} else {
	Debugger::enable(Debugger::PRODUCTION, LOGS_DIR);
}
Debugger::$strictMode = true;

/*
* MySQL PDO Connection Engine
*/
ORM::configure(array(
	'connection_string' => 'mysql:host='.Config::config('mysql')->host.';port='.Config::config('mysql')->port.';dbname='.Config::config('mysql')->database,
	'username' => Config::config('mysql')->username,
	'password' => Config::config('mysql')->password,
	'driver_options' => array(
		PDO::ATTR_PERSISTENT => true,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
	),
	'logging' => false,
	'caching' => true,
	'caching_auto_clear' => true
));

/*
 * Initalize Global core
 */
$core = new stdClass();
$klein = new Klein();

$core->auth        = new Authentication();
$core->user        = new User();
$core->server      = new Server();
$core->email       = new Email();
$core->log         = new Log();
$core->daemon      = new Daemon();
$core->files       = new Files();
$core->twig        = new Twig_Environment(new Twig_Loader_Filesystem(APP_DIR . 'views/'));
$core->language    = new Language();
$core->permissions = new Permissions($core->server->getData('id'));

/*
 * Twig Setup
 */
$core->twig->addGlobal('l', $core->language); // @TODO Change this to addGlobal('language', $core->language) to allow access as {{ language.render('template') }}
$core->twig->addGlobal('settings', Settings::config());
$core->twig->addGlobal('permission', $core->permissions);
$core->twig->addGlobal('fversion', trim(file_get_contents(SRC_DIR.'versions/current')));
$core->twig->addGlobal('admin', (bool) $core->user->getData('root_admin'));
$core->twig->addGlobal('version', Version::get());
$core->twig->addGlobal('basePath', BASE_URL);

$klein->respond('!@^(/auth/|/language/|/api/|/assets/|/oauth2/|/daemon/)', function($request, $response, $service, $app, $klein) use ($core) {

	if(!$core->auth->isLoggedIn()) {

		if(!strpos($request->pathname(), "{{basePath}}/ajax/")) {

			$service->flash('<div class="alert alert-danger">You must be logged in to access that page.</div>');
			$response->redirect('/auth/login');

		} else {

			$response->abort(403);

		}

		$klein->skipRemaining();

	}

});

$klein->respond('@^(/auth/|/oauth2/)', function($request, $response, $service, $app, $klein) use ($core) {

	if($core->auth->isLoggedIn()) {

		// Redirect /auth/* requests to /index if they are logged in
		// Skips redirect on requests to /auth/logout and /auth/remote/*
		if(0 !== strpos($request->pathname(), "{{basePath}}/auth/logout") && 0 !== strpos($request->pathname(), "{{basePath}}/auth/remote/")) {
			$response->redirect('/index');
		}

	}

});

$klein->respond('/node/[*]', function($request, $response, $service, $app, $klein) use($core) {

	if(!$core->auth->isServer()) {

		$response->body($core->twig->render('errors/403.html'))->abort(403);

	}

});

$klein->respond('/admin/[*]', function($request, $response, $service, $app, $klein) use($core) {

	if(!$core->auth->isAdmin()) {

		$response->body($core->twig->render('errors/403.html'))->abort(403);

	}

});

include SRC_DIR.'routes/admin/routes.php';
include SRC_DIR.'routes/ajax/routes.php';
include SRC_DIR.'routes/assets/routes.php';
include SRC_DIR.'routes/auth/routes.php';
include SRC_DIR.'routes/panel/routes.php';
include SRC_DIR.'routes/node/routes.php';
include SRC_DIR.'routes/oauth2/routes.php';
include SRC_DIR.'routes/daemon/routes.php';

$klein->respond(function($request, $response) {
    if(!$response->isAltered()) {
        $response->abort(404);
    }
    if ($response->code() >= 400 || $response->code() < 200) {
        $response->abort($response->code());
    }
});

$klein->onHttpError(function($code, $klein) use ($core) {

    $request = $klein->request();
    $response = $klein->response();

    if(0 !== strpos($request->pathname(), "/daemon/") && 0 !== strpos($request->pathname(), "/daemon/")) {
        return;
    }

    switch($code) {
        case '404': {
            if($request->method('get')) {
                $response->body($core->twig->render('errors/404.html'));
            }
        }
        break;
    }

});

$res = new PPResponse($klein);
$klein->dispatch(null, $res);
