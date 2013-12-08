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
require_once('../../../../../core/framework/framework.core.php');

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}

if(!isset($_POST['sid']) || !isset($_POST['nid']))
	$core->framework->page->redirect('../../find.php');
	
if(!isset($_POST['sftp_pass'], $_POST['sftp_pass_2'], $_POST['nid']))
	$core->framework->page->redirect('../../view.php?id='.$_POST['sid']);
	
if(strlen($_POST['sftp_pass']) < 8)
	$core->framework->page->redirect('../../view.php?id='.$_POST['sid'].'&error=sftp_pass|sftp_pass_2&disp=pass_len');
	
if($_POST['sftp_pass'] != $_POST['sftp_pass_2'])
	$core->framework->page->redirect('../../view.php?id='.$_POST['sid'].'&error=sftp_pass|sftp_pass_2&disp=pass_match');

/* 
 * Select Node, User, & Server Information
 */
$select = $mysql->prepare("SELECT `ftp_user`, `name`, `owner_id` FROM `servers` WHERE `id` = ?");
$select->execute(array($_POST['sid']));
    $server = $select->fetch();

$selectUser = $mysql->prepare("SELECT `email` FROM `users` WHERE `id` = ?");
$selectUser->execute(array($server['owner_id']));
    $user = $selectUser->fetch();

$selectNode = $mysql->prepare("SELECT `username`, `sftp_ip`, `password`, `encryption_iv` FROM `nodes` WHERE `id` = ?");
$selectNode->execute(array($_POST['nid']));
    $node = $selectNode->fetch();

/*
 * Update Server SFTP Information
 */
$iv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
$pass = openssl_encrypt($_POST['sftp_pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));

$mysql->prepare("UPDATE `servers` SET `ftp_pass` = :pass, `encryption_iv` = :iv WHERE `id` = :sid")->execute(array(
    ':sid' => $_POST['sid'],
    ':pass' => $pass,
    ':iv' => $iv
));
	
/*
 * Connect to Node and Execute Password Update
 */
$con = ssh2_connect($node['sftp_ip'], 22);
ssh2_auth_password($con, $node['username'], openssl_decrypt($node['password'], 'AES-256-CBC', file_get_contents(HASH), 0, base64_decode($node['encryption_iv'])));

$stream = ssh2_exec($con, 'cd /srv/scripts; ./update_pass.sh "'.$server['ftp_user'].'" "'.$_POST['sftp_pass'].'"');

/*
 * Send the User an Email
 */
if(isset($_POST['email_user'])){
    
    $core->framework->email->buildEmail('admin_new_sftppass', array(
        'PASS' => $_POST['sftp_pass'],
        'SERVER' => $server['name']
    ))->dispatch($user['email'], $core->framework->settings->get('company_name').' - Your SFTP Password was Reset');
    
}

$core->framework->page->redirect('../../view.php?id='.$_POST['sid']);

?>