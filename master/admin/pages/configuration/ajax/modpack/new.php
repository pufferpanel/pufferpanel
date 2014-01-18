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

if($core->framework->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->framework->auth->getCookie('pp_auth_token'), true) !== true){
	$core->framework->page->redirect('../../../../index.php');
}

//Cookies :3
setcookie("__TMP_pp_admin_newmodpack", json_encode($_POST), time() + 120, '/', $core->framework->settings->get('cookie_website'));

/*
 * All Posted?
 */
if(!isset($_POST['pack_name'], $_POST['pack_version'], $_POST['pack_minram'], $_POST['pack_permgen'], $_FILES['pack_jar']))
	$core->framework->page->redirect('../../modpacks.php?disp=missing_args&tab=install');

/*
 * Validate Modpack Name
 */
if(!preg_match('/^[\s\w.()-]{1,64}$/', $_POST['pack_name']))
	$core->framework->page->redirect('../../modpacks.php?error=pack_name&disp=pn_fail&tab=install');

/*
 * Validate Min. RAM and Permgen
 */	
if(!is_numeric($_POST['pack_minram']) || !is_numeric($_POST['pack_permgen']))
	$core->framework->page->redirect('../../modpacks.php?error=pack_minram|pack_permgen&disp=num_fail&tab=install');

/*
 * Validate Version
 */	
if(!preg_match('/^[\w.-]{1,64}$/', $_POST['pack_version']))
	$core->framework->page->redirect('../../modpacks.php?error=pack_version&disp=ver_fail&tab=install');

/*
 * File Validation
 */
if(!isset($_FILES['pack_jar']['error']) || is_array($_FILES['pack_jar']['error']))
	$core->framework->page->redirect('../../modpacks.php?error=pack_jar&disp=file_error&tab=install');

switch ($_FILES['pack_jar']['error']) {
	
	case UPLOAD_ERR_OK:
		break;
	case UPLOAD_ERR_NO_FILE:
		$core->framework->page->redirect('../../modpacks.php?error=pack_jar&disp=no_file&tab=install');
	case UPLOAD_ERR_INI_SIZE:
	case UPLOAD_ERR_FORM_SIZE:
		$core->framework->page->redirect('../../modpacks.php?error=pack_jar&disp=file_size&tab=install');
		break;
	default:
		$core->framework->page->redirect('../../modpacks.php?error=pack_jar&disp=file_error&tab=install');
		break;
		
}

/*
 * Limit File Size to 35MB
 */
if($_FILES['pack_jar']['size'] > (1024 * 1024 * 35))
	$core->framework->page->redirect('../../modpacks.php?error=pack_jar&disp=file_size&tab=install');

/*
 * Check File Extension
 */
$finfo = new finfo(FILEINFO_MIME_TYPE);
if($finfo->file($_FILES['pack_jar']['tmp_name']) != "application/zip")
    $core->framework->page->redirect('../../modpacks.php?error=pack_jar&disp=file_type&tab=install');

/*
 * File is Legit, Add Modpack
 */
$modpackHash = $core->framework->auth->keygen(8).'-'.$core->framework->auth->keygen(7);
$isDefault = (isset($_POST['pack_default'])) ? 1 : 0;

$addPack = $mysql->prepare("INSERT INTO `modpacks` VALUES(NULL, :hash, :name, :version, :minram, :permgen, :time, :default, 0)");
$addPack->execute(array(
	':hash' => $modpackHash,
	':name' => $_POST['pack_name'],
	':version' => $_POST['pack_version'],
	':minram' => $_POST['pack_minram'],
	':permgen' => $_POST['pack_permgen'],
	':time' => time(),
	':default' => $isDefault
));

$addpackId = $mysql->lastInsertId();

/*
 * Move File
 */
if(!move_uploaded_file($_FILES['pack_jar']['tmp_name'], sprintf($core->framework->settings->get('modpack_dir').'%s.%s', $modpackHash, "zip"))){
	
	//Delete from DB
	$mysql->exec("DELETE FROM `modpacks` WHERE `id` = '".$addpackId."' LIMIT 1");
	
	//Redirect
	$core->framework->page->redirect('../../modpacks.php?error=pack_jar&disp=file_nomove&tab=install');
	
}else{

	//Reset Default
	if($isDefault == 1)
		$mysql->exec("UPDATE `modpacks` SET `default` = 0 WHERE `id` != '".$addpackId."' LIMIT 1");
		
	//Redirect
	$core->framework->page->redirect('../../modpacks.php');

}