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
	$core->framework->page->redirect('../../../../../index.php');
}

if(!isset($_POST['uid']) || !is_numeric($_POST['uid']))
	$core->framework->page->redirect('../find.php?error=UPDATE-USER__undefined_user');
	
if($_POST['pass'] != $_POST['pass_2'])
	$core->framework->page->redirect('../../view.php?id='.$_POST['uid'].'&error=password');
	
$update = $mysql->prepare("UPDATE `users` SET `password` = :password WHERE `id` = :uid");
$update->execute(array(
	':password' => $core->framework->auth->encrypt($_POST['pass']),
	':uid' => $_POST['uid']
));

	if(isset($_POST['email_user'])){
	
		/*
		 * Send Email
		 */
		$message = $core->framework->email->generateNewPasswordEmail(array('NEW_PASS' => $_POST['pass'], 'EMAIL' => $_POST['email']));
		
		$core->framework->email->dispatch($_POST['email'], $core->framework->settings->get('company_name').' - An Admin has Reset Your Password', $message);
	
	}
	
	if(isset($_POST['clear_session'])){
	
		$update = $mysql->prepare("UPDATE `users` SET `session_id` = '', `session_ip` = '', `session_expires` = '0' WHERE `id` = :uid");
		$update->execute(array(
			':uid' => $_POST['uid']
		));
	
	}

$core->framework->page->redirect('../../view.php?id='.$_POST['uid'].'&disp=p_updated');

?>