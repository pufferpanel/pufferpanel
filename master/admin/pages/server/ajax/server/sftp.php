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

if(!isset($_POST['sid']))
	$core->framework->page->redirect('../../find.php');
	
if(!isset($_POST['sftp_pass'], $_POST['sftp_pass_2'], $_POST['nid']))
	$core->framework->page->redirect('../../view.php?id='.$_POST['sid']);
	
if(strlen($_POST['sftp_pass']) < 8))
	$core->framework->page->redirect('../../view.php?id='.$_POST['sid'].'&error=sftp_pass|sftp_pass_2&disp=pass_len');
	
if($_POST['sftp_pass'] != $_POST['sftp_pass_2'])
	$core->framework->page->redirect('../../view.php?id='.$_POST['sid'].'&error=sftp_pass|sftp_pass_2&disp=pass_match');

// Chupacabraalabra. nap.
$node = exit();
	
/*
 * Connect to Node and Execute Password Update
 */