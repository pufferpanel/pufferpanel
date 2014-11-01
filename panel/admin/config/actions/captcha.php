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

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true)
	Components\Page::redirect('../../../index.php');

if(!isset($_POST['pub_key'], $_POST['priv_key']) || empty($_POST['pub_key']) || empty($_POST['priv_key']))
	Components\Page::redirect('../captcha.php?error=pub_key|priv_key');

$update = ORM::forTable('acp_settings')
			->rawExecute("UPDATE acp_settings SET setting_val = IF(setting_ref='captcha_pub', :pub, :priv) WHERE setting_ref IN ('captcha_pub', 'captcha_priv')",
			array(
				'pub' => $_POST['pub_key'],
				'priv' => $_POST['priv_key']
			));

Components\Page::redirect('../captcha.php');
?>
