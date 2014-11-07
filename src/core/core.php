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
session_start();

/*
 * Initalize Start Time
 */
$pageStartTime = microtime(true);

/*
 * Cloudflare IP Fix
 */
$_SERVER['REMOTE_ADDR'] = (isset($_SERVER['HTTP_CF_CONNECTING_IP'])) ? $_SERVER['HTTP_CF_CONNECTING_IP'] : $_SERVER['REMOTE_ADDR'];

/*
* Has Installer been run?
*/
if(!file_exists(__DIR__.'/configuration.php'))
	exit("Installer has not yet been run. Please navigate to the installer and run through the steps to use this software.");

/*
 * Define Directories
 */
define('BASE_DIR', dirname(dirname(__DIR__)).'/');
define('APP_DIR', BASE_DIR.'app/');
define('PANEL_DIR', BASE_DIR.'panel/');
define('SRC_DIR', BASE_DIR.'src/');

/*
 * Include Dependency Libs
 */
require_once(SRC_DIR.'core/configuration.php');
require_once(BASE_DIR.'vendor/autoload.php');

/*
 * To use a local-only debugging option please uncomment the lines
 * below and comment out the bugsnag lines above. This debugging can be
 * used on a live environment if you wish.
 */
use Tracy\Debugger;
Debugger::enable(Debugger::DETECT, dirname(__DIR__).'/logs');
Debugger::$strictMode = TRUE;

/*
* MySQL PDO Connection Engine
*/
use \ORM as ORM;
use \PDO as PDO;
ORM::configure(array(
	'connection_string' => 'mysql:host='.$_INFO['sql_h'].';dbname='.$_INFO['sql_db'],
	'username' => $_INFO['sql_u'],
	'password' => $_INFO['sql_p'],
	'driver_options' => array(
		PDO::ATTR_PERSISTENT => true,
		PDO::ATTR_DEFAULT_FETCH_MODE => PDO::FETCH_ASSOC
	)
));

/*
* Include Required Global Component Files
*/
require_once(__DIR__.'/components/authentication.php');
require_once(__DIR__.'/components/functions.php');
require_once(__DIR__.'/components/page.php');

/*
 * Include Required Global Class Files
 */
require_once(__DIR__.'/authentication.php');
require_once(__DIR__.'/email.php');
require_once(__DIR__.'/files.php');
require_once(__DIR__.'/language.php');
require_once(__DIR__.'/user.php');
require_once(__DIR__.'/log.php');
require_once(__DIR__.'/query.php');
require_once(__DIR__.'/server.php');
require_once(__DIR__.'/settings.php');

/*
 * Initalize Global core
 */
$core = new \stdClass();
$_l = new \stdClass();

/*
 * Initalize cores
 */
$core->settings = new Settings();
$core->auth = new Authentication();
$core->user = new User($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash'));
$core->server = new Server($core->auth->getCookie('pp_server_hash'), $core->user->getData('id'), $core->user->getData('root_admin'));
$core->email = new Email();
$core->log = new Log($core->user->getData('id'));
$core->gsd = new Query($core->server->getData('id'));
$core->files = new Files();

/*
 * Require HTTPS Connection
 */
if($core->settings->get('https') == 1)
	if(!isset($_SERVER['HTTPS']) || $_SERVER['HTTPS'] == "") {
		http_response_code(302);
		header("Location: https://".$_SERVER['HTTP_HOST'].$_SERVER['REQUEST_URI']);
	}

/*
 * Check Language Settings
 */
if($core->user->getData('language') === false)
	if(!isset($_COOKIE['pp_language']) || empty($_COOKIE['pp_language']))
		$_l = new Language($core->settings->get('default_language'));
	else
		$_l = new Language($_COOKIE['pp_language']);
else
	$_l = new Language($core->user->getData('language'));

/*
 * Twig Setup
 */
\Twig_Autoloader::register();
$loader = new \Twig_Loader_Filesystem(APP_DIR.'views/');
$twig = new \Twig_Environment($loader, array(
	'cache' => false,
	'debug' => true
));
$twig->addGlobal('lang', $_l->loadTemplates());
$twig->addGlobal('settings', $core->settings->get());
$twig->addGlobal('get', Components\Page::twigGET());
$twig->addGlobal('permission', $core->user->twigListPermissions());
$twig->addGlobal('fversion', trim(file_get_contents(SRC_DIR.'versions/current')));
if($core->user->getData('root_admin') == 1){ $twig->addGlobal('admin', true); }
?>