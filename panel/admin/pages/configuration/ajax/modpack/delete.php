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
require_once('../../../../../../src/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../../index.php');
}

/*
 * All Posted?
 */
if(!isset($_POST['pack_del'], $_POST['pack_newdefault'], $_POST['conf_pack_hash'], $_POST['confirm_delete']))
	Page\components::redirect('../../edit.php?mid='.$_POST['pack_del'].'&error=pack_delete&disp=missing_params');

/*
 * Does the Pack Delete Match?
 */
if($_POST['pack_del'] != $_POST['conf_pack_hash'])
	Page\components::redirect('../../edit.php?mid='.$_POST['pack_del'].'&error=pack_delete&disp=pack_hash_mismatch');
	
$getPack = $mysql->prepare("SELECT * FROM `modpacks` WHERE `hash` = :hash");
$getPack->execute(array(
	':hash' => $_POST['pack_del']
));

	if($getPack->rowCount() != 1)
		Page\components::redirect('../../edit.php?mid='.$_POST['pack_del'].'&error=pack_delete&disp=pack_hash');
	else
		$pack = $getPack->fetch();

/*
 * Do we need a new default.
 */		
if($pack['default'] == 1 && !isset($_POST['pack_newdefault']) || $_POST['pack_newdefault'] == 'no-continue')
	Page\components::redirect('../../edit.php?mid='.$_POST['pack_del'].'&error=pack_delete&disp=no_new_default');
	
	if($pack['default'] == 1){
	
		/*
		 * Check that new default is legit.
		 */
		$checkHash = $mysql->prepare("SELECT `hash` FROM `modpacks` WHERE `hash` = :hash");
		$checkHash->execute(array(
			':hash' => $_POST['pack_newdefault']
		));
		
		if($checkHash->rowCount() != 1)
			Page\components::redirect('../../edit.php?mid='.$_POST['pack_del'].'&error=pack_delete&disp=new_default_noexist');
		
	}

/*
 * Delete Modpack from MySQL
 */
$deletePack = $mysql->prepare("UPDATE `modpacks` SET `deleted` = 1 WHERE `hash` = :oldpack LIMIT 1");
$deletePack->execute(array(
	':oldpack' => $_POST['pack_del']
));

/*
 * Set New Default
 */
 $deletePack = $mysql->prepare("UPDATE `modpacks` SET `default` = 1 WHERE `hash` = :newpack LIMIT 1");
 $deletePack->execute(array(
 	':newpack' => $_POST['pack_newdefault']
 ));
 
 /*
  * Delete Modpack from Server
  */
 if(file_exists($core->settings->get('modpack_dir').$_POST['pack_del'].'.zip') && is_readable($core->settings->get('modpack_dir').$_POST['pack_del'].'.zip'))
 	unlink($core->settings->get('modpack_dir').$_POST['pack_del'].'.zip');
 	
 Page\components::redirect('../../modpacks.php');
