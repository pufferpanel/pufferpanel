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

require_once('../../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	exit('Not authenticated!');
}

$server = ORM::forTable('servers')
			->select_many('servers.*', 'nodes.ip', 'nodes.gsd_listen')
			->select('nodes.gsd_secret', 'node_gsd_secret')
			->join('nodes', array('servers.node', '=', 'nodes.id'))
			->findOne();

if(!$server) {
	exit('That server does not exist.');
}

$newSecret = $core->auth->generateUniqueUUID('servers', 'gsd_secret');

try {

	$request = Unirest::put(
		'http://' . $server->ip . ':' . $server->gsd_listen . '/gameservers/' . $server->gsd_id,
		array(
			"X-Access-Token" => $server->node_gsd_secret
		),
		array(
			"keys" => json_encode(array(
				$server->gsd_secret => array(),
				$newSecret => array("s:ftp", "s:get", "s:power", "s:files", "s:files:get", "s:files:put", "s:query", "s:console", "s:console:send")
			))
		)
	);

	if($request->code != 200) {
		exit("Error trying to update token. GSD said: {$request->raw_body} [HTTP/1.1 $request->code ]");
	}

} catch(\Exception $e) {

	\Tracy\Debugger::log($e);
	exit('The server management daemon is not responding, we were unable to update the GSD token.');

}

$server->gsd_secret = $newSecret;
$server->save();

exit($server->gsd_secret);