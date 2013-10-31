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