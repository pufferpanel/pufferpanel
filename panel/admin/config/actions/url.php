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
session_start();
require_once('../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php');
}

setcookie("__TMP_pp_admin_updateglobal", json_encode($_POST), time() + 30, '/', $core->settings->get('cookie_website'));

if(!isset($_POST['main_url'], $_POST['master_url'], $_POST['assets_url']))
	Page\components::redirect('../urls.php?error=main_url|master_url|assets_url');

foreach($_POST as $id => $val)
	{

		if(!filter_var($val, FILTER_VALIDATE_URL))
			Page\components::redirect('../urls.php?error='.$id);

	}

$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'main_website'")->execute(array($_POST['main_url']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'master_url'")->execute(array($_POST['master_url']));
$mysql->prepare("UPDATE `acp_settings` SET `setting_val` = ? WHERE `setting_ref` = 'assets_url'")->execute(array($_POST['assets_url']));

Page\components::redirect('../urls.php');

?>
