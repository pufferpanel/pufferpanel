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
	$core->framework->page->redirect('../../../index.php');
}

if(isset($_GET['do']) && $_GET['do'] == 'ipuser') {

	if(!isset($_POST['nid']) || !is_numeric($_POST['nid']))
		$core->framework->page->redirect('../../list.php');
		
	if(!isset($_POST['warning']))
		$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=warning&disp=missing_warn&tab=sftp');
	
	if(!isset($_POST['sftp_ip'], $_POST['sftp_user']))
		$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=sftp_ip|sftp_user&disp=missing_args&tab=sftp');
		
	if(!filter_var($_POST['sftp_ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE))
		$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=sftp_ip&disp=ip_fail&tab=sftp');
			
	if(strlen($_POST['sftp_user']) < 1 || $_POST['sftp_user'] == 'root')
		$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=sftp_user&disp=user_fail&tab=sftp');
	
	/*
	 * Run Update on Node Table
	 */
	$mysql->prepare("UPDATE `nodes` SET `sftp_ip` = :ip, `username` = :name WHERE `id` = :nid")->execute(array(':ip' => $_POST['sftp_ip'], ':name' => $_POST['sftp_user'], ':nid' => $_POST['nid']));
	$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&tab=sftp');

}else if(isset($_GET['do']) && $_GET['do'] == 'pass'){

	if(!isset($_POST['nid']) || !is_numeric($_POST['nid']))
		$core->framework->page->redirect('../../list.php');
		
	if(!isset($_POST['warning']))
		$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=warning&disp=missing_warn&tab=sftp');
	
	if(!isset($_POST['pass']))
		$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=pass&disp=missing_args&tab=sftp');
		
	if(strlen($_POST['pass']) < 12)
		$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=pass&disp=pass_fail&tab=sftp');
	
	$iv = base64_encode(mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_CAST_256, MCRYPT_MODE_CBC), MCRYPT_RAND));
	$_POST['pass'] = openssl_encrypt($_POST['pass'], 'AES-256-CBC', file_get_contents(HASH), false, base64_decode($iv));
	
	/*
	 * Run Update on Node Table
	 */
	$mysql->prepare("UPDATE `nodes` SET `encryption_iv` = :iv, `password` = :pass WHERE `id` = :nid")->execute(array(':iv' => $iv, ':pass' => $_POST['pass'], ':nid' => $_POST['nid']));
	$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&tab=sftp');

}

?>