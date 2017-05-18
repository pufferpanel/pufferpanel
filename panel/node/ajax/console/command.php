<?php //CH
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
require_once('../../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === true){

    /*
	 * Connect and Run Function
	 */
	$context = stream_context_create(array(
		"http" => array(
			"method" => "POST",
			"header" => 'X-Access-Token: '.$core->server->nodeData('gsd_secret'),
			"timeout" => 3,
			"content" => http_build_query(array('command' => $_POST['command'])),
		)
	));
	$gatherData = @file_get_contents("http://".$core->server->nodeData('sftp_ip').":8003/gameservers/".$core->server->getData('gsd_id')."/console", 0, $context);

	if($gatherData != "\"ok\""){
		exit("An error was encountered with this AJAX request. ($gatherData)");
  }
	echo 'ok';

}else{

	die('Invalid Authentication.');

}
?>
