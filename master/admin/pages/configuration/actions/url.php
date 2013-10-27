<?php
session_start();
require_once('../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}

setcookie("__TMP_pp_admin_updateglobal", json_encode($_POST), time() + 10, '/', $core->framework->settings->get('cookie_website'));

if(!isset($_POST['main_url'], $_POST['master_url'], $_POST['assets_url']))
	$core->framework->page->redirect('../global.php?error=main_url|master_url|assets_url');
	
foreach($_POST as $id => $val)
	{
	
		if(!filter_var($val, FILTER_VALIDATE_URL))
			$core->framework->page->redirect('../global.php?error='.$id);
	
	}

$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'main_website'")->execute(array($_POST['main_url']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'master_url'")->execute(array($_POST['master_url']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'assets_url'")->execute(array($_POST['assets_url']));

$core->framework->page->redirect('../global.php');

?>