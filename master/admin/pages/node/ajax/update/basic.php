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

if(!isset($_POST['nid']) || !is_numeric($_POST['nid']))
	$core->framework->page->redirect('../../list.php');

if(!isset($_POST['name'], $_POST['link'], $_POST['ip']))
	$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&disp=missing_args');

/*
 * Validate Node Name
 */
if(!preg_match('/^[\w.-]{1,15}$/', $_POST['name']))
	$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=name&disp=n_fail');
		
/*
 * Validate node_ip & node_sftp_ip
 */
if(!filter_var($_POST['ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE))
	$core->framework->page->redirect('../../view.php?id='.$_POST['nid'].'&error=ip&disp=ip_fail');
	
/*
 * Update Record
 */
$mysql->prepare("UPDATE `nodes` SET `node` = :name, `node_ip` = :ip WHERE `id` = :nid")->execute(array(
	':nid' => $_POST['nid'],
	':name' => $_POST['name'],
	':link' => $_POST['link'],
	':ip' => $_POST['ip']
));

$core->framework->page->redirect('../../view.php?id='.$_POST['nid']);

?>