<?php
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	$core->page->redirect('../../../index.php');
}

if(isset($_GET['do']) && $_GET['do'] == 'generate_password')
	exit($core->auth->keygen(rand(12, 18)));
?>
