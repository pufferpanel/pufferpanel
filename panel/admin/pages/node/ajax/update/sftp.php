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


if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php');
}

if(isset($_GET['do']) && $_GET['do'] == 'ipuser') {

	if(!isset($_POST['nid']) || !is_numeric($_POST['nid']))
		Page\components::redirect('../../list.php');
		
	if(!isset($_POST['warning']))
		Page\components::redirect('../../view.php?id='.$_POST['nid'].'&error=warning&disp=missing_warn&tab=sftp');
	
	if(!isset($_POST['sftp_ip'], $_POST['sftp_user']))
		Page\components::redirect('../../view.php?id='.$_POST['nid'].'&error=sftp_ip|sftp_user&disp=missing_args&tab=sftp');
		
	if(!filter_var($_POST['sftp_ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE))
		Page\components::redirect('../../view.php?id='.$_POST['nid'].'&error=sftp_ip&disp=ip_fail&tab=sftp');
			
	if(strlen($_POST['sftp_user']) < 1 || $_POST['sftp_user'] == 'root')
		Page\components::redirect('../../view.php?id='.$_POST['nid'].'&error=sftp_user&disp=user_fail&tab=sftp');
	
	/*
	 * Run Update on Node Table
	 */
	$mysql->prepare("UPDATE `nodes` SET `sftp_ip` = :ip, `username` = :name WHERE `id` = :nid")->execute(array(':ip' => $_POST['sftp_ip'], ':name' => $_POST['sftp_user'], ':nid' => $_POST['nid']));
	Page\components::redirect('../../view.php?id='.$_POST['nid'].'&tab=sftp');

}else if(isset($_GET['do']) && $_GET['do'] == 'pass'){

	if(!isset($_POST['nid']) || !is_numeric($_POST['nid']))
		Page\components::redirect('../../list.php');
		
	if(!isset($_POST['warning']))
		Page\components::redirect('../../view.php?id='.$_POST['nid'].'&error=warning&disp=missing_warn&tab=sftp');
	
	$_POST['ssh_pub_key'] = trim($_POST['ssh_pub_key']);
	$_POST['ssh_priv_key'] = trim($_POST['ssh_priv_key']);
	
	if(!preg_match('/^\/(.+)\/.ssh\/([^\/]+).pub$/', $_POST['ssh_pub_key']) || !preg_match('/^\/(.+)\/.ssh\/([^\/]+)$/', $_POST['ssh_priv_key']))
		Page\components::redirect('../../view.php?id='.$_POST['nid'].'&error=ssh_pub_key|ssh_priv_key&disp=key_fail&tab=sftp');
	
	/*
	 * Generate Encrypted Version of Secret
	 */
	$ssh_secret_iv = (!empty($_POST['ssh_secret'])) ? $core->auth->generate_iv() : null;
	$ssh_secret = (!empty($_POST['ssh_secret'])) ? $core->auth->encrypt($_POST['ssh_secret'], $ssh_secret_iv) : null;
	
	/*
	 * Run Update on Node Table
	 */
	$mysql->prepare("UPDATE `nodes` SET `ssh_pub` = :ssh_pub, `ssh_priv` = :ssh_priv, `ssh_secret` = :ssh_secret, `ssh_secret_iv` = :ssh_secret_iv WHERE `id` = :nid")->execute(array(
		':ssh_pub' => $_POST['ssh_pub_key'],
		':ssh_priv' => $_POST['ssh_priv_key'],
		':ssh_secret' => $ssh_secret,
		':ssh_secret_iv' => $ssh_secret_iv,
		':nid' => $_POST['nid']
	));
	Page\components::redirect('../../view.php?id='.$_POST['nid'].'&tab=sftp');

}

?>