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
header('Content-Type: application/json');
require_once('../../src/framework/framework.core.php');

/*
 * Illegally Accessed File
 */
if(!isset($_POST['username']) || !isset($_POST['password'])){

	http_response_code(403);

}

if(!preg_match('^([mc-]{3})([\w\d\-]{12})[\-]([\d]+)$^', $_POST['username'], $matches)){

	http_response_code(403);

}else{

	$username = $matches[1].$matches[2];
	$serverid = $matches[3];

}

/*
 * Varify Identity
 */
$query = $mysql->prepare("SELECT `encryption_iv`, `ftp_pass`, `gsd_secret` FROM `servers` WHERE `gsd_id` = :gsdid AND `ftp_user` = :username");
$query->execute(array(
	'gsdid' => $serverid,
	'username' => $username
));

	if($query->rowCount() != 1){

		http_response_code(403);
		exit();

	}else{

		$server = $query->fetch();
		if($core->auth->encrypt($_POST['password'], $server['encryption_iv']) != $server['ftp_pass']){

			http_response_code(403);
			exit();

		}else{

			http_response_code(200);
			die(json_encode(array('authkey' => $server['gsd_secret'])));

		}

	}

?>