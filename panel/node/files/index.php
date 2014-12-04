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
use \ORM;

require_once('../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Components\Page::redirect($core->settings->get('master_url').'index.php?login');
	exit();
}

if($core->user->hasPermission('files.view') !== true)
	Components\Page::redirect('../index.php?error=no_permission');

if(isset($_GET['file']))
	$_GET['file'] = str_replace('..', '', urldecode($_GET['file']));

if(isset($_GET['dir']))
	$_GET['dir'] = str_replace('..', '', urldecode($_GET['dir']));

if(isset($_GET['do']) && $_GET['do'] == 'download'){

	if($core->user->hasPermission('files.download') !== true)
		Components\Page::redirect('../index.php?error=no_permission');

	$url = "http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/file/".$_GET['file'];
	$context = stream_context_create(array(
		"http" => array(
			"method" => "GET",
			"header" => 'X-Access-Token: '.$core->server->getData('gsd_secret'),
			"timeout" => 3
		)
	));
	$content = json_decode(file_get_contents($url, 0, $context), true);

	/*
	* Download a File
	*/
	header("Pragma: public");
	header("Expires: 0");
	header("Cache-Control: must-revalidate, post-check=0, pre-check=0");
	header("Content-Type: application/octet-stream");
	header("Content-Description: File Transfer");
	header('Content-Disposition: attachment; filename="'.$_GET['file'].'"');
	header("Content-Length: ".mb_strlen($content['contents']));

	print( $content['contents'] );
	exit();

}

/*
 * Display Page
 */
echo $twig->render(
		'node/files/index.html', array(
			'server' => $core->server->getData(),
			'footer' => array(
				'seconds' => number_format((microtime(true) - $pageStartTime), 4)
			)
	));
?>
