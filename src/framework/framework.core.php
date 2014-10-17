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

/*
 * Initalize Start Time
 */
$pageStartTime = microtime(true);

/*
 * Cloudflare IP Fix
 */
$_SERVER['REMOTE_ADDR'] = (isset($_SERVER['HTTP_CF_CONNECTING_IP'])) ? $_SERVER['HTTP_CF_CONNECTING_IP'] : $_SERVER['REMOTE_ADDR'];

/*
 * Define Directories
 */
define('BASE_DIR', dirname(dirname(__DIR__)).'/');
define('APP_DIR', dirname(dirname(__DIR__)).'/app/');
define('PANEL_DIR', dirname(dirname(__DIR__)).'/panel/');
define('SRC_DIR', dirname(dirname(__DIR__)).'/src/');

/*
 * Include Dependency Libs
 */
require_once(BASE_DIR.'vendor/autoload.php');

/*
* Bugsnag Settings
*
* Bugsnag is an Cloud-based Error collection system. This system allows
* us to collect error information directly from your panel when they
* occur. This allows us to patch errors that you may not notice, or that
* are hard to track down. To enable this, please uncomment the lines below.
*/
$bugsnag = new Bugsnag_Client("0745694bcbeaab273d7e0095e7947a0d");
$bugsnag->setAppVersion(trim(file_get_contents(SRC_DIR.'versions/current')));
$bugsnag->setMetaData(array('git' => substr(trim(file_get_contents(BASE_DIR.'.git/refs/heads/master')), 0, 8)));
$bugsnag->setFilters(array('pass', 'password', 'hash', 'token')); // Prevent accidental sending of this (filters anything containing that string)
set_error_handler(array($bugsnag, "errorHandler"));
set_exception_handler(array($bugsnag, "exceptionHandler"));

/*
 * To use a local-only debugging option please uncomment the lines
 * below and comment out the bugsnag lines above. This debugging can be
 * used on a live environment if you wish.
 */
//use Tracy\Debugger;
//Debugger::enable(Debugger::DETECT, dirname(__DIR__).'/logs');
//Debugger::$strictMode = TRUE;

/*
 * Has Installer been run?
 */
if(!file_exists(__DIR__.'/configuration.php'))
	exit("Installer has not yet been run. Please navigate to the installer and run through the steps to use this software.");

/*
 * Include Required Global Framework Files
 */
require_once('framework.database.connect.php');
require_once('framework.page.php');
require_once('framework.auth.php');
require_once('framework.files.php');
require_once('framework.user.php');
require_once('framework.server.php');
require_once('framework.settings.php');
require_once('framework.log.php');
require_once('framework.query.php');
require_once('framework.language.php');
require_once('framework.functions.php');
require_once('framework.email.php');

/*
 * Initalize Global Framework
 */
$core = new stdClass();
$_l = new stdClass();

/*
 * Initalize Frameworks
 */
$core->settings = new settings();
$core->auth = new \Auth\auth();
$core->user = new user($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash'));
$core->server = new server($core->auth->getCookie('pp_server_hash'), $core->user->getData('id'), $core->user->getData('root_admin'));
$core->email = new tplMail();
$core->log = new log($core->user->getData('id'));
$core->gsd = new query($core->server->getData('id'));
$core->files = new files();

/*
 * Check Language Settings
 */
if($core->user->getData('language') === false)
	if(!isset($_COOKIE['pp_language']) || empty($_COOKIE['pp_language']))
		$_l = new Language\lang($core->settings->get('default_language'));
	else
		$_l = new Language\lang($_COOKIE['pp_language']);
else
	$_l = new Language\lang($core->user->getData('language'));

/*
 * MySQL PDO Connection Engine
 */
$mysql = Database\database::connect();

/*
 * Twig Setup
 */
Twig_Autoloader::register();

$loader = new Twig_Loader_Filesystem(APP_DIR.'views/');
$twig = new Twig_Environment($loader, array(
	'cache' => false,
	'debug' => true
));
$twig->addGlobal('lang', $_l->loadTemplates());
$twig->addGlobal('settings', $core->settings->get());
$twig->addGlobal('get', Page\components::twigGET());
$twig->addGlobal('fversion', trim(file_get_contents(dirname(__DIR__).'/versions/current')));
if($core->user->getData('root_admin') == 1){ $twig->addGlobal('admin', true); }
?>
