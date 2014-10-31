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

require_once('../../../../../src/core/core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true)
	Components\Page::redirect('../../../../index.php?login');

if(!preg_match('/^[\w-]{4,35}$/', $_POST['username']))
	Components\Page::redirect('../../new.php?disp=u_fail');

if(!filter_var($_POST['email'], FILTER_VALIDATE_EMAIL))
	Components\Page::redirect('../../new.php?disp=e_fail');

if(strlen($_POST['pass']) < 8 || $_POST['pass'] != $_POST['pass_2'])
	Components\Page::redirect('../../new.php?disp=p_fail');

$query = ORM::forTable('users')->where_any_is(array(array('username' => $_POST['username']), array('email' => $_POST['email'])))->findOne();

if($query === false)
	Components\Page::redirect('../../new.php?disp=a_fail');

$user = ORM::forTable('users')->create();
$user->set(array(
	'uuid' => $core->auth->gen_UUID(),
	'username' => $_POST['username'],
	'email' => $_POST['email'],
	'password' => $core->auth->hash($_POST['pass']),
	'language' => $core->settings->get('default_language'),
	'register_time' => time()
));
$user->save();

/*
 * Send Email
 */
$core->email->buildEmail('admin_newaccount', array(
    'PASS' => $_POST['pass'],
    'EMAIL' => $_POST['email']
))->dispatch($_POST['email'], $core->settings->get('company_name').' - Account Created');

Components\Page::redirect('../../view.php?id='.$mysql->lastInsertId());

?>