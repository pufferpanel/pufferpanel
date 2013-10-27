<?php
session_start();
require_once('../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}

if(!isset($_POST['pub_key'], $_POST['priv_key']))
	$core->framework->page->redirect('../global.php?error=pub_key|priv_key');
	
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'captcha_pub'")->execute(array($_POST['pub_key']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'captcha_priv'")->execute(array($_POST['priv_key']));

$core->framework->page->redirect('../global.php');
?>