<?php
session_start();
require_once('../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../index.php');
}

if(!isset($_POST['company_name']))
	$core->framework->page->redirect('../global.php?error=company_name&disp=cname');
	
$query = $mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'company_name'");
$query->execute(array($_POST['company_name']));

$core->framework->page->redirect('../global.php');

?>