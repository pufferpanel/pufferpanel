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
use \ORM, \Flow, \Tracy, \League\Flysystem\Filesystem as Filesystem, \League\Flysystem\Adapter\Ftp as Adapter;

require_once '../../../../src/core/core.php';

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false) {
	http_response_code(403);
	return;
}

if($core->user->hasPermission('files.delete') !== true) {
	http_response_code(403);
	return;
}

if(!isset($_POST['newFilePath'])) {
	http_response_code(404);
	return;
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
	exit();

}

$tempDir = '/tmp/'.$core->server->getData('hash');
$uploadPath = SRC_DIR . 'uploads/' . $core->server->getData('hash') . '/';

if(!is_dir($tempDir)) {
	mkdir($tempDir, 0777);
}

if(!is_dir($uploadPath)) {
	mkdir($uploadPath, 0777);
}

$config = new Flow\Config();
$config->setTempDir($tempDir);
$file = new Flow\File($config);

if($_SERVER['REQUEST_METHOD'] === 'GET') {

	if($file->checkChunk()) {
		http_response_code(200);
	} else {

		http_response_code(404);
		return;

	}

} else {

	if($file->validateChunk()) {
		$file->saveChunk();
	} else {

		http_response_code(400);
		return;

	}

}

if($file->validateFile() && $file->save($uploadPath.$_POST['flowFilename'])) {

	$stream = file_get_contents($uploadPath.$_POST['flowFilename']);
	if(!$filesystem->write(rtrim($_POST['newFilePath'], '/').'/'.$_POST['flowFilename'], $stream)) {
		http_response_code(500);
	} else {
		http_response_code(200);
	}

}else{

}