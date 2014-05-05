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
require_once('../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === true){
	
	if(!isset($_POST['sftp_pass'], $_POST['sftp_pass_2']))
		Page\components::redirect('../../settings.php');
		
	if(strlen($_POST['sftp_pass']) < 8)
		Page\components::redirect('../../settings.php?error=sftp_pass|sftp_pass_2&disp=pass_len');
		
	if($_POST['sftp_pass'] != $_POST['sftp_pass_2'])
		Page\components::redirect('../../settings.php?error=sftp_pass|sftp_pass_2&disp=pass_match');
	
	/*
	 * Update Server SFTP Information
	 */
	$iv = $core->auth->generate_iv();
	$pass = $core->auth->encrypt($_POST['sftp_pass'], $iv);
	
	$mysql->prepare("UPDATE `servers` SET `ftp_pass` = :pass, `encryption_iv` = :iv WHERE `id` = :sid")->execute(array(
	    ':sid' => $core->server->getData('id'),
	    ':pass' => $pass,
	    ':iv' => $iv
	));
		
	/*
	 * Connect to Node and Execute Password Update
	 */
	$core->ssh->generateSSH2Connection($core->server->nodeData('id'), true)->executeSSH2Command('cd /srv/scripts; sudo ./update_pass.sh "'.$core->server->getData('ftp_user').'" "'.$_POST['sftp_pass'].'"', false);
	
	Page\components::redirect('../../settings.php?success');
	
}

?>