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

	if($core->user->hasPermission('console.power') !== true)
		exit('You do not have the required permissions to perform this function.');

	/*
	 * Open Stream for Reading/Writing
	 */
	$rewrite = false;
	$errorMessage = "Unable to process your request. Please try again.";

	$url = "http://".$core->server->nodeData('ip').":8003/gameservers/".$core->server->getData('gsd_id')."/file/server.properties";

	$curl = curl_init($url);
	curl_setopt($curl, CURLOPT_CUSTOMREQUEST, "GET");
	curl_setopt($curl, CURLOPT_HTTPHEADER, array(
	    'X-Access-Token: '.$core->server->getData('gsd_secret')
	));
	curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
	$response = curl_exec($curl);

	/*
	 * Typically Means Server is Off
	 */
	if(empty($response))
		exit($errorMessage."<!--Empty Response-->");

	$json = json_decode($response, true);

	/*
	 * Usually Occurs when there is an authentication error.
	 */
	if(json_last_error() != "JSON_ERROR_NONE")
		exit("An error was encountered with this AJAX request. (Invalid JSON Response - ".json_last_error_msg().")");

	if(!array_key_exists('contents', $json)) {

		/*
		 * Create server.properties
		 */
		if(!file_exists(APP_DIR.'templates/server.properties.tpl') || empty(file_get_contents(APP_DIR.'templates/server.properties.tpl')))
			exit($errorMessage."<!--No Template Avaliable for server.properties-->");

		$data = array("contents" => sprintf(file_get_contents(APP_DIR.'templates/server.properties.tpl'), $core->server->getData('server_port'), $core->server->getData('server_ip')));
		$curl = curl_init($url);
		curl_setopt($curl, CURLOPT_CUSTOMREQUEST, "PUT");
		curl_setopt($curl, CURLOPT_HTTPHEADER, array(
		    'X-Access-Token: '.$core->server->getData('gsd_secret')
		));
		curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($data));
		$response = curl_exec($curl);

	        if(!empty($response))
	        	exit($errorMessage."<!--Unable to make server.properties-->");

		$core->log->getUrl()->addLog(0, 1, array('system.create_serverprops', 'A new server.properties file was created for your server.'));

	}else{

		$lines = explode("\n", $json['contents']);
		$newContents = $json['contents'];
		foreach($lines as $line){

			$var = explode('=', $line);

				if($var[0] == 'server-port' && $var[1] != $core->server->getData('server_port')){
					$newContents = str_replace('server-port='.$var[1], "server-port=".$core->server->getData('server_port')."\n", $newContents);
					$rewrite = true;
				}else if($var[0] == 'online-mode' && $var[1] == 'false'){
					if($core->settings->get('force_online') == 1){
						$newContents = str_replace('online-mode='.$var[1], "online-mode=true\n", $newContents);
						$rewrite = true;
					}
				}else if($var[0] == 'enable-query' && $var[1] != 'true'){
					$newContents = str_replace('enable-query='.$var[1], "enable-query=true\n", $newContents);
					$rewrite = true;
				}else if($var[0] == 'query.port' && $var[1] != $core->server->getData('server_port')){
					$newContents = str_replace('query.port='.$var[1], "query.port=".$core->server->getData('server_port')."\n", $newContents);
					$rewrite = true;
				}else if($var[0] == 'server-ip' && $var[1] != $core->server->getData('server_ip')){
					$newContents = str_replace('server-ip='.$var[1], "server-ip=".$core->server->getData('server_ip')."\n", $newContents);
					$rewrite = true;
				}

        }

	}

	/*
	 * Write New Data
	 */
	if($rewrite === true){

		$data = array("contents" => $newContents);
		$curl = curl_init($url);
		curl_setopt($curl, CURLOPT_CUSTOMREQUEST, "PUT");
		curl_setopt($curl, CURLOPT_HTTPHEADER, array(
		    'X-Access-Token: '.$core->server->getData('gsd_secret')
		));
		curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
		curl_setopt($curl, CURLOPT_POSTFIELDS, http_build_query($data));
		$response = curl_exec($curl);

		    if(!empty($response))
		    	exit($errorMessage."<!--Unable to update server.properties-->");

        $core->log->getUrl()->addLog(0, 0, array('system.serverprops_updated', 'The server properties file was updated to match the assigned information.'));

	}

    /*
	 * Connect and Run Function
	 */
	$context = stream_context_create(array(
		"http" => array(
			"method" => "GET",
			"header" => 'X-Access-Token: '.$core->server->getData('gsd_secret'),
			"timeout" => 3,
			"ignore_errors" => true
		)
	));
	$gatherData = @file_get_contents("http://".$core->server->nodeData('ip').":8003/gameservers/".$core->server->getData('gsd_id')."/on", 0, $context);

	if($gatherData != "\"ok\"")
		exit($errorMessage."<!--Unable to start server (".$gatherData.")-->");

	echo 'ok';

}else{

	die('Invalid Authentication.');

}
?>
