<?php
session_start();
require_once('../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}

setcookie("__TMP_pp_admin_updateglobal", json_encode($_POST), time() + 10, '/', $core->framework->settings->get('cookie_website'));

if(!isset($_POST['smail_method'], $_POST['sendmail_email'], $_POST['postmark_api_key'], $_POST['mandrill_api_key'], $_POST['mailgun_api_key']))
	$core->framework->page->redirect('../global.php?error=smail_method|sendmail_email|postmark_api_key|mandrill_api_key|mailgun_api_key');
	
if(!in_array($_POST['smail_method'], array('php', 'postmark', 'mandrill', 'mailgun')))
	$core->framework->page->redirect('../global.php?error=smail_method');
	
if(!filter_var($_POST['sendmail_email'], FILTER_VALIDATE_EMAIL))
	$core->framework->page->redirect('../global.php?error=sendmail_email');
	
if($_POST['smail_method'] != 'php' && empty($_POST[$_POST['smail_method'].'_api_key']))
	$core->framework->page->redirect('../global.php?error=smail_method|'.$_POST['smail_method'].'_api_key');
	
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'sendmail_method'")->execute(array($_POST['smail_method']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'sendmail_email'")->execute(array($_POST['sendmail_email']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'postmark_api_key'")->execute(array($_POST['postmark_api_key']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'mandrill_api_key'")->execute(array($_POST['mandrill_api_key']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'mailgun_api_key'")->execute(array($_POST['mailgun_api_key']));

$core->framework->page->redirect('../global.php');

?>