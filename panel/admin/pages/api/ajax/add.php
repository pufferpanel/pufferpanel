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
require_once('../../../../core/framework/framework.core.php');

if($core->auth->isLoggedIn($_SERVER['REMOTE_ADDR'], $core->auth->getCookie('pp_auth_token'), null, true) !== true){
	Page\components::redirect('../../../../index.php?login');
}

if(!isset($_POST['allowed_ips'], $_POST['permissions']))
	Page\components::redirect('../view.php?tab=add&error=missing_values');
	
if(empty($_POST['permissions']))
	Page\components::redirect('../view.php?tab=add&error=no_perms');
	
/*
 * Build IP Data
 */
$ips = array();
if(empty($_POST['allowed_ips']) || $_POST['allowed_ips'] = '*')
	$ips = array("*");
else {

	foreach(explode(',', $_POST['allowed_ips']) as $id => $ip){
	
		$ip = trim($ip);
		
		if(!filter_var($ip, FILTER_VALIDATE_IP))
			Page\components::redirect('../view.php?tab=add&error=invalid_ip');
			
		$ips = array_merge($ips, array($ip));
	
	}

}

/*
 * Build Permissions Data
 */
$perms = array();
foreach($_POST['permissions'] as $id => $perm){

	$perms = array_merge($perms, array(trim(strtolower($perm))));

}

$add = $mysql->prepare("INSERT INTO `api` VALUES (NULL, :key, :perms, :ips)");
$add->execute(array(
	':key' => $core->auth->keygen(10).'-'.$core->auth->keygen(5).'-'.$core->auth->keygen(5).'-'.$core->auth->keygen(14),
	':perms' => json_encode($perms),
	':ips' => json_encode($ips)
));

Page\components::redirect('../view.php');

?>