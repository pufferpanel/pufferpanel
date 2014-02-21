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
 * Cloduflare IP Fix
 */
$_SERVER['REMOTE_ADDR'] = (isset($_SERVER['HTTP_CF_CONNECTING_IP'])) ? $_SERVER['HTTP_CF_CONNECTING_IP'] : $_SERVER['REMOTE_ADDR'];

/*
 * Include Dependency Libs
 */
require_once('lib/password.lib.php');

/* 
 * Include Required Global Framework Files
 */
require_once('framework.database.connect.php');
require_once('framework.auth.php');
require_once('framework.page.php');
require_once('framework.files.php');
require_once('framework.settings.php');
require_once('framework.user.php');
require_once('framework.server.php');
require_once('framework.log.php');
require_once('framework.query.php');

/*
 * Include Email Sending Files
 */
require_once('email/core.email.php');

/*
 * Initalize Global Framework
 */
$core = new stdClass();
$core->framework = new stdClass();

/*
 * Initalize Frameworks
 */
$core->framework->settings = new getSettings();
$core->framework->auth = new auth();
$core->framework->user = new user($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), $core->framework->auth->getCookie('pp_server_hash'));
$core->framework->server = new server($core->framework->auth->getCookie('pp_server_hash'), $core->framework->user->getData('id'), $core->framework->user->getData('root_admin'));
$core->framework->email = new tplMail($core->framework->settings);
$core->framework->page = new page($core->framework->user, $core->framework->settings);
$core->framework->log = new log($core->framework->user->getData('id'));
$core->framework->gsd = new GSD_Query($core->framework->server->getData('id'));
$core->framework->files = new files();

/*
 * MySQL PDO Connection Engine
 */
$mysql = dbConn::getConnection();

function pdo_exception_handler($exception) {
    if ($exception instanceof PDOException) {
        
        error_log($exception);
        
        die(json_encode(array('error' => 'A MySQL error was encountered with this request.', 'e_code' => $exception->getCode(), 'e_line' => $exception->getLine(), 'e_time' => date('d-M-Y H:i:s', time()))));
        
    } else {
    
    	die('Exception handler from unknown source.');
    
    }
}
set_exception_handler('pdo_exception_handler');

?>