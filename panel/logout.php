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
require_once('../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token')) !== true){
	
	Page\components::redirect('index.php?login');
	exit();
	
}else{

	/*
	 * Expire Session and Clear Database Details
	 */
	setcookie("pp_auth_token", null, time()-86400, '/', $core->settings->get('cookie_website'));
	setcookie("pp_server_node", null, time()-86400, '/', $core->settings->get('cookie_website'));
	setcookie("pp_server_hash", null, time()-86400, '/', $core->settings->get('cookie_website'));
	
	$logoutUser = $mysql->prepare("UPDATE `users` SET `session_id` = NULL, `session_ip` = NULL WHERE `session_ip` = :sesip AND `session_id` = :sesid");
	$logoutUser->execute(array(':sesip' => $_SERVER['REMOTE_ADDR'], ':sesid' => $_COOKIE['pp_auth_token']));
	
    $core->log->getUrl()->addLog(0, 1, array('auth.user_logout', 'Account logged out from '.$_SERVER['REMOTE_ADDR'].'.'));
	Page\components::redirect('index.php?login');
	exit();

}
?>
