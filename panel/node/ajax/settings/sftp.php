<?php
session_start();
require_once('../../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	$core->page->redirect('../../../index.php');
}
	
if(!isset($_POST['sftp_pass'], $_POST['sftp_pass_2'], $_POST['nid']))
	$core->page->redirect('../../settings.php');
	
if(strlen($_POST['sftp_pass']) < 8)
	$core->page->redirect('../../settings.php?error=sftp_pass|sftp_pass_2&disp=pass_len');
	
if($_POST['sftp_pass'] != $_POST['sftp_pass_2'])
	$core->page->redirect('../../settings.php?error=sftp_pass|sftp_pass_2&disp=pass_match');

/* 
 * Select Node, User, & Server Information
 */
$select = $mysql->prepare("SELECT `ftp_user`, `name`, `owner_id`, `node` FROM `servers` WHERE `id` = :id");
$select->execute(array(
	':id' => $core->server->getData('id')
));
    $server = $select->fetch();

$selectUser = $mysql->prepare("SELECT `email` FROM `users` WHERE `id` = :id");
$selectUser->execute(array(
	':id' => $server['owner_id']
));
    $user = $selectUser->fetch();

$selectNode = $mysql->prepare("SELECT * FROM `nodes` WHERE `id` = :id");
$selectNode->execute(array(
	':id' => $server['node']
));
    $node = $selectNode->fetch();

/*
 * Update Server SFTP Information
 */
$iv = $core->auth->generate_iv();
$pass = openssl_encrypt($_POST['sftp_pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));

$mysql->prepare("UPDATE `servers` SET `ftp_pass` = :pass, `encryption_iv` = :iv WHERE `id` = :sid")->execute(array(
    ':sid' => $core->server->getData('id'),
    ':pass' => $pass,
    ':iv' => $iv
));
	
/*
 * Connect to Node and Execute Password Update
 */
$core->ssh->generateSSH2Connection($core->server->nodeData('id'), true)->executeSSH2Command('cd /srv/scripts; sudo ./update_pass.sh "'.$server['ftp_user'].'" "'.$_POST['sftp_pass'].'"', false);

/*
 * Send the User an Email
 */
if(isset($_POST['email_user'])){
    
    $core->email->buildEmail('admin_new_sftppass', array(
        'PASS' => $_POST['sftp_pass'],
        'SERVER' => $server['name']
    ))->dispatch($user['email'], $core->settings->get('company_name').' - Your SFTP Password was Reset');
    
}

$core->page->redirect('../../settings.php');

?>
