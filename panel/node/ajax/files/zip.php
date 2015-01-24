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
use \ORM, \Unirest;

require_once '../../../../src/core/core.php';

if(!$core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash'))) {
	exit('Not authenticated.');
}

if($core->user->hasPermission('files.zip') !== true) {

	exit('You do not have permission to zip or unzip files.');

}

if(isset($_POST['zipItemPath']) && !empty($_POST['zipItemPath'])) {

	$request = Unirest\Request::put(
		"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/file/".$_POST['zipItemPath'],
		array(
			"X-Access-Token" => $core->server->getData('gsd_secret')
		),
		array(
			"zip" => $_POST['zipItemPath']
		)
	);

	if($request->code !== 200) {

		exit("An error was enocuntered trying to process your request. Server returned the following data\n\n[HTTP/1.1 {$request->code}] {$request->raw_body}");

	}

} elseif(isset($_POST['unzipItemPath']) && !empty($_POST['unzipItemPath'])) {

	$request = Unirest\Request::put(
		"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/file/".$_POST['unzipItemPath'],
		array(
			"X-Access-Token" => $core->server->getData('gsd_secret')
		),
		array(
			"unzip" => $_POST['unzipItemPath']
		)
	);

	if($request->code !== 200) {

		exit("An error was enocuntered trying to process your request. Server returned the following data\n\n[HTTP/1.1 {$request->code}] {$request->raw_body}");

	}

} else {

	exit("Error trying to handle an unknown function.");

}