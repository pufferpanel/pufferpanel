<?php
session_start();
require_once('core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token')) !== true){
	
	$core->framework->page->redirect('index.php');
	exit();
	
}else{

	/*
	 * Expire Session and Clear Database Details
	 */
	setcookie("pp_auth_token", null, time()-86400, '/', $core->framework->settings->get('cookie_website'));
	setcookie("pp_server_node", null, time()-86400, '/', $core->framework->settings->get('cookie_website'));
	
	$logoutUser = $mysql->prepare("UPDATE `users` SET `session_id` = NULL, `session_ip` = NULL, `session_expires` = NULL WHERE `session_ip` = :sesip AND `session_id` = :sesid");
	$logoutUser->execute(array(':sesip' => $_SERVER['REMOTE_ADDR'], ':sesid' => $_COOKIE['pp_auth_token']));
	
	$core->framework->page->redirect('index.php');
	exit();

}
?>