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
use \ORM, \Unirest, \Tracy;

require_once '../../../../src/core/core.php';

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === true) {

	if(!$core->user->hasPermission('console.power')) {
		exit('You do not have the required permissions to perform this function.');
	}

	/*
	 * Open Stream for Reading/Writing
	 */
	$rewrite = false;
	$errorMessage = "Unable to process your request. Please try again.";

	try {

		$response = Unirest::get(
			"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/file/server.properties",
			array(
				"X-Access-Token" => $core->server->getData('gsd_secret')
			)
		);

	} catch(\Exception $e) {

		Tracy\Debugger::log($e);
		exit($errorMessage." Unable to connect to remote host.");

	}

	/*
	 * Typically Means Server is Off
	 */
	if(!in_array($response->code, array(200, 500))) {
		switch($response->code) {

			case 403:
				exit($errorMessage." Authentication error encountered.");
				break;
			default:
				exit("$errorMessage HTTP/$response->code. Invalid response was recieved. ($response->raw_body)");
				break;

		}
	}

	if($response->code == 500 || !isset($response->body->contents) || empty($response->body->contents)) {

		/*
		 * Create server.properties
		 */
		if(!file_exists(APP_DIR.'templates/server.properties.tpl') || empty(file_get_contents(APP_DIR.'templates/server.properties.tpl')))
			exit($errorMessage." No Template Avaliable for server.properties");

		$put = Unirest::put(
			"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/file/server.properties",
			array(
				"X-Access-Token" => $core->server->getData('gsd_secret')
			),
			array(
				"contents" => sprintf(file_get_contents(APP_DIR.'templates/server.properties.tpl'), $core->server->getData('server_port'), $core->server->getData('server_ip'))
			)
		);

        if(!empty($put->body)) {
        	exit($errorMessage." Unable to make server.properties");
		}

		$core->log->getUrl()->addLog(0, 1, array('system.create_serverprops', 'A new server.properties file was created for your server.'));

	} else {

		$lines = explode("\n", $response->body->contents);
		$newContents = $response->body->contents;

		foreach($lines as $line) {

			$var = explode('=', $line);

			if($var[0] == 'server-port' && $var[1] != $core->server->getData('server_port')) {

				$newContents = str_replace('server-port='.$var[1], "server-port=".$core->server->getData('server_port')."\n", $newContents);
				$rewrite = true;

			} else if($var[0] == 'online-mode' && $var[1] == 'false' && $core->settings->get('force_online') == 1) {

				$newContents = str_replace('online-mode='.$var[1], "online-mode=true\n", $newContents);
				$rewrite = true;

			} else if($var[0] == 'enable-query' && $var[1] != 'true') {

				$newContents = str_replace('enable-query='.$var[1], "enable-query=true\n", $newContents);
				$rewrite = true;

			} else if($var[0] == 'query.port' && $var[1] != $core->server->getData('server_port')) {

				$newContents = str_replace('query.port='.$var[1], "query.port=".$core->server->getData('server_port')."\n", $newContents);
				$rewrite = true;

			} else if($var[0] == 'server-ip' && $var[1] != $core->server->getData('server_ip')) {

				$newContents = str_replace('server-ip='.$var[1], "server-ip=".$core->server->getData('server_ip')."\n", $newContents);
				$rewrite = true;

			}

        }

	}

	/*
	 * Write New Data
	 */
	if($rewrite) {

		$put = Unirest::put(
			"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/file/server.properties",
			array(
				"X-Access-Token" => $core->server->getData('gsd_secret')
			),
			array(
				"contents" => $newContents
			)
		);

		if(!empty($put->body)) {
	    	exit($errorMessage." Unable to update server.properties");
		}

        $core->log->getUrl()->addLog(0, 0, array('system.serverprops_updated', 'The server properties file was updated to match the assigned information.'));

	}

    /*
	 * Connect and Run Function
	 */
	$get = Unirest::get(
		"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/on",
		array(
			"X-Access-Token" => $core->server->getData('gsd_secret')
		)
	);

	if($get->body != "ok")
		exit($errorMessage." Unable to start server (".$get->raw_body.")");

	echo 'ok';

} else {

	die('Invalid Authentication.');

}
?>
