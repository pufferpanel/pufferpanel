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
require_once('../../../../../src/framework/framework.core.php');


if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../index.php?login');
}

if(isset($_GET['do']) && $_GET['do'] == 'ipuser') {

	if(!isset($_POST['nid']) || !is_numeric($_POST['nid']))
		Page\components::redirect('../../list.php');

	if(!filter_var($_POST['ip'] , FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE))
		Page\components::redirect('../../view.php?id='.$_POST['nid'].'&error=ip&disp=ip_fail&tab=ftp');

	/*
	 * Run Update on Node Table
	 */
	$mysql->prepare("UPDATE `nodes` SET `ip` = :ip WHERE `id` = :nid")->execute(array(':ip' => $_POST['ip'], ':nid' => $_POST['nid']));
	Page\components::redirect('../../view.php?id='.$_POST['nid'].'&tab=ftp');

}

?>
