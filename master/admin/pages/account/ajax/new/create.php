<?php
session_start();
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../../index.php');
}

if(!preg_match('/^[\w-]{4,35}$/', $_POST['username']))
	$core->framework->page->redirect('../../new.php?disp=u_fail');
	
if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL))
	$core->framework->page->redirect('../../new.php?disp=e_fail');
	
if(strlen($_POST['pass']) < 8 || $_POST['pass'] != $_POST['pass_2'])
	$core->framework->page->redirect('../../new.php?disp=p_fail');

$query = $mysql->prepare("SELECT * FROM `users` WHERE `username` = :user OR `email` = :email");
$query->execute(array(
	':user' => $_POST['username'],
	':email' => $_POST['email']
));

if($query->rowCount() > 0)
	$core->framework->page->redirect('../../new.php?disp=a_fail');

$insert = $mysql->prepare("INSERT INTO `users` VALUES(NULL, NULL, :user, :email, :pass, :time, 'owner', NULL, NULL, NULL, 0, 0, 0)");
$insert->execute(array(
	':user' => $_POST['username'],
	':email' => $_POST['email'],
	':pass' => $core->framework->auth->encrypt($_POST['pass']),
	':time' => time()
));

/*
 * Send Email
 */
$message = $core->framework->email->adminAccountCreated(array('PASS' => $_POST['pass'], 'EMAIL' => $_POST['email']));
$core->framework->email->dispatch($_POST['email'], $core->framework->settings->get('company_name').' - Account Created', $message);

$core->framework->page->redirect('../../view.php?id='.$mysql->lastInsertId());

?>