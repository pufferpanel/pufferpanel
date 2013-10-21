<?php
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../../index.php');
}

if(!isset($_POST['uid']) || !is_numeric($_POST['uid']))
	$core->framework->page->redirect('../find.php?error=UPDATE-USER__undefined_user');

if($_POST['action'] == 'details'){	
	
	$update = $mysql->prepare("UPDATE `users` SET `email` = :email, `root_admin` = :root WHERE `id` = :uid");
	$update->execute(array(
		':email' => $_POST['email'],
		':root' => $_POST['root_admin'],
		':uid' => $_POST['uid']
	));
	
	$core->framework->page->redirect('../../view.php?id='.$_POST['uid'].'&disp=d_updated');
	
}

?>