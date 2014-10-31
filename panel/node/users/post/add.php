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
use \ORM as ORM;

require_once('../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), $core->auth->getCookie('pp_server_hash')) === false){

	Components\Page::redirect($core->settings->get('master_url').'index.php?login');
	exit();

}

if($core->auth->XSRF(@$_POST['xsrf']) !== true)
	Components\Page::redirect('../add.php?error=token');

if(!isset($_POST['email'], $_POST['permissions']))
	Components\Page::redirect('../add.php?error=missing_required');

if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL))
	Components\Page::redirect('../add.php?error=email');

if(empty($_POST['permissions']))
	Components\Page::redirect('../add.php?error=permissions_empty');

$iv = $core->auth->generate_iv();

$account = ORM::forTable('account_change')->create();
$account->set(array(
	'type' => 'subuser',
	'content' => array($core->server->getData('hash') => $_POST['permissions']),
	'key' => $core->auth->encrypt($_POST['email'], $iv).".".$iv,
	'time' => time()
));
$account->save();

$subusers = (!is_null($core->server->getData('subusers')) && !empty($core->server->getData('subusers'))) ? json_decode($core->server->getData('subusers'), true) : array();
$subusers[$_POST['email']] = $iv;

$server = ORM::forTable('servers')->findOne($core->server->getData('id'));
$server->subusers = json_encode($subusers);
$server->save();

/*
* Send Email
*/
$core->email->buildEmail('new_subuser', array(
	'TOKEN' => $core->auth->encrypt($_POST['email'], $iv).".".$iv,
	'URLENCODE_TOKEN' => urlencode($core->auth->encrypt($_POST['email'], $iv).".".$iv),
	'SERVER' => $core->server->getData('name'),
	'EMAIL' => $_POST['email']
))->dispatch($_POST['email'], $core->settings->get('company_name').' - You\'ve Been Invited to Manage a Server');

Components\Page::redirect('../list.php?success');

?>