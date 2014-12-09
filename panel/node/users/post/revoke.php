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

require_once('../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Components\Page::redirect($core->settings->get('master_url').'index.php?login');
	exit();

}

if($core->settings->get('allow_subusers') != 1)
	Components\Page::redirect('../list.php?error=not_enabled');

//id means pending, uid means not pending
if(isset($_GET['id']) && !empty($_GET['id'])){

	$_GET['id'] = urldecode($_GET['id']);

	$query = ORM::forTable('account_change')->where(array('key' => $_GET['id'], 'verified' => 0))->findOne();
	if(!$query) {
		Components\Page::redirect('../list.php?error');
	}

	// verify that this user is assigned to this server
	if(!array_key_exists($core->server->getData('hash'), json_decode($query->content, true)))
		Components\Page::redirect('../list.php?error=c1');

	// remove verification codes
	$query->delete();

	// update server in database
	$_perms = json_decode($core->server->getData('subusers'), true);
	unset($_perms[$_GET['id']]);
	$_perms = json_encode($_perms);

	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->subusers = $_perms;
	$server->save();

	Components\Page::redirect('../list.php?revoked');

}elseif(isset($_GET['uid']) && !empty($_GET['uid'])){

	$_GET['uid'] = urldecode($_GET['uid']);

	if(!$core->gsd->avaliable($core->server->nodeData('ip'), $core->server->nodeData('gsd_listen'))) {
		Components\Page::redirect('../view.php?id='.$_POST['uuid'].'&error');
	}

	$user = ORM::forTable('users')->where('uuid', $_GET['uid'])->findOne();

	if(!$user) {
		Components\Page::redirect('../list.php?error');
	}

	// verify that this user is assigned to this server
	if(!array_key_exists($core->server->getData('hash'), json_decode($user->permissions, true))) {
		Components\Page::redirect('../list.php?error');
	}

	// update server in database
	$_perms = json_decode($core->server->getData('subusers'), true);
	unset($_perms[$user->id]);
	$_perms = json_encode($_perms);

	$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
	$server->subusers = $_perms;
	$server->save();

	// update user
	$_uperms = json_decode($user->permissions, true);

	try {

		$request = Unirest::put(
			"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id'),
			array(
				"X-Access-Token" => $core->server->nodeData('gsd_secret')
			),
			array(
				"keys" => json_encode(array(
					$_uperms[$core->server->getData('hash')]['key'] => array()
				))
			)
		);

	} catch(\Exception $e) {

		\Tracy\Debugger::log($e);
		$exception = true;
		$outputMessage = '<div class="alert alert-danger">The server management daemon is not responding, we were unable to add your permissions. Please try again later.</div>';

	}

	unset($_uperms[$core->server->getData('hash')]);

	$user->permissions = json_encode($_uperms);
	$user->save();

	Components\Page::redirect('../list.php?revoked');

}
?>