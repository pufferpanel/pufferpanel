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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false) {
	Components\Page::redirect($core->settings->get('master_url').'index.php?login');
}

if($core->settings->get('allow_subusers') != 1) {
	Components\Page::redirect('../add.php?error=not_enabled');
}

if($core->auth->XSRF(@$_POST['xsrf']) !== true) {
	Components\Page::redirect('../add.php?error=token');
}

if(!isset($_POST['email'], $_POST['permissions'])) {
	Components\Page::redirect('../add.php?error=missing_required');
}

if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL)) {
	Components\Page::redirect('../add.php?error=email');
}

if(empty($_POST['permissions'])) {
	Components\Page::redirect('../add.php?error=permissions_empty');
}

/*
 * Does user exist already? If not we need to have them register first.
 */
$iv = $core->auth->generate_iv();
$registerKey = $core->auth->encrypt($_POST['email'], $iv).".".$iv;
$findUser = ORM::forTable('users')->where('email', $_POST['email'])->findOne();
if(!$findUser) {

	$account = ORM::forTable('account_change')->create();
	$account->set(array(
		'type' => 'user_register',
		'content' => 'null',
		'key' => $registerKey,
		'time' => time()
	));
	$account->save();

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

/*
 * Add subuser key. They will enter this on the account.php page so that we can handle any issues in a clener format.
 */
$subuserToken = $core->auth->keygen(32);
$account = ORM::forTable('account_change')->create();
$account->set(array(
	'type' => 'subuser',
	'content' => json_encode(array($core->server->getData('hash') => array('key' => $core->auth->generateUniqueUUID('servers', 'gsd_secret'), 'perms' => $_POST['permissions'], 'perms_gsd' => $gsdPermissions))),
	'key' => $subuserToken,
	'time' => time()
));
$account->save();

$subusers = (!is_null($core->server->getData('subusers')) && !empty($core->server->getData('subusers'))) ? json_decode($core->server->getData('subusers'), true) : array();
$subusers[$_POST['email']] = $subuserToken;

$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
$server->subusers = json_encode($subusers);
$server->save();

/*
* Send Email
*/
$core->email->buildEmail('new_subuser', array(
	'REGISTER_TOKEN' => $registerKey,
	'SUBUSER_TOKEN' => $subuserToken,
	'URLENCODE_TOKEN' => urlencode($registerKey),
	'SERVER' => $core->server->getData('name'),
	'EMAIL' => $_POST['email']
))->dispatch($_POST['email'], $core->settings->get('company_name').' - You\'ve Been Invited to Manage a Server');

Components\Page::redirect('../list.php?success');

?>
