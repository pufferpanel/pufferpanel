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
namespace PufferPanel\Core;

require_once('../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === true){

	if($core->user->hasPermission('manage.ftp.password') !== true)
		Components\Page::redirect('../../index.php?error=no_permission');

	if(!isset($_POST['ftp_pass'], $_POST['ftp_pass_2']))
		Components\Page::redirect('../../settings.php');

	if(strlen($_POST['ftp_pass']) < 8)
		Components\Page::redirect('../../settings.php?error=ftp_pass|ftp_pass_2&disp=pass_len');

	if($_POST['ftp_pass'] != $_POST['ftp_pass_2'])
		Components\Page::redirect('../../settings.php?error=ftp_pass|ftp_pass_2&disp=pass_match');

	/*
	 * Update Server ftp Information
	 */
	$iv = $core->auth->generate_iv();

	$mysql->prepare("UPDATE `servers` SET `ftp_pass` = :pass, `encryption_iv` = :iv WHERE `id` = :sid")->execute(array(
	    ':sid' => $core->server->getData('id'),
	    ':pass' => $core->auth->encrypt($_POST['ftp_pass'], $iv),
	    ':iv' => $iv
	));

	Components\Page::redirect('../../settings.php?success');

}

?>
