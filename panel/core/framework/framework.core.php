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
 * PufferPanel Core Framework File
 */
$pageStartTime = microtime(true);

/*
 * Cloudflare IP Fix
 */
$_SERVER['REMOTE_ADDR'] = (isset($_SERVER['HTTP_CF_CONNECTING_IP'])) ? $_SERVER['HTTP_CF_CONNECTING_IP'] : $_SERVER['REMOTE_ADDR'];

/*
 * Include Dependency Libs
 */
require_once(dirname(dirname(dirname(__DIR__))).'/vendor/autoload.php');
require_once('lib/password.lib.php');

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
require_once('framework.ssh2.php');
require_once('framework.log.php');
require_once('framework.query.php');
require_once('framework.language.php');
require_once('framework.functions.php');

/*
 * Include Email Sending Files
 */
require_once('email/core.email.php');

/*
 * Initalize Global Framework
 */
$core = new stdClass();
$_l = new stdClass();
set_exception_handler('pdo_exception_handler');

/*
 * Initalize Frameworks
 */
$core->settings = new settings();
$core->auth = new \Auth\auth($core->settings);
$core->ssh = new ssh($core->settings->get('use_ssh_keys'));
$core->user = new user($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash'));
$core->server = new server($core->auth->getCookie('pp_server_hash'), $core->user->getData('id'), $core->user->getData('root_admin'));
$core->email = new tplMail($core->settings);
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


function pdo_exception_handler($exception) {
    
    if($exception instanceof PDOException) {
        
        error_log($exception);
        exit('<!DOCTYPE html>
        <html lang="en">
        <head>
        	<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">
        </head>
        <body>
        	<div class="container">
        		<h1>Database Error</h1>
        			<div class="col-12">
        				<div class="alert alert-danger"><strong>Error:</strong> An unexpected MySQL Error was encountered with this request. Please try again in a few minutes.</div>
        			</div>
        	</div>
        </body>
        </html>');
        
    } else {
    
       	exit('<!DOCTYPE html>
    	<html lang="en">
    	<head>
    		<link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">
    	</head>
    	<body>
    		<div class="container">
    			<h1>Error</h1>
    				<div class="col-12">
    					<div class="alert alert-danger">'.nl2br($exception).'</div>
    				</div>
    		</div>
    	</body>
    	</html>');
    
    }
}

Twig_Autoloader::register();

$loader = new Twig_Loader_Filesystem(dirname(dirname(dirname(__DIR__))).'/app/views/');
$twig = new Twig_Environment($loader, array(
    'cache' => false,
    'debug' => true
));
$twig->addGlobal('lang', $_l->loadTemplates());
$twig->addGlobal('settings', $core->settings->get());
$twig->addGlobal('get', Page\components::twigGET());
if($core->user->getData('root_admin') == 1){ $twig->addGlobal('admin', true); }
?>
