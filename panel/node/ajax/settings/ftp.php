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

$klein->respond('*', function($request, $response) use ($core) {
	if(!$core->user->hasPermission('manage.ftp.password')) {
		$response->redirect('../index?error=no_permission', 302)->send();
	}

	if(!isset($_POST['ftp_pass'], $_POST['ftp_pass_2'])) {
		$response->redirect('../settings?error=ftp_pass|ftp_pass_2&disp=no_pass', 302)->send();
	}

	if(strlen($_POST['ftp_pass']) < 8) {
		$response->redirect('../settings?error=ftp_pass|ftp_pass_2&disp=pass_len', 302)->send();
	}

	if($_POST['ftp_pass'] != $_POST['ftp_pass_2']) {
		$response->redirect('../settings?error=ftp_pass|ftp_pass_2&disp=pass_match', 302)->send();
	}

	/*
	 * Update Server ftp Information
	 */
	$iv = $core->auth->generate_iv();

	$ftp = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$ftp->ftp_pass = $core->auth->encrypt($_POST['ftp_pass'], $iv);
	$ftp->encryption_iv = $iv;
	$ftp->save();

	$response->redirect('../settings?success', 302)->send();
});
