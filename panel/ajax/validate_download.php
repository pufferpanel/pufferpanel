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
require_once '../../src/core/core.php';

/*
* Illegally Accessed File
*/
if(!isset($_POST['token'], $_POST['server'])) {

	http_response_code(403);
	exit();

}

/*
* Verify Identity
*/
$download = ORM::forTable('downloads')
	->where(array(
		'server' => $_POST['server'],
		'token' => $_POST['token']
	))->findOne();

if(!$download) {

	http_response_code(404);
	exit();

} else {

	$download->delete();

	http_response_code(200);
	exit(json_encode(array(
		'path' => $download->path
	)));

}