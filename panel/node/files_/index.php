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
use \ORM, \Unirest, \League\Flysystem\Filesystem as Filesystem, \League\Flysystem\Adapter\Ftp as Adapter;

require_once('../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Components\Page::redirect(Settings::config('master_url').'index.php?login');
	exit();
}

if($core->user->hasPermission('files.view') !== true)
	Components\Page::redirect('../index.php?error=no_permission');

if(isset($_GET['file']))
	$_GET['file'] = str_replace('..', '', urldecode($_GET['file']));

if(isset($_GET['dir']))
	$_GET['dir'] = str_replace('..', '', urldecode($_GET['dir']));

if(isset($_GET['do'], $_GET['file']) && $_GET['do'] == 'download' && !empty($_GET['file'])){

	if($core->user->hasPermission('files.download') !== true) {
		Components\Page::redirect('../index.php?error=no_permission');
	}

	if(!$core->gsd->avaliable($core->server->nodeData('ip'), $core->server->nodeData('gsd_listen'))) {
		Components\Page::redirect('index.php');
	}

	$_GET['file'] = str_replace("../", "", $_GET['file']);
	$downloadToken = $core->auth->keygen(32);

	$download = ORM::forTable('downloads')->create();
	$download->set(array(
		'server' => $core->server->getData('gsd_id'),
		'token' => $downloadToken,
		'path' => $_GET['file']
	));
	$download->save();

	Components\Page::redirect("http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id')."/download/".$downloadToken);

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
