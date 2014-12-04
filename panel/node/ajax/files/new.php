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
use \ORM, \Tracy, \League\Flysystem\Filesystem as Filesystem, \League\Flysystem\Adapter\Ftp as Adapter;

require_once '../../../../src/core/core.php';

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false) {
	exit('<div class="alert alert-danger">Not authenticated.</div>');
}

if($core->user->hasPermission('files.delete') !== true) {
	exit('<div class="alert alert-danger">You do not have permission to delete files.</div>');
}

if(!isset($_POST['newFilePath'], $_POST['newFileContents'])) {
	exit('<div class="alert alert-danger">Missing some required parameters.</div>');
}

try {

	$filesystem = new Filesystem(new Adapter(array(
		'host' => $core->server->nodeData('ip'),
		'username' => $core->server->getData('ftp_user').'-'.$core->server->getData('gsd_id'),
		'password' => $core->auth->decrypt($core->server->getData('ftp_pass'), $core->server->getData('encryption_iv')),
		'port' => 21,
		'passive' => true,
		'ssl' => true,
		'timeout' => 10
	)));

} catch(\Exception $e) {

	http_response_code(500);
	Tracy\Debugger::log($e);
	exit('<div class="alert alert-danger">An execption occured when trying to connect to the server.</div>');

}

try {

	if(!$filesystem->write(urldecode($_POST['newFilePath']), $_POST['newFileContents'])) {

		http_response_code(500);
		exit('<div class="alert alert-danger">An error occured when trying to write this file to the system.</div>');

	} else {
		exit('ok');
	}

} catch(\Exception $e) {

	http_response_code(500);
	Tracy\Debugger::log($e);
	exit('<div class="alert alert-danger">An execption occured when trying to write the file to the server.</div>');

}