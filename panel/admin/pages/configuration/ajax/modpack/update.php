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
require_once('../../../../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	$core->page->redirect('../../../../index.php');
}

//Cookies :3
setcookie("__TMP_pp_admin_updatemodpack", json_encode($_POST), time() + 120, '/', $core->settings->get('cookie_website'));

/*
 * All Posted?
 */
if(!isset($_POST['pack_name'], $_POST['pack_hash'], $_POST['pack_version'], $_POST['pack_minram'], $_POST['pack_permgen']))
	$core->page->redirect('../../edit.php?mid='.$_POST['pack_hash'].'&disp=missing_args');

/*
 * Validate Modpack Name
 */
if(!preg_match('/^[\s\w.()-]{1,64}$/', $_POST['pack_name']))
	$core->page->redirect('../../edit.php?mid='.$_POST['pack_hash'].'&error=pack_name&disp=pn_fail');

/*
 * Validate Modpack Jar Name
 */
//if(!preg_match('/^[\s\w.-]{1,64}$/', $_POST['server_jar']))
//	$core->page->redirect('../../edit.php?mid='.$_POST['pack_hash'].'&error=server_jar&disp=pn_fail');


/*
 * Validate Min. RAM and Permgen
 */	
if(!is_numeric($_POST['pack_minram']) || !is_numeric($_POST['pack_permgen']))
	$core->page->redirect('../../edit.php?mid='.$_POST['pack_hash'].'&error=pack_minram|pack_permgen&disp=num_fail');

/*
 * Validate Version
 */	
if(!preg_match('/^[\w.-]{1,64}$/', $_POST['pack_version']))
	$core->page->redirect('../../edit.php?mid='.$_POST['pack_hash'].'&error=pack_version&disp=ver_fail');

$mysql->prepare("UPDATE `modpacks` SET `name` = ? WHERE `hash` = ? ")->execute(array($_POST['pack_name'], $_POST['pack_hash']));
$mysql->prepare("UPDATE `modpacks` SET `version` = ? WHERE `hash` = ? ")->execute(array($_POST['pack_version'], $_POST['pack_hash']));
$mysql->prepare("UPDATE `modpacks` SET `min_ram` = ? WHERE `hash` = ? ")->execute(array($_POST['pack_minram'], $_POST['pack_hash']));
$mysql->prepare("UPDATE `modpacks` SET `permgen` = ? WHERE `hash` = ? ")->execute(array($_POST['pack_permgen'], $_POST['pack_hash']));
$mysql->prepare("UPDATE `modpacks` SET `server_jar` = ? WHERE `hash` = ? ")->execute(array($_POST['server_jar'], $_POST['pack_hash']));


if(isset($_POST['pack_default'])){
	$mysql->exec("UPDATE `modpacks` SET `default` = 0");
	$mysql->prepare("UPDATE `modpacks` SET `default` = 1 WHERE `hash` = ?")->execute(array($_POST['pack_hash']));
}
	
//Redirect
$core->page->redirect('../../edit.php?mid='.$_POST['pack_hash'].'&success=true');
