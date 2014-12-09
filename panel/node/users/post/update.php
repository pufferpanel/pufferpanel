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

if(!isset($_POST['uuid'], $_POST['permissions']))
	Components\Page::redirect('../list.php');

if($core->auth->XSRF(@$_POST['xsrf']) !== true)
	Components\Page::redirect('../list.php?id='.$_POST['uuid'].'&error');

if(empty($_POST['permissions']))
	Components\Page::redirect('../view.php?id='.$_POST['uuid'].'&error');

if(!$core->gsd->avaliable($core->server->nodeData('ip'), $core->server->nodeData('gsd_listen'))) {
	Components\Page::redirect('../view.php?id='.$_POST['uuid'].'&error');
}

$query = ORM::forTable('users')->where('uuid', $_POST['uuid'])->findOne();

if(!$query) {
	Components\Page::redirect('../list.php?error');
}

foreach($_POST['permissions'] as $id => $permission) {

	if(in_array($permission, array('files.edit', 'files.save', 'files.download', 'files.delete', 'files.create', 'files.upload', 'files.zip')) && !in_array('files.view', $_POST['permissions'])) {
		$_POST['permissions'] = array_merge($_POST['permissions'], array("files.view"));
	}

	if(in_array($permission, array('manage.rename.jar')) && !in_array('manage.rename.view', $_POST['permissions'])) {
		$_POST['permissions'] = array_merge($_POST['permissions'], array("manage.rename.view"));
	}

	if(in_array($permission, array('manage.ftp.details', 'manage.ftp.password')) && !in_array('manage.ftp.view', $_POST['permissions'])) {
		$_POST['permissions'] = array_merge($_POST['permissions'], array("manage.ftp.view"));
	}

}

$gsdPermissions = array("s:console", "s:query");
foreach($_POST['permissions'] as $id => $permission) {

	switch($permission) {

		case "console.power":
			$gsdPermissions = array_merge($gsdPermissions, array("s:power"));
			break;
		case "console.command":
			$gsdPermissions = array_merge($gsdPermissions, array("s:console:command"));
			break;
		case "files.view":
			$gsdPermissions = array_merge($gsdPermissions, array("s:files"));
			break;
		case "files.edit":
			$gsdPermissions = array_merge($gsdPermissions, array("s:files:get"));
			break;
		case "files.save":
			$gsdPermissions = array_merge($gsdPermissions, array("s:files:put"));
			break;
		case "files.zip":
			$gsdPermissions = array_merge($gsdPermissions, array("s:files:zip"));
			break;

	}

}

$permissions = @json_decode($query->permissions, true);
if(!is_array($permissions) || !array_key_exists($core->server->getData('hash'), $permissions))
	Components\Page::redirect('../view.php?id='.$_POST['uuid'].'&error');

$permissions[$core->server->getData('hash')]['perms'] = $_POST['permissions'];
$permissions[$core->server->getData('hash')]['perms_gsd'] = $gsdPermissions;
$query->permissions = json_encode($permissions);
$query->save();

try {

	$request = Unirest::put(
		"http://".$core->server->nodeData('ip').":".$core->server->nodeData('gsd_listen')."/gameservers/".$core->server->getData('gsd_id'),
		array(
			"X-Access-Token" => $core->server->nodeData('gsd_secret')
		),
		array(
			"keys" => json_encode(array(
				$permissions[$core->server->getData('hash')]['key'] => $gsdPermissions
			))
		)
	);

} catch(\Exception $e) {

	\Tracy\Debugger::log($e);
	$exception = true;
	$outputMessage = '<div class="alert alert-danger">The server management daemon is not responding, we were unable to add your permissions. Please try again later.</div>';

}

Components\Page::redirect('../view.php?id='.$_POST['uuid']);

?>