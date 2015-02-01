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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Components\Page::redirect('../../../index.php');
}

if(!isset($_POST['sid']) || !isset($_POST['nid']))
	Components\Page::redirect('../../find.php');

if(!isset($_POST['ftp_pass'], $_POST['ftp_pass_2'], $_POST['nid']))
	Components\Page::redirect('../../view.php?id='.$_POST['sid']);

if(strlen($_POST['ftp_pass']) < 8)
	Components\Page::redirect('../../view.php?id='.$_POST['sid'].'&error=ftp_pass|ftp_pass_2&disp=pass_len');

if($_POST['ftp_pass'] != $_POST['ftp_pass_2'])
	Components\Page::redirect('../../view.php?id='.$_POST['sid'].'&error=ftp_pass|ftp_pass_2&disp=pass_match');

/*
 * Select Node, User, & Server Information
 */
$core->server->rebuildData($_POST['sid']);
$core->user->rebuildData($core->server->getData('owner_id'));

/*
 * Update Server FTP Information
 */
$iv = $core->auth->generate_iv();
$pass = $core->auth->encrypt($_POST['ftp_pass'], $iv);

$server = ORM::forTable('servers')->findOne($_POST['sid']);
$server->ftp_pass = $pass;
$server->encryption_iv = $iv;
$server->save();

/*
 * Send the User an Email
 */
if(isset($_POST['email_user'])){

    $core->email->buildEmail('admin_new_ftppass', array(
        'PASS' => $_POST['ftp_pass'],
        'SERVER' => $core->server->getData('name')
    ))->dispatch($core->user->getData('email'), Settings::config()->company_name.' - Your FTP Password was Reset');

}

Components\Page::redirect('../../view.php?id='.$_POST['sid']);

?>
