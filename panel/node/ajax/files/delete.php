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
use \ORM, \League\Flysystem\Filesystem, \League\Flysystem\Adapter\Ftp as Adapter;

require_once '../../../../src/core/core.php';

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false) {
	exit('Not authenticated.');
}

if($core->user->hasPermission('files.delete') !== true) {
	exit('You do not have permission to delete files.');
}

if(!isset($_POST['deleteItemType'], $_POST['deleteItemPath'])) {
	exit('Not enough variables were passed.');
}

/*
 * Delete File
 */
if($_POST['deleteItemType'] == 'file' && !empty($_POST['deleteItemPath'])) {

	try {

		$cid = ftp_ssl_connect($core->server->nodeData('ip'));
		$login = ftp_login($cid, $core->server->getData('ftp_user').'-'.$core->server->getData('gsd_id'), $core->auth->decrypt($core->server->getData('ftp_pass'), $core->server->getData('encryption_iv')));

		$_POST['deleteItemPath'] = urldecode($_POST['deleteItemPath']);

		ftp_delete($cid, $_POST['deleteItemPath']);
		ftp_close($cid);
		echo 'ok';

	} catch(\FtpException $e) {
		exit('Error occured trying to delete that file! '.$e->getMessage());
	}

	// try {
	//
	// 	$_POST['deleteItemPath'] = urldecode($_POST['deleteItemPath']);
	// 	$filesystem->delete('1417640539.txt');
	// 	echo 'ok';
	//
	// } catch(\League\Flysystem\FileNotFoundException $e){
	// 	exit('Error occured trying to delete that file! '.$e->getMessage());
	// }

} else if($_POST['deleteItemType'] == 'dir' && !empty($_POST['deleteItemPath'])) {

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
		\Tracy\Debugger::log($e);
		exit('An execption occured.');
	}

	try {

		$_POST['deleteItemPath'] = urldecode($_POST['deleteItemPath']);
		$filesystem->deleteDir($_POST['deleteItemPath']);
		echo 'ok';

	} catch(\League\Flysystem\FileNotFoundException $e){
		exit('Error trying to delete that directory. '.$e->getMessage());
	}

} else {
	var_dump($_POST);
	echo 'Nothing was matched in the script.';
}